# The Future

This project has been more or less abandoned, but I am still working on it.
The last active deployed and used yummy instance at my company isn't running anymore
for quite a while now. But then we had the need for a repository server and with
it came the chance to deploy new versions of yummy.

But, looking at the code after these days I am sure, that the project needs a lot
of polishing and work.

# The Roadmap

First, I had the plan to make yummy not only the repository serving webserver but also a
CLI tool for housekeeping tasks and to may be able to upload rpm packages. Feature creep?
But, as time told, I haven't had the time to implement this.

~~Because of this, I decided to keep yummy as a pure webserver with the simple purpose of
serving rpm packages with the ability to upload and delete `stale` packages. With this
step the cobra library and code can be removed.~~

The code got refactored and the architecture got improved a bit. There is only the server
code left, without the cobra library and the code for the CLI.

~~As of the day when yummy was under active development, dep was a viable option for
dependency management. But today Go modules are the way to go. So it is time to
switch to this instead.~~

Switch from dep to Go modules is finished. There were some issues with the build process
but those could be mitigated.

The next step is to update the router library to a newer release which now satisfies the
HandleFunc() interface.
