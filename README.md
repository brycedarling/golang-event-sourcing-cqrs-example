# Go Event Sourcing & CQRS Example

Events are stored in Postgres using [Message DB](https://github.com/message-db/message-db)

Events are processed in batches by "components" and "aggregators" that subscribe to event streams.

Components handle commands and generate events.

Aggregators generate view data from events.

Commands are written to Postgres. Queries are read from Redis.


## Example gRPC Usage

Start up the http/json web API and gRPC server, and postgres and redis:

```
$ docker-compose up -d
Starting postgres ... done
Starting redis    ... done
Starting golang-event-sourcing-cqrs-example_server_1 ... done
```

Connect to the gRPC server using a universal client like [Evans](https://github.com/ktr0731/evans):

```
$ evans --reflection

...

practical.PracticalService@127.0.0.1:50051>
```

Try call an RPC but fail for auth reasons:

```
practical.PracticalService@127.0.0.1:50051> call Viewing
command call: rpc error: code = Unauthenticated desc = missing authorization token
```

Login to get a JWT:

```
practical.PracticalService@127.0.0.1:50051> call Login
email (TYPE_STRING) => bryce@darling.com
password (TYPE_STRING) => foobarbaz
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiI5YzAzMTcwMS0xYjNhLTExZWItODQxNC0wMjQyNzc5ZjE0ZmUiLCJleHAiOjE2MDQ3NTgyMDYsImlzcyI6Im1pY3JvIn0.WjMubCvTalkN4J-mq63JWH3VXDrs9-GL5EExky0QACA"
}
```

Now set the `authorization` header and your RPC call will now succeed:

```
practical.PracticalService@127.0.0.1:50051> header authorization="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySUQiOiI5YzAzMTcwMS0xYjNhLTExZWItODQxNC0wMjQyNzc5ZjE0ZmUiLCJleHAiOjE2MDQ3NTgyMDYsImlzcyI6Im1pY3JvIn0.WjMubCvTalkN4J-mq63JWH3VXDrs9-GL5EExky0QACA"

practical.PracticalService@127.0.0.1:50051> call Viewing
{
  "viewing": {}
}
```

When you are done, shut 'er all down:

```
docker-compose down
```
