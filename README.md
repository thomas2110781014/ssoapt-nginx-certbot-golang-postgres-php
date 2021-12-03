## SSOA-PT REST APP

### Services

- Backend: REST API in Go
  - Database: PostgreSQL
  - Web Proxy: Nginx
    - Let's Encrypt HTTPS certificates with certbot
- Frontend: PHP App

## Set the environment in the file `.env`

```
PUBLIC_SERVER_NAME=<public_dns_name>
AWS=no
STAGING=1
```

If you are using AWS cloud, set `AWS=yes`, PUBLIC_SERVER_NAME then is not needed. 
Set `STAGING` to `0` if you want to use Let's Encrypt without `--staging` (production mode).

## Set the database password in the file `./db/password.txt`

## Set also environment in the file `docker-compose.yaml`

## Get a Let's Encrypt HTTPS cert with certbot

```
./init-letsencrypt.sh
```

Run this for initialization to get the certificate.

You can edit the script and set `staging=0` to disable staging (testing).

## Deploy with docker-compose

[_docker-compose.yaml_](docker-compose.yaml)

```
docker-compose up -d
```

Stop and remove the containers
```
$ docker-compose down
```
