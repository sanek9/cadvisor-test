#!/bin/sh
set -x

export CIMAGE=cadvisor:19dbf410
docker-compose build
RANDOM_METRIC_COUNT=80 CUSTOM_LABEL=custom_label docker-compose up --abort-on-container-exit #ok
echo status $?
RANDOM_METRIC_COUNT=100 CUSTOM_LABEL=custom_label docker-compose up --abort-on-container-exit  #ok
echo status $?
#RANDOM_METRIC_COUNT=101 CUSTOM_LABEL=custom_label docker-compose up --abort-on-container-exit  #fail
#echo status $?
RANDOM_METRIC_COUNT=80 CUSTOM_LABEL=name docker-compose up --abort-on-container-exit  #ok
echo status $?
