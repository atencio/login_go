FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go mod tidy && CGO_ENABLED=0 go build -o /api ./cmd/api

FROM alpine:3.20
COPY --from=build /api /api
EXPOSE 8080
CMD ["/api"]
