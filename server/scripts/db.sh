#!/bin/zsh

NAME=postgres_inst
SECRET=secret
DBNAME=timestamps
DATA="/data/database/postgres"

docker run -d --rm \
	--name ${NAME} \
	-p 54320:5432 \
	-e POSTGRES_PASSWORD=${SECRET} \
	-e POSTGRES_DB=${DBNAME} \
	-v $(pwd)/${DATA}:/var/lib/postgresql/data \
	postgres
