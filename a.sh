#! /bin/sh

for state in $(cat main.go)
do
echo "$state" >> log1.txt
done