v0.5.0

## API Changes

- Add `Handler#Register(interface{})` endpoint which is shorthand for `Handler#RegisterCommand(Command)` or `RegisterMiddleware(Middleware)`, depending on the passed object instance implementation.

- `Handler#RegisterHandlers(discordgo.Session)` is now **deprecated**. Please use `Handler#Setup(discordgo.Session)` instead.

- `Handler#SetObject(discordgo.Session)` is now **deprecated**. Please pass a `di.Container` to register handler specific instances.

## Dependency Injection Container Implementation

shireikan now uses [sarulabs/di](https://github.com/sarulabs/di) for object reference injection on handler level into command `Context` instances.

See [**this example**](https://github.com/zekroTJA/shireikan/tree/master/examples/di) on how to implement this.