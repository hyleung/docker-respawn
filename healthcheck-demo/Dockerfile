FROM alpine

RUN mkdir -p /opt/work

COPY healthcheck.sh /opt/work/healthcheck.sh

WORKDIR /opt/work
HEALTHCHECK --interval=2s CMD ["/bin/sh", "/opt/work/healthcheck.sh"]

CMD ["/bin/sh", "-c",  "while true; do echo hello world; sleep 1; done"]
