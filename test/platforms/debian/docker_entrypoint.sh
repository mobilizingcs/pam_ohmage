#!/bin/bash
source /common.sh

startRstudio

echo "Starting rsyslog"
touch /var/log/pam_ohmage.log
rsyslogd &

echo "Tailing /var/log/pam_ohmage.log"
tail -f /var/log/pam_ohmage.log