# Usar imagem oficial do Go
FROM golang:1.25-alpine

# Diretório de trabalho no container
WORKDIR /app

# Copiar os arquivos go.mod e go.sum primeiro (pra cache do Docker)
COPY go.mod ./
COPY go.sum ./

# Baixar dependências
RUN go mod download

# Copiar todo o código do backend
COPY . .

# Compilar o binário
RUN go build -o backend .

# Expor porta
EXPOSE 8080

# Comando para rodar a aplicação
CMD ["./backend"]
