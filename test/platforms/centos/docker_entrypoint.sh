#!/bin/bash
source /common.sh

startRstudio

echo "Starting syslog-ng"
syslog-ng

touch /var/log/pam_ohmage.log
echo "Tailing /var/log/pam_ohmage.log"
tail -f /var/log/pam_ohmage.log