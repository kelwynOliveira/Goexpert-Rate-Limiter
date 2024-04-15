# Rate Limiter

Goexpert postgraduation project

## Lab project Description

> **Objective**: To develop a rate limiter in Go that can be configured to limit the maximum number of requests per second based on a specific IP address or access token.
>
> ### Description: 
> The aim of this challenge is to create a rate limiter in Go that can be used to control the traffic of requests to a web service. The rate limiter must be able to limit the number of requests based on two criteria:
>
> 1. **IP address**: The rate limiter must restrict the number of requests received from a single IP address within a defined time interval.
> 1. **Access Token**: The rate limiter must also be able to limit requests based on a unique access token, allowing different expiration time limits for different tokens. The token must be entered in the header in the following format:
>     - `API_KEY: <TOKEN>`
> The limit settings for the access token must override those for the IP. Ex: If the limit per IP is 10 req/s and that of a given token is 100 req/s, the rate limiter must use the token's information.
> 
> ### Requirements
> - The rate limiter must be able to work as middleware that is injected into the web server
> - The rate limiter must be able to configure the maximum number of requests allowed per second.
> - The rate limiter must have the option of choosing the time to block the IP or token if the number of requests has been exceeded.
> - Limit settings must be made via environment variables or in an ".env" file in the root folder.
> - It must be possible to configure the rate limiter for both IP and token limitation.
> - The system should respond appropriately when the limit is exceeded:
>     - HTTP code: 429
>     - Message: you have reached the maximum number of requests or actions allowed within a certain time frame.
> - All limiter information must be stored in and queried from a Redis database. You can use docker-compose to upload Redis.
> - Create a strategy that allows you to easily switch from Redis to another persistence mechanism.
> - The limiter logic must be separate from the middleware.
>
> ### Examples
> 1. Limiting by IP: Suppose the rate limiter is configured to allow a maximum of 5 requests per second per IP. If IP 192.168.1.1 sends 6 requests in one second, the sixth request should be blocked.
> 1. Token limitation: If an abc123 token has a configured limit of 10 requests per second and sends 11 requests in that interval, the eleventh should be blocked.
> 1. In the two cases above, the next requests can only be made when the full expiration time has elapsed Ex: If the expiration time is 5 minutes, a given IP will only be able to make new requests after 5 minutes.
>
> ### Tips
> - Test your rate limiter under different load conditions to ensure that it works as expected in high-traffic situations.
> 
> ### Delivery
> - The complete source code of the implementation.
> - Documentation explaining how the rate limiter works and how it can be configured.
> - Automated tests demonstrating the effectiveness and robustness of the rate limiter.
> - Use docker/docker-compose so that we can test your application.
> - The web server should respond on port 8080.

## How the rate limiter works

The application consists of a web server that receives HTTP requests and a rate limiter middleware that is responsible for controlling the number of requests received. The middleware intercepts all requests and executes the rate limiting logic from the instance of `RateLimiter`, which contains the limiter's business rules and knows how to invoke the storage _Strategy_ instantiated by the dependency manager at `internal/pkg/dependencyinjector/injector.go` to perform the limit check.
 
The rate limiter can be configured to check limits by IP or `API_KEY` token, and uses Redis as _storage_ to store the number of requests made by each IP and/or token. This configuration is carried out using environment variables declared in the `.env` file and injected into the application via the dependency manager, at application boot.

The storage strategy is defined via a `LimiterStrategyInterface` interface which has a `Check` method for obtaining and setting values in the _storage_. At the moment, the application only has one implementation for Redis, but it is possible to add new implementations for other _storages_ such as memory, database, etc, without changing the rate limiting logic, just by injecting the new implementation into the `RateLimiter` instance via dependency manager.

## How to configure

Inside the [`.env`](.env) file in the root update the content and adjust it as necessary.

### How to execute

After updating the [`.env`](.env) file run the command `docker compose up redis api` to start the application and Redis.

### How to make requests

Request with IP check: `$ curl -vvv http://localhost:8080`
 
Request with token check: `$ curl -H 'API_KEY: <TOKEN>' -vvv http://localhost:8080`

## Automated tests

### Unit test
 
To run the unit tests and validate the coverage, run the `make test` command.
 
### Stress tests
To run the stress tests with k6, follow these steps:
1. start the application and Redis with the command `docker compose up redis api`;
2. Run the command `make test_k6_smoke` to start the _smoke_ stress test (duration 1 minute);
3. Run the `make test_k6_stress` command to start the _stress_ stress test (duration 40 minutes).

You can view the results in the `./scripts/k6/smoke` and `./scripts/k6/stress` folders, both in text and HTML format.
 
 

