# build stage
FROM golang:1.16-alpine AS build-stage

WORKDIR /build
COPY . /build/

RUN go mod download

RUN CGO_ENABLED=0 go build -o goapp .

# final stage
FROM busybox:1.32.0

COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

WORKDIR /app
COPY --from=build-stage /build/goapp ./

RUN addgroup -S appgroup && adduser -S --no-create-home appuser -G appgroup
USER appuser

CMD ["./goapp"]
