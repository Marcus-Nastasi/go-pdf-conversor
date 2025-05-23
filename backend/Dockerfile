# 1: build
FROM golang:1.24.2 AS builder

# Define workdir
WORKDIR /app

# Copia arquivos e baixa dependências
COPY go.mod ./
COPY go.sum ./
ENV GOPROXY=direct
RUN go mod tidy

COPY . .

# Compila o binário
RUN go build -o main .

# 2: final image
FROM ubuntu:latest

# Install lowriter package
RUN apt update && apt install -y libreoffice-writer

WORKDIR /root/

# Copia binário final
COPY --from=builder /app/main .

EXPOSE 8081

# Comando de execução
CMD ["./main"]
