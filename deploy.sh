#!/bin/bash

ssh -A apps@vpsargeweb_28153 'cd /apps/hashcd && git pull'
ssh -A apps@vpsargeweb_28153 'cd /apps/hashcd && docker compose up -d --build --force-recreate --no-deps'
