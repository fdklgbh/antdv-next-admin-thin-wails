#!/bin/sh
set -e

# After `apt remove`, dpkg retains conffiles. On purge, remove a policy that
# AppArmor might have reloaded in the meantime, without needing its source file.
[ "${1:-}" = "purge" ] || exit 0
[ -n "${DPKG_MAINTSCRIPT_PACKAGE:-}" ] || exit 0

profile_name="${DPKG_MAINTSCRIPT_PACKAGE}-installed"
if [ -r /sys/kernel/security/apparmor/profiles ] &&
    grep -Fq "$profile_name (" /sys/kernel/security/apparmor/profiles; then
  printf '%s' "$profile_name" > /sys/kernel/security/apparmor/.remove
fi
# dpkg removes the packaged /etc/apparmor.d conffile on purge.
