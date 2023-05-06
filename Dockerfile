FROM golang:1.20
WORKDIR /work
COPY . .
RUN go build
CMD ["/work/nowcast-bot"]