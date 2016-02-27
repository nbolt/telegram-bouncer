Build
-----

`go build`


Run
----


`TELEGRAM_TOKEN=xxx REDIS_URL=xxx REDIS_PASS=xxx REDIS_DB=xxx PORT=xxx ./telegram-bouncer`


API
---

`/telegram/:user_id/:message_type(/:text)`

`:message_type` can be: `message`, `photo`, or `document`

for `message` types, include the text as another url parameter.

for `photo` and `document` types, include a `file` parameter as part of a `multipart/form-data` `POST` request.
