#!/usr/bin/env bash

mongoimport -d saved -c links --jsonArray links.json
mongoexport -d saved -c links --csv --fields url,created_on,category,done > links.csv

sqlite3 saved.sqlite

# Then create table:
# CREATE TABLE 'links' ('id' INTEGER PRIMARY KEY NOT NULL, 'url' VARCHAR NOT NULL, 'category' VARCHAR NOT NULL, 'created_on' DATETIME NOT NULL DEFAULT (CURRENT_TIMESTAMP), 'done' BOOL);

# Then open with GUI to import CSV or try:
# .mode csv
# .import links.csv links
