# Welcome to ZRule

### What is ZRule

ZRule is a Rule Engine built for ZKillboard. If you don't know what ZKillboard is, jump on over to that [bloke's repo](https://github.com/zkillboard/zkillboard) to read more, but in a nut shell, Zkillboard is a Killboard for the MMORPG EVE Online. If you don't know what EVE is, then I assume you are here for the technology and not for the reason this application was built.

ZRule is written in Go. It is a single go program that is run in 4 seperate docker containers backed by Redis and Mongo, also run in their own docker containers. It hosts three workers and an HTTP API.

ZRule works by listening to the [ZKillboard Websocket](https://github.com/zKillboard/zKillboard/wiki/Websocket). When a Killmail comes in, we compare that killmail to the list of rules that currently exist, if a match is found, the policy that the rule belongs to is pulled and the actions belonging to that policy are triggered. Zrule Support three seperate action types. Slack, Discord, and REST. Slack and Github have the url to the killmail posted them via a webhook that is supplied when an Action is created. The intention for REST, is you expose an Endpoint to ZRule that ZRule can make a POST request to containing the id and hash of the killmail that matched your policy. This is all good is theory, but we have no been able to fully test this yet.
