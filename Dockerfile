# syntax=docker/dockerfile:1
FROM golang:alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build


FROM alpine AS runtime

RUN touch /etc/pypihub.toml

COPY --from=build /app/pypihub /

CMD [ "/pypihub", "-c", "/etc/pypihub.toml" ]
