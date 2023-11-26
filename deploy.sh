

ssh ubuntu@13.127.6.232 << EOF

    cd osvauld-web/

    git pull

    sudo docker stop osvauldweb

    sudo docker rm osvauldweb

    sudo docker rmi osvauld-web:latest

    sudo docker build -t osvauld-web:latest .

    sudo docker run -d -p 8000:8000 --name osvauldweb osvauld-web:latest
EOF

