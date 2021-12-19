---
sidebar_position: 1
---

# Admin list filter

Admin list filter populates your right sidebar with options to filter your objects. Like in [Django](https://www.dothedev.com/blog/django-admin-list_filter/)  
An example of how it could be used.
```go
listFilter := &core.ListFilter{
	URLFilteringParam: "IsSuperUser__exact",
	Title:             "Is super user ?",
}
listFilter.OptionsToShow = append(listFilter.OptionsToShow, &core.FieldChoice{DisplayAs: "Yes", Value: true})
listFilter.OptionsToShow = append(listFilter.OptionsToShow, &core.FieldChoice{DisplayAs: "No", Value: false})
```
Later on we will migrate it to interface as well, so it could be used easily for any type of list filter functionality.
