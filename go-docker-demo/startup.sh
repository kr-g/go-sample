#!/bin/bash

export BDIR=${1:-~} && if [ ! -d $BDIR ]; then echo [fail] $BDIR does not exist; exit 1; fi


echo path $PATH
echo goroot $GOROOT


cd $BDIR/repo/go-sample/webapp-static
go run main.go &

cd $BDIR/repo/go-sample/webapp-static-markdown
go run main.go &


vmstat 5

