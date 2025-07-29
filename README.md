## build

```
go build -o <app name> ./cmd/frequentNumbers
```

```
go build -o <app name> ./cmd/result
```

## Docker

This project includes Docker support for easy deployment and execution.

### Building the Docker image

```bash
docker build -t loto7-promise .
```

### Running with Docker

Run the result application:
```bash
# Show usage
docker run --rm loto7-promise

# Get past 10 results  
docker run --rm loto7-promise ./result 10
```

Run the frequent numbers application:
```bash
docker run --rm loto7-promise ./frequentNumbers
```

### Using Docker Compose

For easier management, use Docker Compose:

```bash
# Build all services
docker compose build

# Run result service with argument
docker compose run --rm result ./result 5

# Run frequent numbers service
docker compose run --rm frequent

# Run interactive shell
docker compose run --rm interactive
```

### Image Details

- **Base images**: golang:1.23-alpine (build) + alpine:latest (runtime)
- **Image size**: ~39MB
- **Applications**: Both `result` and `frequentNumbers` binaries included
- **Security**: Runs as non-root user `appuser`
- **Dependencies**: All Go dependencies vendored for offline builds
