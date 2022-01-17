#!/bin/bash
readonly service="$1"
readonly server_to_run="$2"
readonly project_id="$3"

gcloud run deploy "$service-$server_to_run" \
    --image "gcr.io/$project_id/$service" \
    --region asia-southeast1 \
    --platform managed
