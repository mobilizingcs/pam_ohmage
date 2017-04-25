# pam_ohmage
pam_ohmage.so is a Linux Pluggable Authentication Module that authenticates against an Ohmage server. The Mobilize project at UCLA needs students to use RStudio to run the curriculum labs. RStudio Server Pro uses PAM for authentication. This module could be used to provide SSO-like login experience for students. Students can use the same username and password for signing in to RStudio and Ohmage.

## How to build:
#### If you are on a linux system (tested with debian jessie)
- Install libpam headers (in libpam0g-dev): `apt-get install libpam0g-dev`
- `make`
#### If you are not a linux system, you should use the Docker method
- Build the docker image: `docker build . -t mobilizingcs/pam_ohmage`
- Run the docker container: `docker run -v $(pwd):/go/src/pam_ohmage mobilizingcs/pam_ohmage`
- `pam_ohmage.so` should be available at `./bin/pam_ohmage.so`

## Functional testing with RStudio & Ohmage:
Testing the module requires a lot of different software components to play together (ohmage, mysql, rstudio, selenium, webdriverio). Docker Compose is used to orchestrate all these components together.

1. Run: `docker-compose up -d`
2. Tail logs from the test-runner: `docker-compose logs -f test_debian`
3. Once complete, run `docker-compose down -v`

Note: If there have been code changes, you should rebuild the pam_ohmage_build service with:  `docker-compose build pam_ohmage_build` before running `docker-compose up` again.

### Roadmap
  - Add multi-server linux user uid sync support
  - Depend on ohmage username to generate linux uid

### Warning
This module breaks RStudio Server Pro's user impersonation feature when the RStudio user account being impersonated does not already have a corresponding local linux account. This could happen on ephemeral RStudio servers with NFS mounted /home directory.