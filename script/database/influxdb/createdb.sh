#!/usr/bin/env bash

curl -XPOST "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE libtsdbtest"
curl -XPOST "http://localhost:8086/query" --data-urlencode "q=CREATE DATABASE libtsdb"
