FROM apache/apisix:latest

USER root
RUN apt install -y httpie

COPY --chown=apisix:apisix services/apisix/config.yaml /usr/local/apisix/conf/config.yaml
