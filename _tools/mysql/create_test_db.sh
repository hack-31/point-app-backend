#!/bin/bash
mysql -uroot -e "create database $DB_TEST_NAME;"
mysql -uroot -e "grant all on $DB_TEST_NAME.* to $MYSQL_USER@'%';"
