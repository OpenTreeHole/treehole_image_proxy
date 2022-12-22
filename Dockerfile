FROM golang:1.19-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN apk add --no-cache --virtual .build-deps \
        ca-certificates \
        gcc \
        g++ &&  \
    go mod download

COPY . .

RUN go build -tags="release" -ldflags "-s -w" -o image_proxy

FROM alpine

WORKDIR /app

COPY --from=builder /app/treehole /app/

EXPOSE 8000

ENTRYPOINT ["./image_proxy"]