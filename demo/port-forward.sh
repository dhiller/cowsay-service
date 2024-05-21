#!/usr/bin/env bash

(
    while true; do
        kubectl port-forward svc/fortune 9090:9090
    done
)&
ps ax | grep port-forward
