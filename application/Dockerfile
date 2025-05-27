FROM golang:1.24.1-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/halftone ./cmd/halftone

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app

COPY --from=build /app/halftone /app/halftone

EXPOSE 8080

ENTRYPOINT ["/app/halftone", "api"]