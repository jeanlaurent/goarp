#!/bin/sh
cat oui.txt  | grep "(base 16)" | awk '{$2=$3=""; print $0}' | sed 's/   />/g'