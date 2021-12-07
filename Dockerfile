FROM alpine:3.13

RUN adduser -D frozen-throne-user

ADD bin/linux/amd64/frozen-throne /bin/frozen-throne

USER frozen-throne-user 

CMD ["/bin/frozen-throne"]

