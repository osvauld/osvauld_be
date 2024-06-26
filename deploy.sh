ssh ubuntu@3.110.128.10 << 'EOF'

    source setup_env.sh

    cd osvauld_be

    git checkout develop
    
    git pull

    sudo docker stop osvauld_backend

    sudo docker rm osvauld_backend

    sudo docker rmi osvauld_be:latest

    sudo docker build -t osvauld_be:latest .

    sudo docker run --name osvauld_backend \
    --restart=unless-stopped \
    -d \
    -p 80:8000 \
    -e MASTER_DB_HOST=$MASTER_DB_HOST \
    -e MASTER_DB_NAME=$MASTER_DB_NAME \
    -e MASTER_DB_USER=$MASTER_DB_USER \
    -e MASTER_DB_PASSWORD=$MASTER_DB_PASSWORD \
    -e MASTER_DB_PORT=$MASTER_DB_PORT \
    -e MASTER_SSL_MODE=require \
    -e AUTH_SECRET=$AUTH_SECRET \
    osvauld_be:latest
EOF