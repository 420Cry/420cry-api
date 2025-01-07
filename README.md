# 420cry-api

This is a Go-based API server for the 420cry application.
## Prerequisites

- Go 1.23.4

## Preparation

1. **Add Development Hosts to the `/etc/hosts` File**:
    * On **Linux/macOS**, add the following lines to your `/etc/hosts` file.
    * On **Windows**, add them to the `C:\Windows\System32\drivers\etc\hosts` file.

    Add the following lines to the file:
    ```bash
    127.0.0.1 api.420.crypto.test
    ```

## Installation

1. Clone the repository
2. Install Go dependencies:
    ```bash
    make install
    ```
3. Build the Go application and create a binary:
    ```bash
    make build
    ```
4. Run the server:
    ```bash
    make dev
    ```
### Lint:

This project uses `golangci-lint`. You can install it using the following commands based on your OS:

#### macOS:
You can install it with `brew`:
```bash
brew install golangci-lint
```

### Linux & Win
You can install it with curl:
```bash
curl -sSfL https://github.com/golangci/golangci-lint/releases/download/v1.52.0/golangci-lint-1.52.0-linux-amd64.tar.gz | tar -xz -C /tmp
sudo mv /tmp/golangci-lint-1.52.0-linux-amd64/golangci-lint /usr/local/bin/

```

#### Run Lint:
You can install it with `brew`:
```bash
make lint
```

### Test:
```bash
make test
```
### With Docker
1. Shutdown the dev server docker compose for this project.
    ```bash
    docker compose down
    ```

2. Build and start application in production mode.
    ```bash
    docker compose build
    ```

3. Start the application in DEV mode.
    ```bash
    docker compose up -d
   ```