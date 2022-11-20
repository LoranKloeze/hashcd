#!/bin/bash

ssh -A apps@vpsargeweb_28153 'cd /apps/finalcd && git pull'
ssh -A apps@vpsargeweb_28153 'cd /apps/finalcd && docker compose build'
