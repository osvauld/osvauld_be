

ssh ubuntu@13.127.6.232 << EOF

    cd osvauld_be

    git pull

    sudo docker stop osvauld_backend

    sudo docker rm osvauld_backend

    sudo docker rmi osvauld_be:latest

    sudo docker build -t osvauld_be:latest .

    sudo docker run -d -p 8000:8000 --name osvauld_backend osvauld_be:latest
EOF

