#!/usr/bin/env bash

./protoc --micro_out=. --go_out=. $1/*.proto