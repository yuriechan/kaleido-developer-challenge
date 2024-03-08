#!/bin/bash

set -ex

echo "** Creating default DB"

mysql -u "root" -p"$MYSQL_ROOT_PASSWORD" --execute \
"CREATE USER '$MYSQL_USER'@'$MYSQL_ROOT_HOST' IDENTIFIED WITH mysql_native_password BY '{$MYSQL_PASSWORD}';
GRANT ALL ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'$MYSQL_ROOT_HOST';
FLUSH PRIVILEGES;

USE $MYSQL_DATABASE;

CREATE TABLE Users (
    Id varchar(255) NOT NULL,
    Name varchar(255) NOT NULL,
    CreatedBy date DEFAULT (CURRENT_DATE),
    UpdatedBy date DEFAULT (CURRENT_DATE)
);
"



echo "** Finished creating DB and root user"