FROM alpine:3.3

# add logger user
RUN addgroup -S stdout-metrics && \
	adduser -S -G stdout-metrics -H -h /opt/stdout-metrics -D stdout-metrics

ADD . /

# Fix some permission since we'll be running as a non-root user.
RUN chown -R stdout-metrics:stdout-metrics /opt/stdout-metrics
USER stdout-metrics
EXPOSE 1514 8088
CMD /opt/stdout-metrics/sbin/stdout-metrics
