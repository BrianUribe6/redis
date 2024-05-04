[![progress-banner](https://backend.codecrafters.io/progress/redis/d009a84a-41ba-4ce2-8377-90067b6c151a)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

This is my solution to
["Build Your Own Redis" Challenge](https://codecrafters.io/challenges/redis).

It supports the following commands:

- `PING`
- `SET`
- `GET`
- `ECHO`
- `PSYNC` \*
- `INFO`\*

\* _partial support_

You can start a local server by running

```
./spawn_redis_server.sh
```

This will start a redis server on port **6379** by default.

You can now spawn another replica and have it synchronize with master by
running

```
./spawn_redis_server.sh --port 3000 --replicaof localhost 6379
```

Check if the server is running correctly by running the original `redis-cli`

```
$ redis-cli ping
PONG
```
