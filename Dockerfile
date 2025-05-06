# syntax=docker/dockerfile:1

FROM golang:1.24.2 AS builder

WORKDIR /code

COPY . .
RUN CGO_ENABLED=0 go build -o build/ek main.go


FROM scratch

COPY --from=builder /code/build/ek /

EXPOSE 8080

ENTRYPOINT [ "/ek" ]
