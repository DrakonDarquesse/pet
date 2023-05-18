#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git

# set the working directory
WORKDIR /go/src/app

# copy files from src to dest
COPY . .
RUN go get -d -v ./...
RUN CGO_ENABLED=0 go build -o /go/bin/app

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
EXPOSE 9898 5432
