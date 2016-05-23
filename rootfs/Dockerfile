FROM quay.io/deis/base:0.2.0

RUN adduser --system \
	--shell /bin/bash \
	--disabled-password \
	--home /opt/stdout-metrics \
	--group \
	stdout-metrics

COPY . /

# Fix some permission since we'll be running as a non-root user.
RUN chown -R stdout-metrics:stdout-metrics /opt/stdout-metrics

USER stdout-metrics

CMD /opt/stdout-metrics/sbin/stdout-metrics
EXPOSE 1514 8088

ENV WORKFLOW_RELEASE 2.0.0
