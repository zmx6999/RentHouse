FROM alpine:3.2
ADD conf /conf
ADD html /html
ADD web /web
ADD crypto-config /crypto-config
WORKDIR /
ENTRYPOINT [ "/web" ]
EXPOSE 8000
