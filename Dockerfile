# Start by building the application.
FROM golang:1.13 as build

WORKDIR /go/src/github.com/philips/gce-shell
COPY . .

ENV GO111MODULE=on
ENV GOFLAGS=-mod=vendor

RUN go install -v ./

# Now copy it into our base image.
FROM gcr.io/distroless/base
COPY --from=build /go/bin/gce-shell /
CMD ["/gce-shell", "server"]
