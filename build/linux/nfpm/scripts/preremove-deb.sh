#!/bin/sh
set -e

# Keep the policy active during upgrades; the new postinst replaces it.
[ "${1:-}" = "remove" ] || exit 0
[ -n "${DPKG_MAINTSCRIPT_PACKAGE:-}" ] || exit 0

profile_name="${DPKG_MAINTSCRIPT_PACKAGE}-installed"
profile="/etc/apparmor.d/$profile_name"
if [ -f "$profile" ] && [ -r /sys/kernel/security/apparmor/profiles ] &&
    grep -Fq "$profile_name (" /sys/kernel/security/apparmor/profiles; then
  apparmor_parser -R "$profile"
fi
