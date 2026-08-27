FROM golang:1.22-alpine
WORKDIR /app
COPY . .
RUN go build -o pusula-serve .
EXPOSE 8080
CMD ["./pusula-serve"]
