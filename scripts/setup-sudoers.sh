#!/usr/bin/env bash

########################################################################################################################
# This script creates the "/etc/sudoers.d/greenstar" file with the commands that "make setup-environment" needs to run
# with "sudo".
#
# DANGER: modifying the "sudo" configuration files is dangerous, and can prevent administrative access to the machine
#         if done incorrectly. This script was tested but since it's super small, we recommend you review and make sure
#         it applies & fits to your specific OS type & version.
#
# WARNINGS:
#   - *** DO NOT RUN THIS SCRIPT WITH "SUDO"! IT SHOULD BE RUN WITH YOUR NORMAL USERNAME. ***
#   - This script WILL invoke "sudo" multiple times; be prepared to enter a password to authenticate.
#   - This script IS OPTIONAL and IS NOT A REQUIREMENT of the development experience. It is meant to save you from
#     entering your password during "make setup"
#   - Currently, there is one command that `make setup-environment` runs with `sudo` which is not covered by this script
#     -> sudo bash -c "echo \"nameserver 127.0.0.1\" > /etc/resolver/$(DOMAIN)"
#     Unfortunately I'm not sure yet how to codify and escape it properly - help is welcome in that aspect! :)
########################################################################################################################

DOMAIN="${1}"
[[ -z "${DOMAIN}" ]] && echo "usage: ${0} <domain>" && exit 1

set -euxo pipefail

# Create a temporary sudoers file
TMP_SUDOERS_FILE="/tmp/greenstar-sudoers"
cat > "${TMP_SUDOERS_FILE}" <<EOF
${USER} ALL=(ALL) NOPASSWD: $(brew --prefix)/bin/brew services restart dnsmasq
${USER} ALL=(ALL) NOPASSWD: /bin/mkdir -p /etc/resolver
${USER} ALL=(ALL) NOPASSWD: /usr/bin/dscacheutil -flushcache
${USER} ALL=(ALL) NOPASSWD: /usr/bin/killall -HUP mDNSResponder
EOF

# Critical! Verify the file syntax
visudo -c -f "${TMP_SUDOERS_FILE}"

# Critical! Sudoers files MUST have specific ownership & permissions!
sudo chown root:wheel "${TMP_SUDOERS_FILE}"
sudo chmod 440 "${TMP_SUDOERS_FILE}"

# Store file in the system configuration
sudo mv -v "${TMP_SUDOERS_FILE}" "/etc/sudoers.d/greenstar"
