## Go Microservice Test

This project aim is to create a Go microservice with HTTP endpoints, command line utility and Docker.


### Setup

Before procceeding with setup make sure you have at least basic understanding of Go, Make (GNU), MySQL (MariaDB) and Docker (optional). This project was intended to be built and ran on MacOS and Linux. Windows is not supported.


#### Requirement:

1. Go-lang
2. MySQL (MariaDB) database
3. GNU Make
4. Docker (optional)


#### Usage:

1. Start database server. Create database in MySQL (MariaDB) and update `config.toml` configuration and run setup command which will create all database tables. MariaDB can be setup with docker and docker compose as well. See `docker-compose.yml` file.
    ```
    make setup
    ```

2. Fetch and update latest currency feed.
    ```
    make update
    ```

3. Start local HTTP server. Access url/address is set in `config.toml` file.
    ```
    make serve
    ```

4. Build binary executable. The built excutable will be located in `dist` directory. Executable should be shipped together with `config.toml` file and should be in the same directory.
    ```
    make build
    ```

#### Docker:

Microservice can be built and ran inside docker container. Before building docker image make sure `config.toml` is all set and located in project's root directory alongside the Dockerfile. Make sure HTTP server address is set to `0.0.0.0:3035`, otherwise endpoints won't be accessible outside the container.

1. Build docker image.
    ```
    make build-docker
    ```

2. Run docker image. By default it will run `serve` command and will start the HTTP server.
    ```
    make run-docker
    ```

#### HTTP Endpoints:

There are two endpoints that can be accessed by an HTTP request. 

1. Get latest currencies.
    ```
    GET /currency/latest
    ```

3. Get currency historical values by symbol. There are 32 available currencies. `AUD BGN BRL CAD CHF CNY CZK DKK GBP HKD HRK HUF IDR ILS INR ISK JPY KRW MXN MYR NOK NZD PHP PLN RON RUB SEK SGD THB TRY USD ZAR`. Replace `:symbol` with one of the currency symbols to access data for particular currency. Returned values can be limited by adding query param `?limit=20`. Limit is in range from 1 to 100.
    ```
    GET /currency/:symbol
    ```

#### Test:

Run unit tests
```
make test
```

Test coverage
```
make coverage
```


### License

Copyright (c) 2020-present [matiss](https://github.com/matiss). Go Microservice Test is free and open-source software licensed under the MIT License.