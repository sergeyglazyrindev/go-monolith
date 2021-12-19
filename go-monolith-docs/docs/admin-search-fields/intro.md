---
sidebar_position: 1
---

# Admin search field

You may search for the records by any field in the database. You can initialize search field registry using helper.  
```go
searchFieldRegistry = NewSearchFieldRegistryFromGormModel(modelI)
```
If you want to search in the admin page by any field, just tag this field with `gomonolith:"search"`.  
Also, you can customize search by providing CustomSearch function for searchField.
```go
CustomSearch func(afo IAdminFilterObjects, searchString string)
```
Later on we will migrate it to interface as well, so it could be used easily for any type of search functionality.
