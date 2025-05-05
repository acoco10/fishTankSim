#!/bin/zsh
GOOS=js GOARCH=wasm go build -tags=g -o main.wasm

bucket_name="fish-fish-fish-assets"
object_key="wasm/"
local_file="main.wasm"

echo "Uploading $local_file to s3://$bucket_name/$object_key"

aws s3 cp "$local_file" "s3://$bucket_name/$object_key" --content-type "application/wasm"


if [ $? -eq 0 ]; then
  echo "Successfully uploaded to S3!"
else
  echo "Error during S3 upload."
  exit 1
fi