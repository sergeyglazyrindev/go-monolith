define PROTOC_GEN
go get github.com/gogo/protobuf/protoc-gen-gogofaster@v1.3.2
go get github.com/golang/protobuf/protoc-gen-go@v1.3.2
export GOMONOLITH_PATH=$(shell pwd)
protoc -I$${GOPKGSETPATH}/pkg/mod/github.com/gogo/protobuf@v1.3.2 --plugin=$${GOPKGSETPATH}/bin/protoc-gen-gogofaster --gogofaster_out . --proto_path $${GOMONOLITH_PATH} $1
endef

GEN_PROTO_FILES = $(patsubst %.proto,%.pb.go,$(shell find . -name *.proto | grep -v ^./vendor))

%.pb.go: %.proto
	$(call PROTOC_GEN,$<)

blueprint/user/models/generatemodels.pb.go: blueprint/user/models/generatemodels.proto
	$(call PROTOC_GEN,blueprint/user/models/generatemodels.proto)

	# This is to allow calling go generate on generatemodels.pb.go
	sed -e 's/DO NOT EDIT./DO NOT MODIFY/' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	sed -e 's/proto "github.com\/gogo\/protobuf\/proto"/proto "github.com\/gogo\/protobuf\/proto"\n"github.com\/go-monolith\/go-monolith\/model"\n/' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	sed -e 's/math_bits "math\/bits"/math_bits "math\/bits"\n"time"\n/' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	sed -e 's/type User struct {/type User struct {\nmodel.Model\n/' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	sed -e 's/UserGroup[ ]*string/UserGroup UserGroup/' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	sed -e 's/UserGroupID[ ]*int32/UserGroupID uint/' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	sed -e 's/LastLogin[ ]*string/LastLogin *time.Time/' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	sed -e 's/ExpiresOn[ ]*string/ExpiresOn *time.Time/' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	sed '1 i //go:generate go run github.com/sergeyglazyrindev/go-monolith/extras/gendecoder' -i github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	gofmt -s -w github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go
	mv github.com/sergeyglazyrindev/go-monolith/blueprint/user/models/generatemodels.pb.go blueprint/user/models/
	rm -rf github.com

.proto: $(GEN_PROTO_FILES)

.PHONY: .proto.touch
.proto.touch:
	@echo $(GEN_PROTO_FILES) | xargs touch

.PHONY: .proto.clean
.proto.clean:
	find . \( -name *.pb.go ! -path './vendor/*' \) -exec rm {} \;

# rm -rf flow/layers/generated.proto