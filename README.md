# EasyTrip Backend API

Este é o backend da API EasyTrip, uma aplicação para organização colaborativa de viagens em grupo. A API foi desenvolvida em Go e utiliza o PostgreSQL como banco de dados.

O projeto segue uma arquitetura em camadas (**handlers**, **services**, **repositories**) para separar as responsabilidades e garantir uma base de código limpa e modular.

---

## 🚀 Como Rodar a Aplicação

Este guia irá te ajudar a configurar e executar a API, independentemente do seu sistema operacional. A aplicação não utiliza Docker, rodando diretamente na sua máquina.

### Pré-requisitos

Certifique-se de que as seguintes ferramentas estão instaladas em seu sistema:

- **Go (versão 1.21+):** [Download e Instalação](https://go.dev/doc/install)
- **PostgreSQL:** [Download e Instalação](https://www.postgresql.org/download/)

### 🛠️ Configuração do Projeto

Siga os passos abaixo para configurar o ambiente de desenvolvimento.

#### Passo 1: Clone o Repositório

Clone o repositório para o seu ambiente local:

```bash
git clone <url-do-seu-repositorio>
cd <nome-do-seu-projeto>
```

#### Passo 2: Instale e Configure o PostgreSQL

**Para usuários de Linux (Pop!_OS, Ubuntu):**

Abra o terminal e instale os pacotes:

```bash
sudo apt update
sudo apt install postgresql
```

Crie o usuário e o banco de dados. Acesse o shell do PostgreSQL:

```bash
sudo -i -u postgres
psql
```

Execute os comandos para criar o usuário e o banco.  
**Lembre-se de substituir `sua_senha_aqui` pela sua senha.**

```sql
CREATE USER admin WITH PASSWORD 'sua_senha_aqui';
CREATE DATABASE project_lab OWNER admin;
\q
```

**Para usuários de Windows:**

1. Baixe e instale o PostgreSQL a partir do site oficial.
2. Durante a instalação, o instalador irá pedir para você criar uma senha para o usuário padrão `postgres`. Anote essa senha.
3. Abra a ferramenta **pgAdmin** (que vem com a instalação) e crie um novo banco de dados chamado `project_lab`.
4. Crie um novo usuário chamado `admin` e defina a senha que você irá usar no arquivo `.env`.

---

#### Passo 3: Configure as Variáveis de Ambiente

Crie um arquivo na raiz do projeto chamado `.env` para armazenar suas credenciais e configurações de forma segura.

```ini
DB_HOST=localhost
DB_PORT=5432
DB_USER=admin
DB_PASSWORD=sua_senha_aqui
DB_NAME=project_lab
```

⚠️ Importante: O arquivo `.env` já está no `.gitignore` para que suas credenciais não sejam enviadas para o Git.

---

#### Passo 4: Baixe as Dependências do Go

Na raiz do seu projeto, no terminal, execute o comando para instalar todas as dependências:

```bash
go mod tidy
```

---

### ▶️ Executando a Aplicação

Com todas as configurações prontas, você pode iniciar o servidor da API com um único comando no terminal:

```bash
go run .
```

A aplicação irá se conectar ao banco de dados e criar as tabelas automaticamente na primeira execução, de acordo com o `schema.go`.  
Você verá a mensagem:

