``` sh
docker build -t rungoals/database:latest .

docker run -d --name database -p 5432:5432 rungoals/database:latest

docker exec -it database psql -U postgres -d rungoals

psql -U postgres -d rungoals -W

psql \dt
```
