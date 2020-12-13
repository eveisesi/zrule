# Welcome to ZRule

### What is ZRule

ZRule is a Rule Engine built for ZKillboard. If you don't know what ZKillboard is, jump on over to that [bloke's repo](https://github.com/zkillboard/zkillboard) to read more, but in a nut shell, Zkillboard is a Killboard for the MMORPG EVE Online. If you don't know what EVE is, then I assume you are here for the technology and not for the reason.

ZRule is written in Go. It is a single go program run in 4 seperate containers backed by Redis and Mongo. Its host three workers and on HTTP API.

ZRule works by listening to the [ZKillboard Websocket](https://github.com/zKillboard/zKillboard/wiki/Websocket). When a Killmail comes in, we compare that killmail to the list of rules that currently exist, if a match is found, the policy that the rule belongs to is pulled and the actions belonging to that policy are triggered.
