# syntax=docker/dockerfile:1
FROM golang:1.24.1-bookworm AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/halftone ./cmd/halftone

RUN --mount=type=secret,id=firebase_service_account,target=/tmp/firebase-service-account.json \
    cp /tmp/firebase-service-account.json /app/firebase-service-account.json

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /app

COPY --from=build /app/halftone /app/halftone
COPY --from=build /app/firebase-service-account.json /app/firebase-service-account.json

EXPOSE 8080

ADD https://truststore.pki.rds.amazonaws.com/global/global-bundle.pem /opt/global-bundle.pem

ENTRYPOINT ["/app/halftone", "api"]