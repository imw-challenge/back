
# Back  Challenge
 Submitted by Isaac Wilder
 
## Notes
In general, the exercise was straightforward. I decided to use an in-memory database for fun, and since the brief mentioned a maximum possible dataset of a few gigabytes. It should be fairly snappy! I've included a Dockerfile for your convenience, as well as a trivial Travis CI buildfile. The service is also deployed on my VPS, and reachable at 0x539.lol:9000. The API Spec follows, as well as Docker and Curl examples.

## API Spec
### Post Message
This is the public post method, located at `/public/message` - it only listens to `POST` requests. It expects a JSON in the body, of the form:
```
{
  "id": "B5D99898-7DE9-7E69-C311-763310C9AA54",
  "name": "Isaac Wilder",
  "email": "name@fake.domain",
  "text": "hi there",
  "time": "2018-04-10T13:15:17-07:00"
}
```
ID and Text are mandatory fields, all other fields are optional. 

### Put Message
This is the private message update method, located at `/private/message` - it only listens to `put` requests, and requires correct HTTP basic auth headers. It expects a JSON body, of the form:
```
{
  "id": "E5D99898-7DE9-7E69-C311-763310C9AA54",
  "text": "hi there"
}
```
ID and Text are mandatory fields, all other fields are ignored.

### Get Message
This is the private message retrieval method, located at `/private/message` - it only listens to `GET` requests, and requires correct HTTP basic auth headers. It expects a JSON body, of the form:
```
{
  "id": "B5D99898-7DE9-7E69-C311-763310C9AA54"
}
```
ID is a mandatory field, all other fields are ignored.

### Get Dump
This is the private reverse chronological dump method, located at `/private/dump` - it only listens to `GET` requests, and requires correct HTTP basic auth headers. It ignores any request body.

## Docker Commands
```
sudo docker build -t imw-back .
sudo docker run -d -p 9000:9000 imw-back
```

## Example Requests
Post Message Request:
```
curl -i -H "Content-Type: application/json" -H "Accept: application/json" -X POST -d '{"id":"B5D99898-7DE9-7E69-C311-763310C9AA54","name":"Isaac Wilder","email":"name@fake.domain","text":"hi there","time":"2018-04-10T13:15:17-07:00"}' http://localhost:9000/public/message 
```

Dump Messages Request:
```
curl -i --user admin:back-challenge -H "Content-Type: application/json" -H "Accept: application/json" -X GET  http://localhost:9000/private/dump
```

Get Message Request:
```
curl -i --user admin:back-challenge -H "Content-Type: application/json" -H "Accept: application/json" -X GET -d '{"id":"E5D99898-7DE9-7E69-C311-763310C9AA54"}' http://localhost:9000/private/message
```

Update Message Request:
```
curl -i --user admin:back-challenge -H "Content-Type: application/json" -H "Accept: application/json" -X PUT -d '{"id":"E5D99898-7DE9-7E69-C311-763310C9AA54","text":"hi there"}' http://localhost:9000/private/message
```
