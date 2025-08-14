ARG securityBase=security-base-image


FROM $securityBase AS security

FROM golang AS build
COPY --from=security /etc/ssl/certs/*WK* /etc/ssl/certs/

COPY . /app/

WORKDIR /app

RUN go mod tidy
RUN go build

FROM golang AS runtime
COPY --from=security /etc/ssl/certs/*WK* /etc/ssl/certs/
COPY --from=hairyhenderson/gomplate:stable /gomplate /bin/gomplate

WORKDIR /app

RUN apt-get update; apt-get install -y protobuf-compiler

RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
COPY --from=build /app/protoc-gen-go-gin /go/bin/

EXPOSE 8080

COPY docker-contents ./

ENTRYPOINT ["/app/generate.sh"]
