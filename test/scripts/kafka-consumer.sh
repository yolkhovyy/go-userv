#!/bin/bash

docker compose exec -it kafka kafka-console-consumer \
--bootstrap-server kafka:9092 \
--topic postgres.public.users \
--from-beginning
