

ssh ubuntu@3.110.128.10 << 'EOF'

    source setup_env.sh

    cd osvauld_be

    git checkout main
    
    git pull

    sudo docker build -t osvauld_be:latest .

    sudo docker stop osvauld_backend

    sudo docker rm osvauld_backend

    sudo docker rmi osvauld_be:latest


    sudo docker run --name osvauld_backend \
    -d \
    -p 80:8000 \
    -e MASTER_DB_HOST=$MASTER_DB_HOST \
    -e MASTER_DB_NAME=$MASTER_DB_NAME \
    -e MASTER_DB_USER=$MASTER_DB_USER \
    -e MASTER_DB_PASSWORD=$MASTER_DB_PASSWORD \
    -e MASTER_DB_PORT=$MASTER_DB_PORT \
    -e MASTER_SSL_MODE=require \
    osvauld_be:latest
EOF

