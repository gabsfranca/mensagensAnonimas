# Primeira etapa: compilação
FROM golang:1.23.7-alpine AS builder

# Diretório de trabalho
WORKDIR /app

# Copiar arquivos de dependências
COPY backend/go.mod backend/go.sum ./

# Baixar dependências
RUN go mod download && go mod tidy

# Copiar código-fonte
COPY backend/ ./

# Compilar o aplicativo
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Segunda etapa: imagem mínima para execução
FROM alpine:latest

# Instalar certificados para HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar o binário compilado da etapa anterior
COPY --from=builder /app/main .

# Copiar os arquivos frontend necessários
COPY frontend/ ./frontend/

# Criar .env vazio para evitar erro de carregamento
RUN touch .env

# Expor a porta configurada
EXPOSE 8080

# Definir a variável de ambiente para o caminho do template
ENV TEMPLATE_PATH="./frontend/index.html"

# Comando para executar o aplicativo
CMD ["./main"]