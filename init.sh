#! /bin/bash

# pull docker image of rabbitmq and launch it
if [[ $( docker images | grep rabbitmq) ]]; then
    echo "rabbitmq image already pulled."
else
	echo "pulling rabbitmq docker image"
	docker pull rabbitmq:3.8-management    
fi

if [[ $( docker ps -a | grep my-rabbit | head -c12 ) ]]; then
    echo "my-rabbit container already present..."
    if [[ $(docker ps | grep my-rabbit | head -c12 ) ]]; then 
    	echo "...and running!"
    else
    	docker start my-rabbit
    	echo "...starting container"
    fi
else
	docker run -d --name my-rabbit -e RABBITMQ_DEFAULT_USER=myuser -e RABBITMQ_DEFAULT_PASS=password -p 5672:5672 -p 15672:15672 rabbitmq:3.8-management
fi

# build executable
echo "...building Go executables..."
cd ./publisher && go build -o rabbit-publisher && cd ../
cd ./consumer && go build -o rabbit-consumer && cd ../

echo "wait until the rabbitmq instance is ready"
sleep 10s

# export env vars
export RABBITMQ_USER=myuser
export RABBITMQ_PASSWORD=password
export RABBITMQ_URL=localhost:5672

# launch 1 publisher app
echo "...starting producer app..."
export PORT=8080
./publisher/rabbit-publisher >> publisher-logs.txt 2>&1 &

# launch 3 consumer apps
echo "...starting 3 consumer app..."
export PORT=8081
./consumer/rabbit-consumer >> consumer-1-logs.txt 2>&1 &

export PORT=8082
./consumer/rabbit-consumer >> consumer-2-logs.txt 2>&1 &

export PORT=8083
./consumer/rabbit-consumer >> consumer-3-logs.txt 2>&1 &
