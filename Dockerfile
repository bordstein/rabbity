from debian:jessie
VOLUME ["/data"]

# install needed services
RUN apt-get update && \
	apt-get upgrade -y -q && \
	apt-get -q -y install golang git mongodb-server supervisor

# setup GOPATH
ENV PATH /usr/src/go/bin:$PATH
RUN mkdir -p /go/src
ENV GOPATH /go
ENV PATH /go/bin:$PATH

# install go dependencies
RUN go get github.com/gin-gonic/gin
RUN go get gopkg.in/mgo.v2
RUN go get github.com/tonnerre/golang-go.crypto/sha3

# run supervisor with supplied config per default
ADD deployment/supervisord.conf /etc/supervisor/supervisord.conf
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisor/supervisord.conf"]

# Document the external reachable ports
# 8080	rabbity webservice
# 9001	supervisor interface
# 27017	mongodb
# 28017	mongodb status page
EXPOSE 8080 9001 27017 28017

# Copy the local package files to the container's workspace.
ADD ./node /go/src/github.com/bordstein/rabbity/node

# Build the rabbity project inside the container.
RUN go install github.com/bordstein/rabbity/node
