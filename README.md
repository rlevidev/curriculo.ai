# curriculo.ai
This project is all about solving a problem a lot of people face: getting overlooked by companies during job applications.

## Desenvolvimento local

O ambiente local executa o frontend Vite e o backend Go separadamente:

```text
http://localhost:5173 (React/Vite) -> http://localhost:8080 (Go/TeX Live)
```

### Pré-requisitos

- Node.js 22 ou compatível com o workflow de CI
- npm
- Docker com Docker Compose

### Iniciar o backend

Na raiz do projeto:

```bash
docker compose up --build backend
```

O backend ficará disponível em `http://localhost:8080`. A opção `--build` é necessária na primeira execução e após alterações no backend ou no Dockerfile.

Também é possível executá-lo sem Docker, desde que o TeX Live esteja instalado localmente:

```bash
cd backend
FRONTEND_ORIGIN=http://localhost:5173 go run .
```

### Iniciar o frontend

Em outro terminal:

```bash
cd frontend
npm ci
npm run dev
```

Acesse `http://localhost:5173` no navegador. O ambiente de desenvolvimento usa `frontend/.env.development`, que aponta as chamadas para `http://localhost:8080`.

### Verificar o backend

```bash
curl http://localhost:8080/health
```

A resposta esperada é `200 OK`. Com os dois serviços ativos, preencha um currículo e use **Exportar PDF** para testar o fluxo completo localmente.

### Produção x desenvolvimento

- Desenvolvimento: `http://localhost:5173` chama `http://localhost:8080`.
- Produção: o frontend publicado no GitHub Pages chama `https://curriculo-ai.onrender.com`.
- O caminho `/curriculo.ai/` do GitHub Pages é usado apenas no build de produção; localmente o frontend é servido em `/`.
- `FRONTEND_ORIGIN` controla a origem permitida pelo CORS do backend e é configurado pelo Docker Compose no ambiente local.

### Verificações

```bash
cd backend
go vet ./...
go test -race ./...

cd ../frontend
npm run lint
npx tsc -b --noEmit
npm test
npm run build
```
