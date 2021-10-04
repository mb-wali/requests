FROM golang:1.16-alpine

RUN apk add --no-cache make
RUN apk add --no-cache git
RUN go get -u github.com/jstemmer/go-junit-report

ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/cyverse-de/requests
COPY . .
RUN make

FROM scratch

WORKDIR /app

COPY --from=0 /go/src/github.com/cyverse-de/requests/requests /bin/requests
COPY --from=0 /go/src/github.com/cyverse-de/requests/swagger.json swagger.json

# copy config file 
COPY jobservices.yml /etc/iplant/de/jobservices.yml

ENTRYPOINT ["requests"]

EXPOSE 8080


# build
# docker build -t mbwali/requests:latest .

# run
# docker run -it -p 8080:8080 mbwali/requests:latest

# config
# /etc/iplant/de/jobservices.yml
