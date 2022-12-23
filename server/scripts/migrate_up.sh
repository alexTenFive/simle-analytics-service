#!/bin/zsh
migrate -path $(pwd)/migrations -database "postgres://postgres:secret@localhost:54320/timestamps?sslmode=disable" up 1
