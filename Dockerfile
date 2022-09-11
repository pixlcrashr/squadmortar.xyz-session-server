FROM golang:alpine as build

WORKDIR /app/

COPY . .

RUN go mod download
RUN go build -o smss main.go

FROM alpine:3.9

ARG USER=smss
ARG GROUP=smss

RUN addgroup -S ${GROUP} && adduser -S ${USER} -G ${GROUP}

USER ${USER}:${GROUP}

WORKDIR /opt/smss

COPY --from=build --chown=${USER}:${GROUP} /app/smss /opt/smss/bin/

ENTRYPOINT [ "/opt/smss/bin/smss" ]
