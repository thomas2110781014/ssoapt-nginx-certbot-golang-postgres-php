#!/bin/bash

if ! [ -x "$(command -v docker-compose)" ]; then
  echo 'Error: docker-compose is not installed.' >&2
  exit 1
fi

source .env

# if using AWS cloud, get DNS name 
if [ "$AWS" = "yes" ] ; then
  PUBLIC_SERVER_NAME=$(curl http://169.254.169.254/latest/meta-data/public-hostname)
else
# else use nip.io to get a dns name
  if [ -z $PUBLIC_SERVER_NAME ] ; then
    PUBLIC_IPV4_ADDRESS=$(curl https://ipinfo.io/ip)
    PUBLIC_SERVER_NAME=$PUBLIC_IPV4_ADDRESS.nip.io
  fi
fi
sed -i "s/public_dns_name_here/$PUBLIC_SERVER_NAME/" docker-compose.yaml

staging=${STAGING:-1} # Let's Encrypt, set to 1 if you're testing your setup to avoid hitting request limits, set to 0 for production
domains=${PUBLIC_SERVER_NAME}
#domains=(${PUBLIC_SERVER_NAME} www.${PUBLIC_SERVER_NAME})
rsa_key_size=4096
email="${EMAIL_ADDRESS:-guest@hotel.com}"
data_path="./data/certbot"

if [ -d "$data_path" ]; then
  read -p "Existing data found for $domains. Continue and replace existing certificate? (y/N) " decision
  if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
    exit
  fi
fi

grep 'NGINX_HTTP' .env ||
  cat << EOF >> .env
NGINX_HTTP_PORT=80
NGINX_HTTPS_PORT=443
EOF


# start other services, otherwise nginx will fail
docker-compose up -d db
docker-compose up -d apache-php
docker-compose up -d backend 
sleep 4

if [ ! -e "$data_path/conf/options-ssl-nginx.conf" ] || [ ! -e "$data_path/conf/ssl-dhparams.pem" ]; then
  echo "### Downloading recommended TLS parameters ..."
  mkdir -p "$data_path/conf"
  curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot-nginx/certbot_nginx/_internal/tls_configs/options-ssl-nginx.conf > "$data_path/conf/options-ssl-nginx.conf"
  curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot/certbot/ssl-dhparams.pem > "$data_path/conf/ssl-dhparams.pem"
  echo
fi

echo "### Creating dummy certificate for $domains ..."
path="/etc/letsencrypt/live/$domains"
mkdir -p "$data_path/conf/live/$domains"
docker-compose run --rm --entrypoint "\
  openssl req -x509 -nodes -newkey rsa:$rsa_key_size -days 1\
    -keyout '$path/privkey.pem' \
    -out '$path/fullchain.pem' \
    -subj '/CN=localhost'" certbot
echo


echo "### Starting nginx ..."
docker-compose up --force-recreate -d nginx
echo

echo "### Deleting dummy certificate for $domains ..."
docker-compose run --rm --entrypoint "\
  rm -Rf /etc/letsencrypt/live/$domains && \
  rm -Rf /etc/letsencrypt/archive/$domains && \
  rm -Rf /etc/letsencrypt/renewal/$domains.conf" certbot
echo


echo "### Requesting Let's Encrypt certificate for $domains ..."
#Join $domains to -d args
domain_args=""
for domain in "${domains[@]}"; do
  domain_args="$domain_args -d $domain"
done

# Select appropriate email arg
case "$email" in
  "") email_arg="--register-unsafely-without-email" ;;
  *) email_arg="--email $email" ;;
esac

# Enable staging mode if needed
if [ $staging != "0" ]; then staging_arg="--staging"; fi

docker-compose run --rm --entrypoint "\
  certbot certonly --webroot -w /var/www/certbot \
    $staging_arg \
    $email_arg \
    $domain_args \
    --rsa-key-size $rsa_key_size \
    --agree-tos \
    --non-interactive \
    --force-renewal" certbot
echo


docker-compose down
echo "To start:  docker-compose up -d"
