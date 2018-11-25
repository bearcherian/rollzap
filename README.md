# Rollbar Zap

A simple zapcore.Core implementation to integrate with Rollbar. 

To use, initialize rollbar like normal, create a new RollbarCore, then wrap with a NewTee. [See the example code](example/main.go) for a detailed example.

## Testing 

To test this code use `RC_TOKEN=MY_ROLLBAR_TOKEN go test`
