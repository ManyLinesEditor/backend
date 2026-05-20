#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "postgres" \
  -c "CREATE DATABASE ${STORAGE_DB};" \
  -c "CREATE DATABASE ${PAYMENT_DB};"