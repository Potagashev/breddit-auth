Чтобы запустить, введите команду 
`docker-compose up`

Локальная разработка
`go run main.go`

For MacOS:

Чтобы применить миграции, введите команду
(Креды смотреть в .env, а порты в docker-compose.yml)
`brew install golang-migrate`
`migrate -database "postgresql://postgres:postgrespw@localhost:49153/postgres?sslmode=disable" -path ./db/migrations up`

Чтобы создавать миграции, пишем
`migrate create -ext sql -dir db/migrations -seq name_of_revision`
