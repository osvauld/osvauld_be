# Osvauld
An open-source credential management tool intended to solve password-mess teams face. Osvauld helps to enforce password/token hygiene and visibilty across your team thus resulting in better overall security.

## Start osvauld instance
### Database
Osvauld requires a postgres instance running. 

### Pull image
```docker pull osvauld/backend:0.1.0```

### Start Instance
```
docker run --name osvauld_be \
    -d \
    -p 80:8000 \
    -e MASTER_DB_HOST=<db_host> \
    -e MASTER_DB_NAME=<db_name> \
    -e MASTER_DB_USER=<db_user> \
    -e MASTER_DB_PASSWORD=<db_password> \
    -e MASTER_DB_PORT=<db_port> \
    -e MASTER_SSL_MODE=require \
    osvauld/backend:0.1.0
```

