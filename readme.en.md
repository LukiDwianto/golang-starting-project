# Golang Starter Project

This is a starter project for a Golang application, utilizing modern technologies for development and deployment.

## Technologies Used

- **Gin**: A minimal and fast web framework for Golang.
- **PostgreSQL**: A reliable and feature-rich relational database.
- **Docker Compose**: Manages and orchestrates services in Docker containers.
- **Air**: A live-reloading tool for Golang development.
- **Makefile**: Automates common tasks such as building, running, and fetching logs.

## Project Structure

```plaintext
.
|-- .data                 # Data volume for PostgreSQL
|-- constants             # Folder for variable constants
|-- controllers           # Controllers for routes
|-- helpers               # General utilities
|-- middleware            # Middleware for routes
|-- models                # Database model definitions
|-- router                # API route definitions
├-- air.toml              # Configuration for Air (live-reloading)
├-- docker-compose.yml    # Configuration for Docker Compose
├-- Dockerfile.dev        # Development Dockerfile for Golang application
├-- Dockerfile.prod       # Production Dockerfile for Golang application
├── go.mod                # Dependency manager for Golang
├── go.sum                # Dependency checksums
├── main.go               # Application entry point
└-- Makefile              # Makefile to automate tasks
```

## How to Use

### Prerequisites

- **Docker** and **Docker Compose** must be installed on your system.
- **Golang** version 1.17 or later installed on your system.
- **Make** installed on your system to run `Makefile` commands.

### Steps

1. **Clone the Repository**

   ```bash
   git clone <repository-url>
   cd <project-directory>
   ```

2. **Set Up the Environment**

   - Create a `.env` file based on the `.env.example` file and adjust the configuration as needed.

3. **Running the Application**

   You can run the application using the `Makefile`:

   - **up_build**: Stops any running Docker containers, builds, and starts the Docker containers in detached mode.

     ```bash
     make up_build
     ```

   - **logs**: Fetches logs from all services.

     ```bash
     make logs
     ```

   - **init**: Builds and starts the Docker containers in detached mode.

     ```bash
     make init
     ```

4. **Development with Live Reloading**

   The application is configured to use Air for live-reloading during development. Run the following command:

   ```bash
   air
   ```

5. **Docker Setup**

   - **docker-compose.yml**: Manages two main services:
     - `db`: PostgreSQL database service running PostgreSQL 13.
     - `web`: Golang application service built from `Dockerfile.dev` for development.
   - **Dockerfile.dev**: Used for development with Air.
   - **Dockerfile.prod**: Used to build the application for production.

## Contribution

Contributions are welcome! Please create a pull request or report an issue.
