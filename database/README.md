```sh
docker build -t run-goals/db:latest .

docker run -d --name database -p 5432:5432 run-goals/db:latest

docker exec -it database psql -U postgres -d run_goals

psql -U postgres -d run_goals -W

psql \dt
```
