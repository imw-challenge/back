=Back Coding Challenge=
=Isaac Wilder=


Post Message Request:
curl -i -H "Content-Type: application/json" -H "Accept: application/json" -X POST -d '{"id":"B5D99898-7DE9-7E69-C311-763310C9AA54","name":"Isaac Wilder","email":"name@fake.domain","text":"hi there","time":"2018-04-10T13:15:17-07:00"}' http://localhost:9000/public/message 

Dump Messages Request:
curl -i --user admin:back-challenge -H "Content-Type: application/json" -H "Accept: application/json" -X GET  http://localhost:9000/private/dump

Get Message Request:
curl -i --user admin:back-challenge -H "Content-Type: application/json" -H "Accept: application/json" -X GET -d '{"id":"E5D99898-7DE9-7E69-C311-763310C9AA54"}' http://localhost:9000/private/message

Update Message Request:
curl -i --user admin:back-challenge -H "Content-Type: application/json" -H "Accept: application/json" -X PUT -d '{"id":"E5D99898-7DE9-7E69-C311-763310C9AA54","text":"hi there"}' http://localhost:9000/private/message


