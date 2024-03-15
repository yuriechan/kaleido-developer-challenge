#!/bin/bash

set -ex

echo "** Creating default DB"

mysql -u "root" -p"$MYSQL_ROOT_PASSWORD" --execute \
"CREATE USER '$MYSQL_USER'@'$MYSQL_ROOT_HOST' IDENTIFIED WITH mysql_native_password BY '{$MYSQL_PASSWORD}';
GRANT ALL ON $MYSQL_DATABASE.* TO '$MYSQL_USER'@'$MYSQL_ROOT_HOST';
FLUSH PRIVILEGES;

CREATE TABLE IF NOT EXISTS $MYSQL_DATABASE.users (
    id varchar(255) NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    ethereum_address varchar(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

INSERT IGNORE INTO $MYSQL_DATABASE.users
    (id, name, ethereum_address)
VALUES ('1', 'User 1', '0x68173f054e5b5a588dab145f4afd86e99c35616f');

INSERT IGNORE INTO $MYSQL_DATABASE.users
    (id, name, ethereum_address)
VALUES ('2', 'User 2', '0xa6b1b1848a1cdf9fb19c0e22f0bd9094a06ed886');

INSERT IGNORE INTO $MYSQL_DATABASE.users
    (id, name, ethereum_address)
VALUES ('3', 'User 3', '0x63591f52aa7b150aa7a6e881e8fc61918f064665');

CREATE TABLE IF NOT EXISTS $MYSQL_DATABASE.items (
    id varchar(255) NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    state varchar(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);"



echo "** Finished creating DB and root user"