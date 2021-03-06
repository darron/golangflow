test:
	go test
build:
	buffalo build
deploy:
	heroku container:login
	heroku container:push web --app golangflow
	heroku run ./bin/app migrate --app golangflow
psql:
	heroku pg:psql --app golangflow
backup:
	heroku pg:backups:capture --app golangflow
pgdl:
	heroku pg:backups:download --app golangflow