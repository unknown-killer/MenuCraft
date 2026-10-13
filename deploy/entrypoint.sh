#!/bin/sh
# Container entrypoint.
#
# Demo seeding is opt-in and only runs on a brand-new database. Keep it disabled
# in production: the sample account has a publicly documented password.
set -e

if [ "${MENU_SEED_DEMO:-0}" = "1" ] && [ ! -f "${MENU_DB}" ]; then
    echo "first run: seeding the demo restaurant"
    seed || echo "seeding failed; starting with an empty database"
fi

exec menu
