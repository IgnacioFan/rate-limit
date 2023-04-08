# Rate limit
The Rate limit helps control the rate of requests sent to an API. It accurately limits excessive requests and provides clear exceptions when throttling occurs. It also supports various rate limit strategies such as the token bucket, leaking bucket, fixed window, and sliding window algorithms. It can handle high-volume requests with minimal delay and resource usage. It's highly available, memory efficient and low latent.

## Technologies Used
- `Gin` as the web framework
- `Cobra` as the command-line interface framework

## Installation and Setup

### Prerequisite
- Go 1.16+
- Docker, we use `docker-compose` to boot up all required services

### Run it via Docker compose

1. Clone this repository to local
2. Navigate to the repo
```
// start all services
docker compose up

// stop all services
docker compose down
```

## Project Structure
```
├── build/
├── cmd/
│   ├── root.go
│   └── server.go
├── internal/
│   ├── delivery/
│   │   ├── handler/
│   │   └── http/
│   │       ├── server.go
│   │       └── util.go
│   ├── conn/
│   │   └── redis/
│   │       ├── impl.go
│   │       └── redis.go
│   └── services/
│       ├── base/
│       └── ratelimiter/
│           ├── strategy/
│           │     ├── token_bucket/
│           │     ├── leaky_bucket/
│           │     ├── fixed_window_counter/
│           │     ├── sliding_window_logs/
│           │     └── strategy.go
│           ├── impl.go
│           └── ratelimiter.go
├── scripts
└── main.go
```


## How to use?

## Contributors
- Weilong Fan: developer and maintainer

## References
- https://github.com/cosmtrek/air (hot reload)
- https://systemsdesign.cloud/SystemDesign/RateLimiter
