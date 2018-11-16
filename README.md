# node-monitor

## requirement
- ```docker``` and ```docker-compose``` 

## deploy
- clone this repository to your local computer
- open terminal and go to repository's folder
- on termial, run command ```docker-compose up --build```

## api
- ```/api/stress-data```: get all stress test data.
- ```/api/normal-data:```: get normal test data from ```fromTime``` to ```toTime```.
- ```/api/stat-normal-data```: get stat data normal test from the beginning
