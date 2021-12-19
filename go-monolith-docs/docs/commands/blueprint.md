# Blueprint command

Use this command to create new blueprint for your project. **Make sure you execute it in the root of your project**

# Usage

To use it, execute following command:
```bash
  ./gomonolith_binary blueprint create -m "Message that describes the idea of your blueprint" -n "{{SHORT_NAME_OF_YOUR_BLUEPRINT}}"
```
As result of this command execution, you would have a directory structure created for your blueprint, will be added first empty migration.  
You can keep it empty if your blueprint doesn't interact with database.  
Don't forget to register this blueprint in your application, you can do it using following syntax
```go
appInstance.BlueprintRegistry.Register(userblueprint.ConcreteBlueprint)
```
