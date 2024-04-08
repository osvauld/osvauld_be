

ssh ubuntu@3.110.128.10 << EOF

    cd osvauld_be

    git checkout main

    git pull

    sudo docker build -t osvauld_be:latest .

    sudo docker stop osvauld_backend

    sudo docker rm osvauld_backend

    sudo docker rmi osvauld_be:latest

    sudo docker run -d -p 80:8000 --name osvauld_backend osvauld_be:latest
EOF

