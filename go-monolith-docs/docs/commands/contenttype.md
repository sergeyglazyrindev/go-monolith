# Content type command

Use this command to sync up Content Types for the models in your project. You have to do that if you want to use this model in the admin panel. Or if you want to be able to manipulate who from users in your project would have an access to this data model. Please don't forget to register your model that you want to add contenttype for in the ProjectModelRegistry, you can do that in blueprint `Init()` func
Please keep in mind that if you add model as admin page, it will done automatically
```go
func (b Blueprint) Init() {
  core.ProjectModels.RegisterModel(func() (interface{}, interface{}) { return &core.User{}, &[]*core.User{} })
}
```

# Usage

To use it, execute following command:
```bash
  ./gomonolith_binary contenttype sync
```
As result of this command execution, you would have a directory structure created for your blueprint, will be added first empty migration.  
You can keep it empty if your blueprint doesn't interact with database.
