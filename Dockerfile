FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR $GOPATH/src/shopee-mania

COPY . .

RUN go get -d -v ./...

RUN go install -v ./...

EXPOSE 6001
# RUN go build -o main .

# CMD ["mkdir -p test"] 
RUN mkdir -p result

# COPY /static $GOPATH/src/shopee-mania/static

CMD ["shopee-mania"]
# ENTRYPOINT ["/app/main"]