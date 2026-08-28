FROM golang:1.22-alpine AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -o /pusula-serve ./cmd/pusula-serve

FROM alpine:3.20
COPY --from=build /pusula-serve /usr/local/bin/pusula-serve
EXPOSE 8080
USER nobody
ENTRYPOINT ["pusula-serve"]
