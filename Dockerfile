FROM golang:1.26-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o /bin/dotcom .
RUN /bin/dotcom generate

FROM alpine:3.22

WORKDIR /app

COPY --from=build /bin/dotcom /usr/local/bin/dotcom
COPY --from=build /app/public ./public

ENV ADDR=:8080
ENV OUTPUT_DIR=public

EXPOSE 8080

CMD ["dotcom", "serve"]
