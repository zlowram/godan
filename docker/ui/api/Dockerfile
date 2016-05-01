FROM alpine:latest

MAINTAINER Sergi Martinez <zlowram@gmail.com>

ADD godanui.toml /etc/
ADD godanui start-godanui.sh /usr/local/bin/
ADD keys /keys
RUN chmod 755 /usr/local/bin/godanui /usr/local/bin/start-godanui.sh

ENTRYPOINT ["/usr/local/bin/start-godanui.sh"]

EXPOSE 8000
