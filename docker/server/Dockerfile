FROM alpine:latest

MAINTAINER Sergi Martinez <zlowram@gmail.com>

ADD godanserver.toml /etc/
ADD godanserver start-godanserver.sh /usr/local/bin/
RUN chmod 755 /usr/local/bin/godanserver /usr/local/bin/start-godanserver.sh

ENTRYPOINT ["/usr/local/bin/start-godanserver.sh"]

EXPOSE 8080
