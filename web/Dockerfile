FROM alpine:3.2
ADD html /html
ADD web /web
ADD conf /conf
ADD crypto-config /crypto-config
WORKDIR /
ENTRYPOINT [ "/web" ]
EXPOSE 8000
