FROM golang:1.25-bookworm AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o server

FROM gcr.io/distroless/static-debian12:latest

WORKDIR /

COPY --from=build /app/server /server 

EXPOSE 2500

CMD ["/server", "webhookReceiver"]