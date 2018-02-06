#!/bin/bash

hostport="$1"
source="$2"
name="$3"
season="$4"
episode="$5"
title="$6"
when="$7"
duration="$8"

b64source=$(echo -n "$source" | base64)
b64name=$(echo -n "$name" | base64)
b64season=$(echo -n "$season" | base64)
b64episode=$(echo -n "$episode" | base64)
b64title=$(echo -n "$title" | base64)

curl http://${hostport}/rec/${b64source}/${b64name}/${b64season}/${b64episode}/${b64title}/${when}/${duration}
