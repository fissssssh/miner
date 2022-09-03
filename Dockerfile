FROM golang:alpine AS build

RUN  mkdir /app
COPY . /app
WORKDIR /app
RUN  CGO_ENABLED=0 GOOS=linux go build -o miner ./cmd/

FROM scratch as fianl
COPY --from=build /app/miner .
COPY --from=build /app/index.html .
COPY --from=build /app/config.yaml .
EXPOSE 8080

ENTRYPOINT ["./miner"]