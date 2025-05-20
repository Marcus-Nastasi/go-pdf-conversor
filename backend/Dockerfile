# 1: build
FROM golang:1.24.2 AS builder

# Instala git e outras dependências básicas
RUN apk update && apk add --no-cache git

# Define workdir
WORKDIR /app

# Copia arquivos e baixa dependências
COPY go.mod ./
COPY go.sum ./
ENV GOPROXY=direct
RUN go mod download

COPY . .

# Compila o binário
RUN go build -o main .

# 2: final image
FROM ubuntu:latest

WORKDIR /root/

# Copia binário final
COPY --from=builder /app/main .

EXPOSE 8081

# Comando de execução
CMD ["./main"]
