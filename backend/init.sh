#!/bin/bash

set -ex

echo "** Creating default DB"

mysql -u "root" -p"$MYSQL_ROOT_PASSWORD" --execute \
"CREATE USER '$MYSQL_USER'@'$MYSQL_ROOT_HOST' IDENTIFIED WITH mysql_native_password BY '{$MYSQL_PASSWORD}';
GRANT ALL ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'$MYSQL_ROOT_HOST';
FLUSH PRIVILEGES;

CREATE TABLE IF NOT EXISTS $MYSQL_DATABASE.listing (
    id varchar(255) NOT NULL PRIMARY KEY,
    item_name varchar(255) NOT NULL,
    item_state varchar(255) NOT NULL,
    item_price INT NOT NULL,
    nft_id varchar(255) NOT NULL,
    smart_contract_address varchar(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);"

echo "** Finished creating DB and root user"