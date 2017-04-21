.build:
	go build -buildmode=c-shared -o ./bin/pam_ohmage.so

.test:
	sudo /usr/lib/rstudio-server/bin/pamtester --verbose rstudio vagrant authenticate

make: .build .test

build: .build