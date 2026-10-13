#!/bin/sh
# Container entrypoint.
#
# On a brand-new volume this seeds the demo restaurant so whoever opens the URL
# lands on a working menu rather than an empty panel. On every later start the
# database already exists and seeding is skipped, so real data is never touched.
set -e

if [ "${MENU_SEED_DEMO:-1}" = "1" ] && [ ! -f "${MENU_DB}" ]; then
    echo "first run: seeding the demo restaurant"
    seed || echo "seeding failed; starting with an empty database"
fi

exec menu
