# SSOA-PT REST APP

## Services

- Backend: REST API in Go
  - Database: PostgreSQL
  - Web Proxy: Nginx
    - Let's Encrypt HTTPS certificates with certbot
- Frontend: PHP App

## How to run the app

Install `docker` https://docs.docker.com/engine/install/ and `docker-compose` https://docs.docker.com/compose/install/.

Clone this repo to your machine.

Before the first run, get a Let's Encrypt certificate by executing the bash script `./init-letsencrypt.sh`.

After that, you can run the app via `docker-compose up -d`.

## Optional configuration

By default, the app gets the current public ip address and creates a domain name with the help of `nip.io`.
If you want to use your own dns domain, supply the environment variable `PUBLIC_SERVER_NAME` to the init script `./init-letsencrypt.sh`, 
or by setting it in a `.env` file in the root project directory.

```
PUBLIC_SERVER_NAME=<public_dns_name>
AWS=no
STAGING=1
```

If you are using AWS cloud, set `AWS=yes`, PUBLIC_SERVER_NAME then is not needed. 

Set `STAGING` to `0` if you want to use Let's Encrypt without `--staging` (production mode).

Set the database password in the file `./db/password.txt`

Lastly, have a look at the main configuration file [_docker-compose.yaml_](docker-compose.yaml).

## Deploy with docker-compose

Run the stack
```
docker-compose up -d
```

View logs
```
docker-compose logs
```

Restart the stack
```
docker-compose restart
```

Stop and remove the containers
```
$ docker-compose down
```

All persistent data will be stored in `./data/`. Remove this directory too.
