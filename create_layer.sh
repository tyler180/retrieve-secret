#!/usr/bin/env bash

# Exit on errors or undefined variables, and pipe failures
set -euo pipefail

########################################
# Variables (customize as needed)
########################################

# Directory where your "mysecrets" source code resides.
# It should contain secrets.go, go.mod, go.sum, etc.
MYSECRETS_SOURCE_DIR="./retrieve-secret"

# The module path where you'd place code in /opt/go/src/... 
# inside the layer (e.g., "github.com/myorg/mysecrets").
# This should match your actual Go module path if you have one.
MODULE_PATH="github.com/tyler180/retrieve-secret"

# The parent folder used to assemble the layer structure
LAYER_BUILD_DIR="./layer"

# Name of the final zip file
OUTPUT_ZIP="retrieve_secret_layer.zip"

########################################
# 1. Clean up & Recreate the Layer Build Directory
########################################

echo "Cleaning up any existing layer build directory..."
rm -rf "${LAYER_BUILD_DIR}"
mkdir -p "${LAYER_BUILD_DIR}/opt/go/src/${MODULE_PATH}"

########################################
# 2. Copy Your Go Source Into the Layer Structure
########################################

echo "Copying source files from '${MYSECRETS_SOURCE_DIR}'..."
cp "${MYSECRETS_SOURCE_DIR}/secrets.go"  "${LAYER_BUILD_DIR}/opt/go/src/${MODULE_PATH}/"
cp "${MYSECRETS_SOURCE_DIR}/go.mod"      "${LAYER_BUILD_DIR}/opt/go/src/${MODULE_PATH}/"
cp "${MYSECRETS_SOURCE_DIR}/go.sum"      "${LAYER_BUILD_DIR}/opt/go/src/${MODULE_PATH}/" || true

# If you have additional files, copy them as well:
# cp "${MYSECRETS_SOURCE_DIR}/any_other_file.go" "${LAYER_BUILD_DIR}/opt/go/src/${MODULE_PATH}/"

########################################
# 3. Zip Everything
########################################

echo "Zipping the layer contents into '${OUTPUT_ZIP}'..."
(
  cd "${LAYER_BUILD_DIR}" 
  zip -r "../${OUTPUT_ZIP}" .
)

echo "Done! Created '${OUTPUT_ZIP}'."
