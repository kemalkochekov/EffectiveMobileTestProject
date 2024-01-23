# Effective Mobile Company's Test Project
Developed a service that retrieves full names via an API, enriches the response with the most probable age, gender, and nationality from open APIs, and stores the data in a database. 

## Table of Contents
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Usage](#usage)
- [API Endpoints](#api-endpoints)
- [Linting and Code Quality](#linting-and-code-quality)
  - [Linting Installation](#linting-installation)
  - [Linting Usage](#linting-usage)

## Prerequisites

Before running this application, ensure that you have the following prerequisites installed:

- Go: [Install Go](https://go.dev/doc/install/)
- Docker: [Install Docker](https://docs.docker.com/get-docker/)
- Docker Compose: [Install Docker Compose](https://docs.docker.com/compose/install/)

## Installation

1. Clone this repository
  ```bash
    git clone https://github.com/kemalkochekov/EffectiveMobileTestProject.git
  ```
2. Navigate to the project directory:
  ```
    cd EffectiveMobileTestProject
  ```
## Usage
1. Start the Docker containers:
  ```
    docker-compose up
  ```
2. The application will be accessible at Graphql Playground:
  ```
    localhost:8080
  ```

## API Endpoints
The following API endpoints are available at:

<a href="https://documenter.getpostman.com/view/31073105/2s9YymHjru" target="_blank">
    <img alt="View API Doc Button" src="https://github.com/kemalkochekov/JWT-Backend-Development-App/assets/85355663/0c231cef-ee76-4cdf-bc41-e900845da493" width="200" height="60"/>
</a>


## Linting and Code Quality

This project maintains code quality using `golangci-lint`, a fast and customizable Go linter. `golangci-lint` checks for various issues, ensures code consistency, and enforces best practices, helping maintain a clean and standardized codebase.

### Linting Installation

To install `golangci-lint`, you can use `brew`:

```bash
  brew install golangci-lint
```

### Linting Usage
1. Configuration: 

After installing golangci-lint, create or use a personal configuration file (e.g., .golangci.yml) to define specific linting rules and settings:
```bash
  golangci-lint run --config=.golangci.yml
```
This command initializes linting based on the specified configuration file.

2. Run the linter:

Once configuration is completed, you can execute the following command at the root directory of your project to run golangci-lint:

```bash
  golangci-lint run
```
This command performs linting checks on your entire project and provides a detailed report highlighting any issues or violations found.

3. Customize Linting Rules:

You can customize the linting rules by modifying the `.golangci.yml` file.

For more information on using golangci-lint, refer to the golangci-lint documentation.

