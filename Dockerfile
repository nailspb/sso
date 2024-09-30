FROM golang:latest AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o ./sso ./cmd/sso/main.go

#FROM ubuntu:20.04 AS ubuntu
#RUN apt-get update
#RUN DEBIAN_FRONTEND=noninteractive apt-get install -y upx
#RUN apt-get clean && \
#    rm -rf /var/lib/apt/lists/*
#COPY --from=build /build/sso /build/sso
#RUN upx /build/sso


FROM scratch
COPY ./build/migrations/ /app/migrations/
COPY --from=build /build/sso /app/sso
WORKDIR /app
CMD ["/app/sso"]