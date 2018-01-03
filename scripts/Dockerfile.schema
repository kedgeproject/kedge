FROM fedora:27

RUN dnf install nodejs -y

WORKDIR /data

CMD ["/bin/sh","scripts/update-schema.sh"]
