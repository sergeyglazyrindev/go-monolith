---
sidebar_position: 1
---

# GoMonolith command system intro

GoMonolith provides interface to add easily your own commands, like: migrate data or anything you want.  
Try to execute following command in your CLI
```bash
  ./gomonolith_binary
```
You'll see all commands that registered in your application

# Example

You can check example [here](https://github.com/sergeyglazyrindev/go-monolith/blob/master/superadmin.go). This is a command to create superuser for your project.
To use it, execute following command:
```bash
  ./gomonolith_binary superuser create -n {YOURUSERNAME} -e {USERNAMEEMAIL}
```
and provide a password for your user.

# API

Each command has to implement [following](https://github.com/sergeyglazyrindev/go-monolith/blob/master/core/command_interfaces.go#L3) interface.  
If you want to implement subcommands, then initialize CommandRegistry described [here](https://github.com/sergeyglazyrindev/go-monolith/blob/master/core/command_registry.go#L8).
`Proceed` method gets `subaction` as string and `slice with args` passed from cli.  
Then do what you created the command for.

# Register your command in app

After you added your own command, don't forget to register it in your app. An example is [here](https://github.com/sergeyglazyrindev/go-monolith/blob/81bb7a244db87d71fa7d5d0bc38aee40242c76de/app.go#L123)  
You can use application instance in your cmd, like it's shown [here](../intro/#generate-a-new-project)
