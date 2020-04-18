# weather-reporting-api
API building exercise

## Install Go
Before running the api please make sure you have Go installed.
https://golang.org/doc/install

## Clone the repo
You can clone the repo using `git clone` or `go get`.

## Buiding, Testing and Running
Use makefile commands to get things working:

Command | Description
------------ | -------------
`make build` | Build the API
`make test` | Run all automated tests
`make run` | Run the API

## API Endpoints Examples

### Auth

Request:
```
POST http://localhost:8080/auth/
{
	"name": "kirang",
	"password": "secret"
}
```
Success Response:
```
{
    "token": "3ac9f318f426aef056f46a9e02b69d08b8a92646"
}
```

### Save

Request:
```
POST http://localhost:8080/weather/
"Authorization": "3ac9f318f426aef056f46a9e02b69d08b8a92646"
{
	"city": "vancouver",
	"weather": [{
		"date": "2020-04-17",
		"temperature": 17
	},
	{
		"date": "2020-04-18",
		"temperature": 18
	},
	{
		"date": "2020-05-18",
		"temperature": 16
	}]
}
```
Success Response:
```
{
    "message": "The weather was saved succesfully!"
}
```

### Get

Request:
```
GET http://localhost:8080/weather/
"Authorization": "3ac9f318f426aef056f46a9e02b69d08b8a92646"
{
	"city": "vancouver",
	"initial_date": "2020-04-01",
	"end_date": "2020-04-30"
}
```
Success Response:
```
{
    "city": "vancouver",
    "weather": [
        {
            "date": "2020-04-17",
            "temperature": 17
        },
        {
            "date": "2020-04-18",
            "temperature": 18
        }
    ]
}
```

### Delete

Request:
```
DELETE http://localhost:8080/weather/
"Authorization": "3ac9f318f426aef056f46a9e02b69d08b8a92646"
{
	"city": "vancouver"
}
```
Success Response:
```
{
    "message": "The weather was deleted succesfully!"
}
```