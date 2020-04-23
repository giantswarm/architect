#!/bin/bash

TAG=$1
PERSONAL_ACCESS_TOKEN=$2

echo "Creating GitHub Release"
release_output=$(curl \
    --request GET \
    --header "Authorization: token ${PERSONAL_ACCESS_TOKEN}" \
    --header "Content-Type: application/json" \
    https://api.github.com/repos/giantswarm/architect/releases/tags/$TAG
)
echo $release_output | jq

# fetch the release id for the upload
RELEASE_ID=$(echo $release_output | jq '.id')

echo "Upload binary to GitHub Release"
upload_status=$(curl \
    -s -o /dev/null -w "%{http_code}" \
    --header "Authorization: token ${PERSONAL_ACCESS_TOKEN}" \
    --header "Content-Type: application/octet-stream" \
    --data-binary @architect \
    https://uploads.github.com/repos/giantswarm/architect/releases/${RELEASE_ID}/assets?name=architect
)

code=${upload_status:0:1}
if [ "$code" != "2" ]; then
    echo "Upload failed, status code $upload_status"
    exit 1
fi

echo "Done!"
