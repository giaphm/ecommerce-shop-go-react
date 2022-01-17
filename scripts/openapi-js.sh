#!/bin/bash
for service in checkouts orders products users; do
    /usr/local/bin/docker-entrypoint.sh generate \
        -i "./api/openapi/$service.yml" \
        -g javascript \
        -o "./web/src/repositories/clients/$service"
done