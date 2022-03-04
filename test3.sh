#!/bin/sh
curl -v -X POST -d "{\"coordinates\":\"50.7184,-3.5339\"}" \
localhost:3000/what3words
curl -v -X POST -d "{\"coordinates\":\"52.0907,5.1214\"}" \
localhost:3000/what3words
curl -v -X POST -d "{\"coordinates\":\"44.4056,8.9463\"}" \
localhost:3000/what3words
