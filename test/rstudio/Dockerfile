FROM node:7.9.0-slim

RUN apt-get update \
    && apt-get -y install netcat python g++ build-essential

# Context must be [PROJECT ROOT]
COPY ./test/rstudio /tests
COPY ./test/setup_config.json /tests/setup_config.json

RUN chmod +x /tests/docker_entrypoint.sh
WORKDIR /tests
RUN npm install

ENTRYPOINT /tests/docker_entrypoint.sh