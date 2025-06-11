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
    127.0.0.1 db.420.crypto.test
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
This command applies gofumpt and goimports to fix formatting and organize imports.
```bash
make lint-fix
```
### Test:
```bash
make test
```
## ⚠️ Ensure Go Tools Are in Your PATH (MAC OS)

If you encounter a "command not found" error when running `make lint-fix`, it's likely because your Go-installed binaries are not in your system `PATH`.

Add this to your shell profile (`~/.zshrc`, `~/.bashrc`, or `~/.bash_profile`):

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```
Then reload your shell:
```bash
source ~/.zshrc  # or source ~/.bashrc
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

## MailHog
You can access MailHog at 
```bash
    http://localhost:8025/#
```

## Project Structure

This project follows a Domain-Driven Design (DDD) approach, with a well-defined folder structure to separate concerns and ensure clarity in the codebase. Below is a breakdown of the project structure:

### 1. **api**  
Contains the API-related components, including routes, controllers, and any logic related to the HTTP API. This folder is responsible for exposing the domain logic through the server interface.

### 2. **services**  
Contains services and orchestrates the interaction between the domain layer and the external world. This is where business logic is executed, like sending emails or processing user actions.

### 3. **core**  
Contains the foundational code of the services, including utilities, helpers, and common services that are used across other parts of the services.

### 4. **domain**  
The domain layer represents the heart of the business logic and contains entities, value objects, aggregates, and domain services. This is where the core business rules and logic reside.

### 5. **server**  
This folder contains the configuration and setup for running the server, including setting up the routes, middleware, and server initialization.

### 6. **database**  
Contains database-related code, including migrations, schema definitions, and database models. It is responsible for managing data persistence.

### 7. **migration**  
Contains database migration files, which are used to manage changes to the database schema over time.

### 8. **templates**  
Contains the HTML or email templates used in the services, like the verification email templates.

### 9. **types**  
Contains type definitions.

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