FROM golang as build

WORKDIR /app
COPY . .
WORKDIR contrib/renderpdf
RUN go build -o /app/main main.go

FROM ubuntu as runtime

RUN apt-get update
RUN apt-get install -y -q poppler-utils

COPY --from=build /app/main /usr/bin/main

USER 1000:1000
ENTRYPOINT ["/usr/bin/main"]
