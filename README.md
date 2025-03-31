# Rate Limiter

Goexpert postgraduation project

## Lab project Description

> **Objective**: To develop a rate limiter in Go that can be configured to limit the maximum number of requests per second based on a specific IP address or access token.
>
> ### Description:
>
> The aim of this challenge is to create a rate limiter in Go that can be used to control the traffic of requests to a web service. The rate limiter must be able to limit the number of requests based on two criteria:
>
> 1. **IP address**: The rate limiter must restrict the number of requests received from a single IP address within a defined time interval.
> 1. **Access Token**: The rate limiter must also be able to limit requests based on a unique access token, allowing different expiration time limits for different tokens. The token must be entered in the header in the following format: - `API_KEY: <TOKEN>`
>    The limit settings for the access token must override those for the IP. Ex: If the limit per IP is 10 req/s and that of a given token is 100 req/s, the rate limiter must use the token's information.
>
> ### Requirements
>
> - The rate limiter must be able to work as middleware that is injected into the web server
> - The rate limiter must be able to configure the maximum number of requests allowed per second.
> - The rate limiter must have the option of choosing the time to block the IP or token if the number of requests has been exceeded.
> - Limit settings must be made via environment variables or in an ".env" file in the root folder.
> - It must be possible to configure the rate limiter for both IP and token limitation.
> - The system should respond appropriately when the limit is exceeded:
>   - HTTP code: 429
>   - Message: you have reached the maximum number of requests or actions allowed within a certain time frame.
> - All limiter information must be stored in and queried from a Redis database. You can use docker-compose to upload Redis.
> - Create a strategy that allows you to easily switch from Redis to another persistence mechanism.
> - The limiter logic must be separate from the middleware.
>
> ### Examples
>
> 1. Limiting by IP: Suppose the rate limiter is configured to allow a maximum of 5 requests per second per IP. If IP 192.168.1.1 sends 6 requests in one second, the sixth request should be blocked.
> 1. Token limitation: If an abc123 token has a configured limit of 10 requests per second and sends 11 requests in that interval, the eleventh should be blocked.
> 1. In the two cases above, the next requests can only be made when the full expiration time has elapsed Ex: If the expiration time is 5 minutes, a given IP will only be able to make new requests after 5 minutes.
>
> ### Tips
>
> - Test your rate limiter under different load conditions to ensure that it works as expected in high-traffic situations.
>
> ### Delivery
>
> - The complete source code of the implementation.
> - Documentation explaining how the rate limiter works and how it can be configured.
> - Automated tests demonstrating the effectiveness and robustness of the rate limiter.
> - Use docker/docker-compose so that we can test your application.
> - The web server should respond on port 8080.

## How to configure

Inside the [`.env`](.env) file in the root update the content and adjust it as necessary.

### How to execute

After updating the [`.env`](.env) file run the command `make up` or `docker compose up -d`.

### How to make requests

Request with IP check: `make run SCENARIO=ip`

Request with token check: `make run SCENARIO=token`

## Automated tests

To run the tests and validate the coverage, run the `make test` command.
