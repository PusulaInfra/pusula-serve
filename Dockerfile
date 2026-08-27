FROM golang:1.22-alpine AS build
WORKDIR /src
COPY . .
RUN go build -o /pusula-serve .

FROM alpine:3.20
WORKDIR /app
COPY --from=build /pusula-serve /usr/local/bin/pusula-serve
EXPOSE 8080
USER nobody
CMD ["pusula-serve"]
