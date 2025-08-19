# Usar imagem oficial do Node
FROM node:20-alpine as build

WORKDIR /app

# Copiar arquivos e instalar dependências
COPY package*.json ./
RUN npm install

# Copiar todo o frontend
COPY . .

# Build da aplicação Angular
RUN npm run build --prod

# Servir com Nginx
FROM nginx:alpine
COPY --from=build /app/dist/frontend /usr/share/nginx/html

# Expor porta
EXPOSE 4200

CMD ["nginx", "-g", "daemon off;"]
