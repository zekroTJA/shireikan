v0.6.0

- shireikan now makes use of `sync.Pool` for context and rate limiter instances which might make it more performant when handling a lot of commands.
- Fixed registration of command handler instance in context object map. [#3]