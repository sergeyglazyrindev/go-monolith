---
sidebar_position: 1
---

# File storages

By default we support only file and s3 storage. Later on we will provide storages for another providers, etc.  
It could be used like here:
```go
if storage == nil {
  storage = NewFsStorage()
}
ret := make([]string, 0)
var filename string
// files from POST request
for _, file := range files {
  f, err := file.Open()
  if err != nil {
    return err
  }
  bytecontent := make([]byte, file.Size)
  _, err = f.Read(bytecontent)
  if err != nil {
    return err
  }
  filename, err = storage.Save(&FileForStorage{
    Content:           bytecontent,
    PatternForTheFile: "*." + strings.Split(file.Filename, ".")[1],
    Filename:          file.Filename,
  })
  if err != nil {
    return err
  }
  err = f.Close()
  if err != nil {
    return err
  }
  ret = append(ret, filename)
}
```

Or you can use S3 storage provider for your field in the form, like below:
```
form := core.NewFormFromModelFromGinContext(ctx, modelI, make([]string, 0), fields, true, "", true)
iconField, _ := form.FieldRegistry.GetByName("Icon")
s3Storage := core.NewAWSS3Storage("test", &core.AWSConfig{
  S3: &core.AWSS3Config{
    Region:    "{YOUR_REGION}",
    AccessKey: "{YOUR_ACCESS_KEY}",
    SecretKey: "{YOUR_SECRET_KEY}",
  },
})
s3Storage.(*core.AWSS3Storage).Bucket = "test"
s3Storage.(*core.AWSS3Storage).Domain = "https://{YOUR_PART}.s3.eu-central-1.amazonaws.com/test"
iconField.FieldConfig.Widget.(*core.FileWidget).Storage = s3Storage
```
