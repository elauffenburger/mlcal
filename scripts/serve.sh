#!/usr/bin/env bash
while true; do cat ./index | nc -l 80; done