#!/bin/bash

go build -o baltijaskauss ./cmd/web/*.go 
./baltijaskauss -dbname=baltijas_kauss -dbuser=op -cache=false -production=true