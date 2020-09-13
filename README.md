## IRC

An IRC bot that supports plugins which communicate over gRPC.  Plugins can be restart as needed without the
bot disconnecting.

Current plugins:
 - Github notification plugin
 - Webhook notification plugin
 
### Bot

IRC bot that connects to IRC and sits on a single channel.  Supports password and TLS auth for connecting to the server.
Both server and channel are mandatory, all other flags are optional and have sensible defaults.  The bot will always 
have an RPC port and a http port open, so firewall/expose as required.

This can either be run directly with cli arguments, or in docker.  All CLI flags are support as environment variables.

 - go build github.com/greboid/irc/cmd/bot
 - docker run greboid/irc

#### Github plugin

Receives notifications from github and outputs them to a channel.  The RPC token is required.

This can either be run directly with cli arguments, or in docker.  All CLI flags are support as environment variables.

 - go build github.com/greboid/irc/cmd/github
 - docker run greboid/irc-github

#### Webhook plugin

Uses http server on the bot to receive generic notification requests over, auths via API keys, stores data in a 
database.  The RPC token is required.  The list of keys is stored in a database, the CLI argument is the full path to
the database, so if using docker, you can either mount a directory, or the file itself.

This can either be run directly with cli arguments, or in docker.  All CLI flags are support as environment variables.

 - go build github.com/greboid/irc/cmd/web
 - docker run greboid/irc-webhook