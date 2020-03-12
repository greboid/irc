## IRC

Hacking around with IRC in go.  

Bot with plugins to do the various logic, connecting via gRPC.  In theory meaning the core bot can just stay on 
IRC and its functionality change without needing to leave and come back.

You'll need to run go generate before you build.  Currently need to build the bot cmd and maybe the test plugin.