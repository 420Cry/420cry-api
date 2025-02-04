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
2. **Copy .env.example to .env**:
    ```bash
    cp .env.example .env
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
4. Migration:
    ```bash
    make migrate
    ```
5. Run the server:
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

#### Linux & Windows:
You can install it with curl:
```bash
curl -sSfL https://github.com/golangci/golangci-lint/releases/download/v1.52.0/golangci-lint-1.52.0-linux-amd64.tar.gz | tar -xz -C /tmp
sudo mv /tmp/golangci-lint-1.52.0-linux-amd64/golangci-lint /usr/local/bin/

```

#### Run Lint:
You can run the linter with:
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

### Without Docker
1. Log into MySQL.
    ```bash
    mysql -u root
    ```

2. Create a new user (we use 420cry-user for this project): In the MySQL shell, run the following SQL command to create the new user with a password:
    ```bash
    CREATE USER '420cry-user'@'localhost' IDENTIFIED BY 'Password';
    ```

3. Grant privileges to the new user: Now, grant the necessary privileges to the new user for the 420cry-db database:
    ```bash
    GRANT ALL PRIVILEGES ON `420cry-db`.* TO '420cry-user'@'localhost';
   ```

4. Flush privileges: Apply the changes to the user privileges:
    ```bash
    FLUSH PRIVILEGES;
   ```

5. Exit MySQL: Exit the MySQL shell:
    ```bash
    EXIT;
   ```

6. Verify the new user:
    ```bash
    mysql -u 420cry-user -p
   ```

7. Create the database: Once you're logged in to the MySQL shell, run the following SQL command to create the database:
    ```bash
   CREATE DATABASE `420cry-db`;
   ```

## Frequently asked questions
### How can I see which application uses a port?
You can easily check this with the command below.
```shell
sudo netstat -tulpn | grep -E "(80|443|3306)"
```

This is very useful if you get an error like
```
ERROR: for dev-server_mysql_1  Cannot start service mysql: Ports are not available: listen tcp 0.0.0.0:3306: bind: address already in use
```
or
```
WARNING: Host is already in use by another container
ERROR: for dev-server_proxy_1  Cannot start service proxy: driver failed program
```

### What should I do if I encounter a port issue?
If you encounter a port issue, you have two options:

1. Stop MySQL locally: If MySQL is running locally on your machine and using the port, you can stop it to free up the port.

2. Update Docker port: You can modify your Docker configuration to use a different port for MySQL.