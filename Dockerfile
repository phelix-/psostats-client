FROM golang:1.16-alpine as build-env

RUN mkdir /build
WORKDIR /build

COPY server /build/server
COPY pkg /build/pkg
COPY go.mod /build/go.mod
COPY go.sum /build/go.sum

RUN go build -o psostats_server /build/server/cmd/main.go

FROM golang:1.16-alpine as prod-env
RUN mkdir /app/
COPY --from=build-env /build/psostats_server /app/
COPY ./server/internal/templates /app/server/internal/templates
COPY ./static /app/static
WORKDIR /app
CMD [ "/app/psostats_server" ]