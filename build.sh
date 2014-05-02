#!/bin/sh -x
cd common && go build common.go && cd ..
cd check_nrpe && go build check_nrpe.go && cd ..
go build nrped.go
