FROM golang:1.16-alpine AS build

WORKDIR /src/
COPY main.go go.* /src/
RUN CGO_ENABLED=0 go build -o /bin/conn-test

FROM scratch
COPY --from=build /bin/conn-test /bin/conn-test
ENTRYPOINT ["/bin/conn-test"]