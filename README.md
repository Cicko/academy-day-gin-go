The docker-compose.yml file should be in the parent folder.

```
version: '3.2'

services:
  api:
    build: api
    volumes:
      - "./academy-day-gin-go:/asdfasdf"
    command: go run app.go
    ports:
    - "8080:8080"
```


To run just the tests:

```docker-compose run chakram npm run test:users ```
