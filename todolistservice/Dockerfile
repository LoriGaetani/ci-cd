FROM golang:1.27-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /todolistservice ./cmd/todolistservice

FROM gcr.io/distroless/static-debian13:nonroot
COPY --from=build /todolistservice /todolistservice
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/todolistservice"]

