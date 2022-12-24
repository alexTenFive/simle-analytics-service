# Simple analytics service (test)
### Run
##### First terminal window
```
cd simple-analytics-service
./service_up.sh
```
this will up database, application and run migrations
##### Second terminal window
```
cd simple-analytics-service
./test_load.sh
```
this will run http load test (300 requests per second) for 60 seconds and place `plot.html` file with test statistics in `simple-analytics-service/cmd/vegeta/reports`

### Scripts
```
./service_up.sh
```
Run `docker-compose up`, builds application and postgres database and runs `docker-compose logs -f`
```
./service_down.sh
```
Run `docker-compose down` and remove containers
```
./service_restart.sh
```
For down and up all services
```
./test_load.sh
```
Run [vegeta](https://github.com/tsenart/vegeta) client for simple http load test
