FROM wernight/kubectl

USER root
RUN apk add --update openssl

COPY ./generate_certificate.sh /app/generate_certificate.sh
RUN chmod +x /app/generate_certificate.sh

WORKDIR /app

ENTRYPOINT "/bin/ash"
CMD ["./generate_certificate.sh"]
