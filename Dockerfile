FROM golang:1.24 AS build

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ENV CGO_ENABLED=0
RUN go build -o memwaste main.go

FROM busybox:latest AS final

COPY --from=build /app/memwaste /usr/local/bin/memwaste
ENTRYPOINT ["/usr/local/bin/memwaste"]
CMD ["--amount=10M"]
