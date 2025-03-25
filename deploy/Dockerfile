# syntax=docker/dockerfile:1

FROM golang:1.21

WORKDIR /app

COPY . .
RUN go mod download

WORKDIR /app/cmd
RUN GOOS=linux go build -o ugc

EXPOSE 8080

CMD ["./ugc"]

# docker build -f .\deploy\Dockerfile --progress=plain -t dmitriytitov/ugc:latest .
# docker run --rm -p 8080:8080 dmitriytitov/ugc:latest