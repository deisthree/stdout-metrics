FROM alpine:3.3

# add logger user
RUN addgroup -S nginxpusher && \
	adduser -S -G nginxpusher -H -h /opt/nginxpusher -D nginxpusher

ADD . /

# Fix some permission since we'll be running as a non-root user.
RUN chown -R nginxpusher:nginxpusher /opt/nginxpusher

USER nginxpusher

CMD /opt/nginxpusher/sbin/nginxpusher &
EXPOSE 1514 8088

ENV DEIS_RELEASE v2-beta
