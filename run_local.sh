#!/bin/bash

go build -o baltijaskauss ./cmd/web/*.go 
./baltijaskauss -dbname=production_local -dbuser=postgres -cache=false -production=true
