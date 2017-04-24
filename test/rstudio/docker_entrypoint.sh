#!/bin/bash
while ! nc -z $OHMAGE_SERVER $OHMAGE_SERVER_PORT; do
  echo "Waiting for ohmage to be available "
  sleep 3;
done

while ! nc -z $RSTUDIO_SERVER $RSTUDIO_SERVER_PORT; do
  echo "Waiting for rstudio to be available "
  sleep 3;
done

while ! nc -z $SELENIUM_SERVER $SELENIUM_SERVER_PORT; do
  echo "Waiting for selenium to be available "
  sleep 3;
done

set -e
node runner.js