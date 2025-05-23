# PDF Converter

## About the project

This is a Golang API designed to convert various file types to pdf.

## Technologies

### Back-end
- **Language**: Golang 1.24.2
- **Framework**: Gin v1.10.0

## How to run
Follow the steps below to set up and run the project on your local machine.

## Prerequisites
- Git
- Golang 1.24.2
- Docker and Docker Compose

## Steps
**Make sure you have opened the ports 8081 (application) on your machine locally**

1. **Clone this repository.**
   ```
   https://github.com/Marcus-Nastasi/go-pdf-conversor.git
   
2. **Run the app, or run the docker-compose.yml file on "docker" directory:**
    ```bash
    go run main.go

    # or

    [sudo] docker-compose up --build -d

3. **Then, you can access the login endpoint, and the swagger UI interface:**
    ```bash
    # get:
    http://localhost:8081/status

    # get on "/" to access the html
    http://localhost:8081/

    # post with file:
    http://localhost:8081/convert
