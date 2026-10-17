FROM golang:1.26.5-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal
RUN CGO_ENABLED=0 go build -trimpath -o /out/api ./cmd/api

FROM alpine:3.23
WORKDIR /app
COPY --from=build /out/api /app/api
USER 10001:10001
ENV HTTP_ADDR=:8000
EXPOSE 8000
CMD ["/app/api"]
