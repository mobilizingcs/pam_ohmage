#!/bin/bash
while ! nc -z $OHMAGE_SERVER $OHMAGE_SERVER_PORT; do
  echo "Waiting for ohmage to be available "
  sleep 3;
done
/test_setup/node_modules/babel-cli/bin/babel-node.js setup.js