FROM golang:1.10 AS BUILD

# RUN go get -v github.com/Soulou/curl-unix-socket

#doing dependency build separated from source build optimizes time for developer, but is not required
#install external dependencies first
# ADD go-plugins-helpers/Gopkg.toml $GOPATH/src/go-plugins-helpers/
ADD /main.dep $GOPATH/src/gittagger/main.go
RUN go get -v gittagger

#now build source code
ADD gittagger $GOPATH/src/gittagger
RUN go get -v gittagger
# RUN go test -v gittager


FROM golang:1.10

ENV LOG_LEVEL ''
ENV GIT_REPO_URL ''
ENV GIT_USERNAME ''
ENV GIT_EMAIL ''

COPY --from=BUILD /go/bin/* /bin/
ADD startup.sh /
EXPOSE 50000

CMD [ "/startup.sh" ]
