FROM alpine:latest

MAINTAINER Sergi Martinez <zlowram@gmail.com>

ADD godanworker.toml /etc/
ADD godanworker start-godanworker.sh /usr/local/bin/
RUN chmod 755 /usr/local/bin/godanworker /usr/local/bin/start-godanworker.sh

ENTRYPOINT ["/usr/local/bin/start-godanworker.sh"]
