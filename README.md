How to build:
  - If you are on a linux machine (tested with debian jessie)
    - Install libpam headers (in libpam0g-dev): `apt-get install libpam0g-dev`
    - `make`
  - If you are not a linux machine, you should use the Docker method
    - Make sure docker is available
    - Run `docker run -v $(pwd):/go/src/pam_ohmage mobilizingcs/pam_ohmage`
    - `pam_ohmage.so` should be available at `./bin/pam_ohmage.so`

How to do functional testing w/ RStudio:
  - Using docker
    - Make sure docker is available
    - Build pam_ohmage.so: `docker run -v $(pwd):/go/src/pam_ohmage mobilizingcs/pam_ohmage`
    - Test on debian
      - Build the container: `docker build -f ./test/debian/Dockerfile -t pam_ohmage/debian .`
      - Run the container: `docker run -v $(pwd)/bin:/pam_ohmage/bin -p 8787:8787 pam_ohmage/debian`

Roadmap
  - Add multi-server uid sync support
    - Depend on ohmage username to generate linux uid

Warning
  - This module breaks RStudio Server Pro's user impersonation feature when the
RStudio user account being impersonated does not already have a corresponding local
linux account