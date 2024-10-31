run:
	go run cmd/main.go
sdocker:
	sudo sysctl -w kernel.apparmor_restrict_unprivileged_userns=0
startdocker: sdocker
	systemctl --user start docker-desktop
# staring from this till v4.17.0 it is docker command, after that it is to create migrate file u can see on their documentation
startmigration:
	docker run -it --rm --network host --volume "$(shell pwd)/db:/db" migrate/migrate:v4.17.0 create -ext sql -dir /db/migrations -seq add_alter_table
# this is to create mysql container if doen't exist else it will start the container
rundocker:
	docker run --name ecom -p 3306:3306 -e MYSQL_ROOT_PASSWORD=abate -d mysql:9.1
# create database in docker container #simple-api should be the name of the container we created above
createdb:
	docker exec -it ecom mysql -uroot -pabate -e "CREATE DATABASE ecom;"

# migrate up
migrateup:
	docker run -it --rm --network host --volume "$(shell pwd)/db:/db" migrate/migrate:v4.17.0 -path /db/migrations -database "mysql://root:abate@tcp(localhost:3306)/ecom" up
migratedown:
	docker run -it --rm --network host --volume "$(shell pwd)/db:/db" migrate/migrate:v4.17.0 -path /db/migrations -database "mysql://root:abate@tcp(localhost:3306)/ecom" down
