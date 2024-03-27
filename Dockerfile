FROM golang:1.22 as build

ENV CGO_ENABLED 0

WORKDIR /go/src/app
COPY go.* .
RUN go mod download

COPY . .

RUN go build -o /go/bin/aether cmd/exporter/main.go

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/aether /
CMD ["/aether"]
