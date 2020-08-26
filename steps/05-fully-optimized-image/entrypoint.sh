#!/bin/sh

# Check capabilities
HAS_CHOWN=$(capsh --print | grep cap_chown)
if [[ "${HAS_CHOWN}" != "" ]]; then
    echo "[ERROR] Kernel capability 'cap_chown' is enabled!"
    echo "exiting..."
    exit 1
fi

# Start command
$@
