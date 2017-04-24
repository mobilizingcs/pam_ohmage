#!/bin/bash

rstudio() {
  echo "Starting rstudio server..."
  /usr/lib/rstudio-server/bin/rserver --server-daemonize 0
}

addUser() {
  if [ -z "$1" ]
  then
    return 0
  fi
  /usr/sbin/adduser \
    --no-create-home \
    --shell /usr/sbin/nologin \
    --disabled-password \
    --disabled-login \
    -q --gecos "" \
    $1
}

addUserAndHomeDir() {
  if [ -z "$1" ]
  then
    return 0
  fi
  /usr/sbin/adduser \
    --shell /usr/sbin/nologin \
    --disabled-password \
    --disabled-login \
    -q --gecos "" \
    $1
}

rstudio & rstudio_pid=${!}

echo "Setting up user accounts for the tests"
# a: user account & home directory
# scenario a1: user account does not exist, home directory does not exist
# username: uclaids-34573 | canSignIn

# scenario a2: user account exists, home directory does not exist
# username: uclaids-51231 | canSignIn
addUser uclaids-51231

# scenario a3: user account exists, home directory exists
# username: uclaids-53651 | canSignIN
addUserAndHomeDir uclaids-53651

# scenario a4: user account does not exist, home directory exists
# username: uclaids-68912 | canSignIN
# this account creation should be deferred... must be the last account to be created
# see #deferredScenarioA4 comment below for code


# b: home directory ownership
# scenario b1: user account does not exist, home directory is owned by root / uid < 500
# username: uclaids-78945 | canSignIN = false
mkdir /home/uclaids-78945
chown root:root /home/uclaids-78945

# scenario b2: user account does not exist, home directory is owned by uid > 500
# username: uclaids-32478  | canSignIN
mkdir /home/uclaids-32478
chown 2500:2500 /home/uclaids-32478

# scenario b3: user account exists, home directory is owned by root / uid < 500
# username: uclaids-67154 | canSignIN = false
addUserAndHomeDir uclaids-67154
chown 300:300 /home/uclaids-67154

# scenario b4: user account exists, home directory is owned by uid > 500
# username: uclaids-25987 | canSignIN
addUserAndHomeDir uclaids-25987
chown 4525:4525 /home/uclaids-25987

#deferredScenarioA4
addUserAndHomeDir uclaids-68912
# directory is owned by a user that once existed
/usr/sbin/deluser uclaids-68912 > /test_entrypoint.log 2>&1

echo "User account setup complete. Waiting for tests to start and finish."
while true
do
  tail -f /dev/null & wait ${!}
done