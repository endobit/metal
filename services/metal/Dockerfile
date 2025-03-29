FROM golang:1.24rc3-alpine AS builder
WORKDIR /src

RUN apk add figlet git make

RUN figlet Stack > banner.txt

COPY go.mod go.sum .
RUN  go mod download

COPY . .
RUN --mount=type=cache,target=/root/.cache/go-build make build


FROM alpine:latest
WORKDIR /

EXPOSE 8080

COPY --from=builder /src/stackd /bin/stackd
COPY --from=builder /src/banner.txt .
COPY --from=builder /src/services/stack/start.sh .

CMD [ "/bin/sh", "start.sh" ]
