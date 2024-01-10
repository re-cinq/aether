FROM golang:1.21 as build

WORKDIR /go/src/app
COPY go.* .
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -o /go/bin/cloud-carbon .

FROM gcr.io/distroless/static-debian11

COPY --from=build /go/bin/cloud-carbon /
CMD ["/cloud-carbon"]
