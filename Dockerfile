FROM alpine:3.8

WORKDIR "/opt"

COPY release/cev /opt/

COPY templates /opt/templates/
COPY static /opt/static/

RUN chmod -R 777 /opt


WORKDIR /opt/

CMD ["./cev"]