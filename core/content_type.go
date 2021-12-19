package core

import "fmt"

type ContentType struct {
	Model
	BlueprintName string `sql:"unique_index:idx_contenttype_content_type"`
	ModelName     string `sql:"unique_index:idx_contenttype_content_type"`
}

func (ct *ContentType) String() string {
	return fmt.Sprintf("Content type for blueprint %s and model name %s", ct.BlueprintName, ct.ModelName)
}
