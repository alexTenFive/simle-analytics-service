#!/bin/zsh
docker-compose up -d --build postgres app && docker-compose logs -f
