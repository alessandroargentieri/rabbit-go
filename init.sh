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
    	echo "... starting container"
    fi
else
	docker run -d --name my-rabbit -e RABBITMQ_DEFAULT_USER=myuser -e RABBITMQ_DEFAULT_PASS=password -p 5672:5672 -p 15672:15672 rabbitmq:3.8-management
fi

# build executable
cd ./publisher && go build -o rabbit-publisher && cd ../
cd ./consumer && go build -o rabbit-consumer && cd ../

# export env vars
export RABBITMQ_USER=myuser
export RABBITMQ_PASSWORD=password
export RABBITMQ_URL=localhost:5672

# launch 1 publisher app
export PORT=8080
./publisher/rabbit-publisher &

# launch 3 consumer apps
export PORT=8081
./consumer/rabbit-consumer &

export PORT=8082
./consumer/rabbit-consumer &

export PORT=8083
./consumer/rabbit-consumer &
