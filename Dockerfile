FROM golang:latest

WORKDIR /app


# アプリケーションのソースをコピーする
COPY ./src .
