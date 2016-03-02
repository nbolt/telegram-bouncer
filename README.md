Build
-----

`go build`


Run
----


`TELEGRAM_TOKEN=xxx REDIS_URL=xxx REDIS_PASS=xxx REDIS_DB=xxx ./telegram-bouncer`


API
---

`POST` `/telegram/:user_id/:message_type`

`:message_type` can be: `message`, `photo`, or `document`

for `photo` and `document` types, include `file` as a `multipart/form-data` parameter.

for `message` types, include `text`.

