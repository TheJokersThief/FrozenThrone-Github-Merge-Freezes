FROM alpine:3.13

RUN adduser -D frozen-throne-user

ADD bin/linux/frozen_throne /bin/frozen-throne

USER frozen-throne-user 

CMD ["/bin/frozen-throne"]

