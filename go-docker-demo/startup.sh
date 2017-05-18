#!/bin/bash

export BDIR=${1:-~} && if [ ! -d $BDIR ]; then echo [fail] $BDIR does not exist; exit 1; fi


echo path $PATH
echo goroot $GOROOT


for i in $BDIR/repo/go-sample/webapp-* 
do 

	echo $i
	cd $i 
	go run main.go &

done


vmstat 5

