FROM golang as build

WORKDIR /app
COPY . .
WORKDIR contrib/zbarimg
RUN go build -o /app/main main.go

FROM ubuntu as runtime

RUN apt-get update
RUN apt-get install -y -q zbar-tools
RUN apt-get install -y -q dbus

COPY --from=build /app/main /usr/bin/main

USER 1000:1000
ENTRYPOINT ["/usr/bin/main"]
