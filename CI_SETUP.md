# Configurazione CI/CD

I workflow in `.github/workflows` assumono che `code` sia la radice del
repository GitHub e che il branch di rilascio sia `main`.

## Cose da fare su Docker Hub

1. Crea nel tuo namespace i repository `todolistservice` e
   `todolist-frontend` (Docker Hub > **Repositories** > **Create repository**).
2. Apri **Account settings > Personal access tokens > Generate new token**.
   Assegna un nome riconoscibile, una scadenza e il permesso **Repo Read & Write**.
   Copia il token: Docker Hub lo mostra una sola volta. Per un namespace di
   organizzazione, usa un organization access token con permesso di push su
   entrambi i repository.

## Cose da fare su GitHub

1. In **Settings > Secrets and variables > Actions > Variables**, crea
   `DOCKERHUB_USERNAME` con il namespace Docker Hub. Per un account personale
   coincide con lo username.
2. Nella stessa schermata, in **Secrets**, crea `DOCKERHUB_TOKEN` e incolla il
   token Docker Hub. Non inserirlo nei file del repository.
3. Abilita GitHub Actions per il repository. I workflow usano il
   `GITHUB_TOKEN` generato automaticamente da GitHub, con i permessi necessari
   dichiarati nei file YAML. Non serve creare un personal access token GitHub.
4. Usa messaggi di commit in formato Conventional Commits, ad esempio
   `feat: aggiunge filtro` o `fix: corregge validazione`. Semantic Release
   crea tag separati (`todolistservice-v*` e `todolist-frontend-v*`) e le
   immagini ricevono sia il tag di versione sia il commit SHA.
5. Apri una pull request verso `main` per verificare build e test. Dopo il
   merge, controlla **Actions**, **Releases** e i due repository Docker Hub.

I workflow installano automaticamente `semantic-release-monorepo` tramite
`extra_plugins` dell'action. `@semantic-release/commit-analyzer`,
`@semantic-release/release-notes-generator` e `@semantic-release/github` sono
inclusi in Semantic Release. Non serve installare plugin su GitHub o Docker Hub,
né aggiungere una dipendenza npm all'applicazione frontend per la release.

Con il `GITHUB_TOKEN` automatico, una release GitHub non avvia altri workflow
basati sull'evento `release`. Se in futuro servirà quel trigger, crea un token
GitHub dedicato da **Settings > Developer settings > Personal access tokens >
Fine-grained tokens > Generate new token**. Seleziona solo questo repository,
imposta una scadenza e **Contents: Read and write**. Copia il token in
**Settings > Secrets and variables > Actions > Secrets** con il nome
`RELEASE_GH_TOKEN`, poi usa quel secret sia nel checkout sia nello step Semantic
Release. Il token automatico resta la scelta per i workflow attuali.

Riferimenti: [token Docker Hub](https://docs.docker.com/security/access-tokens/personal-access-tokens/),
[token automatico GitHub](https://docs.github.com/en/actions/tutorials/authenticate-with-github_token),
[Semantic Release in monorepo](https://github.com/pmowrer/semantic-release-monorepo).
