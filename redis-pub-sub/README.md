#### Redis pub sub sample scripts

##### Usecase:

```
1. Run server.rb on external network
2. Run subscriber.rb on internal network
3. Hit server.rb from external app.
4. server.rb publishes it to a channel
5. subscriber.rb listens to that channel.
6. you now have the message in internal network
7. you can always hit the external app using the webhooks.

TESTED AND IT WORKS AS EXPECTED
```
