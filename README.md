docker exec -it shah-postgres-db psql -U root -d postgres


# Starting
docker compose up
templ generate -path="./view"
go run cmd/shah/main.go