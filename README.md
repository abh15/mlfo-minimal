# mlfo-minimal
This work is done as part of ITU AI/ML 5G Challenge 2020

## Requires
go v1.14.5

Docker v19.03.13

Docker compose v1.27.4

## Usage
### Setup model repo db

`docker-compose up`

`docker exec -it db bash -c "mysql -uroot -pmlfo1234 modelrepo < modelrepo.sql"`

### Run

For edge use case:

`go run main.go edge_intent.yaml`

For central cloud use case:

`go run main.go cloud_intent.yaml`
