#!/bin/bash

if [ -z "$ENV_FILE" ]
then
  echo "CONFIG FILE is required"
  exit
fi

# shellcheck disable=SC1090
. "$ENV_FILE"

# check all env vars needed have been imported
if [ -z "$DB_NAME" ] || [ -z "$DB_PASSWORD" ] || [ -z "$DB_USER" ] || [ -z "$DB_PORT" ] || [ -z "$DB_HOST" ]
then
      echo "\$DB_NAME, \$DB_PASSWORD, \$DB_USER, \DB_PORT, \DB_HOST are required"
      exit
fi

echo "Waiting postgres to launch on $DB_PORT..."

until docker exec postgres psql -U postgres -h "$DB_HOST" -c "select 1" &>/dev/null; do sleep 1; done

echo "Postgres script launching"

# Create postgresql user
USER_CREATED=$(docker exec postgres psql -h "$DB_HOST" -U postgres -tc "SELECT 'created' FROM pg_user WHERE usename = '$DB_USER'" | grep created)
if [ ${#USER_CREATED} -gt 0 ]; then
  echo "User exists, continuing"
else
  echo "Creating user"
  docker exec postgres psql -h "$DB_HOST" -U postgres -c "CREATE ROLE $DB_USER LOGIN PASSWORD '$DB_PASSWORD';"
fi

# Check if db already exists
DB_CREATED=$(docker exec postgres psql -h "$DB_HOST" -U postgres -tc "SELECT 'created' FROM pg_database WHERE datname = '$DB_NAME'" | grep created)
if [ ${#DB_CREATED} -gt 0 ]; then
  echo "Deleting existing $DB_NAME database"
  docker exec postgres psql -h "$DB_HOST" -U postgres -c "DROP DATABASE $DB_NAME;"
fi

# Create new db
echo "Creating $DB_NAME database"
docker exec postgres psql -h "$DB_HOST" -U postgres -c "CREATE DATABASE $DB_NAME;"