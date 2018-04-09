docker stop $(docker ps -a -q)

sleep 15

docker rm $(docker ps -a -q)

sleep 5

docker volume rm $(docker volume ls -q)

sleep 5


