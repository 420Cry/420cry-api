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
    127.0.0.1 420.api.crypto.test
    127.0.0.1 420.db.crypto.test
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
## Docker
