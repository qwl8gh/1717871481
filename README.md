```
docker rm -f $(docker ps -a -q)
docker rmi -f $(docker images -aq)
docker volume rm $(docker volume ls -q)

sudo kill -9 $(sudo lsof -t -i:8080)

docker compose up --build
```

|![1](imgs/1.png)|![5](imgs/5.png)|
|:---:|:---:|

|![2](imgs/2.png)|![3](imgs/3.png)|![4](imgs/4.png)|
|:---:|:---:|:---:|