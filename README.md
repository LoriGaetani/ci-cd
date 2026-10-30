
# Mini progetto Todo per esercizi CI/CD

Il progetto contiene un frontend React/Vite in `code/frontend`, un'API Go in
`code/todolistservice` e PostgreSQL in `code/compose.yaml`. La pipeline GitHub Actions,
le immagini pubblicate e i manifest Argo CD restano esercizi da costruire a partire
da questa applicazione.

## Avvio locale

Prerequisiti: Node.js 20.19+ o 22.12+, Go 1.27 e Docker Compose.

1. Da `code`, avvia PostgreSQL: `docker compose up -d postgres`.
2. Da `code/todolistservice`, imposta
   `$env:DB_URL='postgres://todolist:local-dev-password@localhost:5432/todolist?sslmode=disable'`
   in PowerShell e avvia `go run ./cmd/todolistservice`.
3. Da `code/frontend`, esegui `npm ci` e `npm run dev`.
4. Apri `http://localhost:5173` e aggiungi qualche elemento. I dati restano nel volume
   Docker `postgres_data`.

L'API espone `GET /api/items` e `POST /api/items` con body JSON
`{"description":"Comprare il pane"}`. Le probe sono disponibili su
`/health/live` e `/health/ready`.

## Configurazione per gli ambienti

| Componente | Variabile | Valore locale / significato |
| --- | --- | --- |
| Backend | `DB_URL` | Obbligatoria; URL PostgreSQL completo, da un Secret nel deploy |
| Backend | `SERVER_PORT` | `8080` |
| Frontend in Vite | `VITE_API_PROXY_TARGET` | `http://localhost:8080` |
| Frontend in container | `API_BASE_URL` | `/api` |
| Compose | `POSTGRES_PASSWORD` | `local-dev-password`, solo per sviluppo locale |
| Compose | `POSTGRES_PORT` | `5432` |

Il browser usa `/api` sullo stesso host del frontend. Vite inoltra le richieste al
backend durante lo sviluppo. Nel deploy configura l'Ingress per instradare `/api`
al backend e `/` al frontend. Se serve un percorso diverso, `API_BASE_URL` viene
scritto in `config.js` all'avvio del container frontend: non occorre ricostruire
l'immagine. `config.js` ha `Cache-Control: no-store`. Mantieni preferibilmente la
stessa origine per frontend e API.

I Dockerfile di entrambe le applicazioni sono nelle rispettive cartelle. Per esempio,
da `code`: `docker build -t todolist-frontend:local frontend` e
`docker build -t todolistservice:local todolistservice`.

Nel backend, `cmd/todolistservice` avvia il server; `internal/controller` gestisce
HTTP e validazione del JSON, `internal/service` contiene le regole applicative,
`internal/repository` gestisce PostgreSQL e lo schema iniziale, `internal/model`
contiene il modello condiviso. Lo schema usa `CREATE TABLE IF NOT EXISTS` per
questo prototipo; per evoluzioni successive introduci migrazioni versionate.

Per esercitarti, puoi creare una workflow che esegue `npm ci`/`npm run build` e
`go test ./...`, costruisce le due immagini, poi aggiorna i tag nei manifest GitOps
per i tuoi ambienti Argo CD. In produzione configura `DB_URL` tramite Secret
per ciascun ambiente e usa un PostgreSQL gestito separatamente.
