FROM node:7.9.0-slim

RUN apt-get update \
    && apt-get -y install netcat

# Context must be [PROJECT ROOT]
ADD ./test/ohmage_test_setup /test_setup
ADD ./test/setup_config.json /test_setup/setup_config.json

RUN chmod +x /test_setup/docker_entrypoint.sh
WORKDIR /test_setup
RUN npm install

ENTRYPOINT /test_setup/docker_entrypoint.sh