How to build:
  - If you are on a linux machine (tested with debian jessie)
    - Install libpam headers (in libpam0g-dev): `apt-get install libpam0g-dev`
    - `make`
  - If you are not a linux machine, you should use the Docker method
    - Make sure docker is available
    - Build the docker image: `docker build . -t mobilizingcs/pam_ohmage`
    - Run the docker container to build pam_ohmage: `docker run -v $(pwd):/go/src/pam_ohmage mobilizingcs/pam_ohmage`
    - `pam_ohmage.so` should be available at `./bin/pam_ohmage.so`

How to build for data container:
  - Run : `docker-compose build pam_ohmage_build`

How to do functional testing w/ RStudio:

  - Using Docker Compose
    - Run: `docker-compose up -d`
    - If there have been code changes, you should rebuild the pam_ohmage_build service with:
        `docker-compose build pam_ohmage_build`
    - And re-run all services with: `docker-compose up -d`

Roadmap
  - Add multi-server uid sync support
    - Depend on ohmage username to generate linux uid

Warning
  - This module breaks RStudio Server Pro's user impersonation feature when the
RStudio user account being impersonated does not already have a corresponding local
linux account. This could happen on ephemeral RStudio servers with NFS mounted /home
directory.