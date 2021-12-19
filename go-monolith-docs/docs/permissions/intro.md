---
sidebar_position: 1
---

# User permissions

By default we support following permissions:
```go
ProjectPermRegistry = NewPerm()
ProjectPermRegistry.AddPermission("read", ReadPermBit)
ProjectPermRegistry.AddPermission("add", AddPermBit)
ProjectPermRegistry.AddPermission("edit", EditPermBit)
ProjectPermRegistry.AddPermission("delete", DeletePermBit)
ProjectPermRegistry.AddPermission("publish", PublishPermBit)
ProjectPermRegistry.AddPermission("revert", RevertPermBit)
```
Each admin page has CustomPermission property. So, you can associate Admin page with custom permissions, if you want.  
Don't forget to sync up content types after adding custom permission. And also you should add your permission to ProjectPermRegistry.
```go
ProjectPermRegistry.AddPermission("send_notification", NotificationPermBit)
```
