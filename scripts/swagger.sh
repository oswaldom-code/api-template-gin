#!/bin/bash

sudo docker run --rm -i yousan/swagger-yaml-to-html < swagger/swagger.yml | sudo tee doc/api.html > /dev/null
