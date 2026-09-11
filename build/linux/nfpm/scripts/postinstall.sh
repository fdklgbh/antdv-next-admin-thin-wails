#!/bin/sh

# dpkg supplies the package name; other package formats do not load this profile.
# GTK3 packages have no profile, so Ubuntu 22.04 never parses AppArmor 4 syntax.
if [ "${1:-}" = "configure" ] && [ -n "${DPKG_MAINTSCRIPT_PACKAGE:-}" ]; then
  profile="/etc/apparmor.d/${DPKG_MAINTSCRIPT_PACKAGE}-installed"
  if [ -f "$profile" ] && [ -r /sys/module/apparmor/parameters/enabled ] &&
      [ "$(cat /sys/module/apparmor/parameters/enabled)" = "Y" ]; then
    # Read the current policy rather than a potentially stale compiled cache.
    apparmor_parser -r -T "$profile" || exit 1
  fi
fi

# Refresh the theme cache; pixmaps also supplies an unthemed lookup fallback.
if command -v gtk-update-icon-cache >/dev/null 2>&1; then
  gtk-update-icon-cache -f -t /usr/share/icons/hicolor
fi

# Update desktop database for .desktop file changes
# This makes the application appear in application menus and registers its capabilities.
if command -v update-desktop-database >/dev/null 2>&1; then
  echo "Updating desktop database..."
  update-desktop-database -q /usr/share/applications
else
  echo "Warning: update-desktop-database command not found. Desktop file may not be immediately recognized." >&2
fi

# Update MIME database for custom URL schemes (x-scheme-handler)
# This ensures the system knows how to handle your custom protocols.
if command -v update-mime-database >/dev/null 2>&1; then
  echo "Updating MIME database..."
  update-mime-database -n /usr/share/mime
else
  echo "Warning: update-mime-database command not found. Custom URL schemes may not be immediately recognized." >&2
fi

exit 0
