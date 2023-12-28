FROM golang:1.21
LABEL authors="dasu"

WORKDIR /app
COPY go.mod  ./
RUN go mod download

COPY config.yml ./
COPY *.go ./

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-source

ENTRYPOINT ["/docker-source"]