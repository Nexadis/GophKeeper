FROM golang:1.21.2-alpine3.18 AS builder

LABEL stage=gobuilder
ENV CGO_ENABLED 0
ENV GOOS linux
RUN apk update --no-cache && apk add --no-cache tzdata

WORKDIR /build
ADD go.mod .
ADD go.sum .
RUN go mod download

COPY . .
RUN go build -ldflags="-s -w" -o /app/server /build/cmd/server/main.go
RUN go build -ldflags="-s -w" -o /app/clients/linux-client /build/cmd/client/main.go

ENV GOOS windows
RUN go build -ldflags="-s -w" -o /app/clients/windows-client /build/cmd/client/main.go

ENV GOARCH arm
ENV GOOS linux
RUN go build -ldflags="-s -w" -o /app/clients/arm-linux-client /build/cmd/client/main.go

FROM alpine:3.18 as runner

RUN apk update --no-cache && apk add --no-cache ca-certificates
WORKDIR /app

COPY --from=builder /usr/share/zoneinfo/America/New_York /usr/share/zoneinfo/America/New_York

ENV TZ America/New_York


COPY --from=builder /app/server /app/server

CMD ["/app/server"]
