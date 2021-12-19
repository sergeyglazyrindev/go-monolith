package core

import (
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type ServiceSwaggerDefinitionInfoContact struct {
	Name  string `json:"name"`
	URL   string `json:"url"`
	Email string `json:"email"`
}

type ServiceSwaggerDefinitionInfoLicense struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type ServiceSwaggerDefinitionInfo struct {
	Contact *ServiceSwaggerDefinitionInfoContact `json:"contact"`
	License *ServiceSwaggerDefinitionInfoLicense `json:"license"`
	Version string                               `json:"version"`
}

type ServiceSwaggerDefinition struct {
	Consumes       []string                      `json:"consumes"`
	Produces       []string                      `json:"produces"`
	Schemes        []string                      `json:"schemes"`
	SwaggerVersion string                        `json:"swagger"`
	Info           *ServiceSwaggerDefinitionInfo `json:"info"`
	Host           string                        `json:"host"`
	BasePath       string                        `json:"basePath"`
}

type Microservice struct {
	Name                     string
	Prefix                   string
	AuthBackend              string
	URLPrefix                string
	Port                     int
	SwaggerPort              int
	IncludeTags              []string
	ServiceSwaggerDefinition *ServiceSwaggerDefinition
}

func NewServiceSwaggerDefinition() *ServiceSwaggerDefinition {
	return &ServiceSwaggerDefinition{
		Schemes: make([]string, 0), Info: &ServiceSwaggerDefinitionInfo{
			Contact: &ServiceSwaggerDefinitionInfoContact{},
			License: &ServiceSwaggerDefinitionInfoLicense{},
		},
		Consumes: make([]string, 0),
		Produces: make([]string, 0),
	}
}

func (m Microservice) RegisterEndpoints(app IApp) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	if m.AuthBackend != "" {
		authAdapter, _ := app.GetAuthAdapterRegistry().GetAdapter(m.AuthBackend)
		adapterGroup := r.Group("/" + authAdapter.GetName())
		adapterGroup.POST("/signin/", authAdapter.Signin)
		adapterGroup.POST("/signup/", authAdapter.Signup)
		adapterGroup.POST("/logout/", authAdapter.Logout)
		adapterGroup.GET("/status/", authAdapter.IsAuthenticated)
	}
	return r
}

func (m Microservice) Start(r *gin.Engine) {
	r.Run(fmt.Sprintf(":%d", m.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

func (m Microservice) StartSwagger(app IApp) error {
	file, err := ioutil.TempFile("", "swagger.*.json")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer os.Remove(file.Name())
	params := []string{"generate", "spec", "-o", file.Name()}
	controllerSuffixes := make([]string, 0)
	if m.AuthBackend != "" {
		controllerSuffixes = append(
			controllerSuffixes, m.AuthBackend,
		)
	}
	if m.Prefix != "" {
		controllerSuffixes = append(controllerSuffixes, m.Prefix)
	}
	if len(controllerSuffixes) > 0 {
		for _, controllerSuffix := range controllerSuffixes {
			params = append(params, fmt.Sprintf("--include-tag=%s", controllerSuffix))
		}
	}
	if m.IncludeTags != nil && len(m.IncludeTags) > 0 {
		for _, includeTag := range m.IncludeTags {
			params = append(params, fmt.Sprintf("--include-tag=%s", includeTag))
		}
	}
	commandToExecute := exec.Command(
		"swagger", params...,
	)
	stderr, err := commandToExecute.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}
	if err := commandToExecute.Start(); err != nil {
		log.Fatal(err)
		return err
	}

	slurp, _ := io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := commandToExecute.Wait(); err != nil {
		log.Fatal(err)
		return err
	}
	if m.ServiceSwaggerDefinition != nil {
		currentServiceDefinition := NewServiceSwaggerDefinition()
		if m.ServiceSwaggerDefinition.Produces != nil {
			currentServiceDefinition.Produces = m.ServiceSwaggerDefinition.Produces
		}
		if m.ServiceSwaggerDefinition.Consumes != nil {
			currentServiceDefinition.Consumes = m.ServiceSwaggerDefinition.Consumes
		}
		if m.ServiceSwaggerDefinition.Schemes != nil {
			currentServiceDefinition.Schemes = m.ServiceSwaggerDefinition.Schemes
		}
		if m.ServiceSwaggerDefinition.SwaggerVersion != "" {
			currentServiceDefinition.SwaggerVersion = m.ServiceSwaggerDefinition.SwaggerVersion
		}
		if m.ServiceSwaggerDefinition.Info != nil {
			if m.ServiceSwaggerDefinition.Info.Contact != nil {
				if m.ServiceSwaggerDefinition.Info.Contact.Name != "" {
					currentServiceDefinition.Info.Contact.Name = m.ServiceSwaggerDefinition.Info.Contact.Name
				}
				if m.ServiceSwaggerDefinition.Info.Contact.URL != "" {
					currentServiceDefinition.Info.Contact.URL = m.ServiceSwaggerDefinition.Info.Contact.URL
				}
				if m.ServiceSwaggerDefinition.Info.Contact.Email != "" {
					currentServiceDefinition.Info.Contact.Email = m.ServiceSwaggerDefinition.Info.Contact.Email
				}
			}
			if m.ServiceSwaggerDefinition.Info.License != nil {
				if m.ServiceSwaggerDefinition.Info.License.Name != "" {
					currentServiceDefinition.Info.License.Name = m.ServiceSwaggerDefinition.Info.License.Name
				}
				if m.ServiceSwaggerDefinition.Info.License.URL != "" {
					currentServiceDefinition.Info.License.URL = m.ServiceSwaggerDefinition.Info.License.URL
				}
			}
			if m.ServiceSwaggerDefinition.Info.Version != "" {
				currentServiceDefinition.Info.Version = m.ServiceSwaggerDefinition.Info.Version
			}
		}
		if m.ServiceSwaggerDefinition.Host != "" {
			currentServiceDefinition.Host = m.ServiceSwaggerDefinition.Host
		}
		if m.ServiceSwaggerDefinition.BasePath != "" {
			currentServiceDefinition.BasePath = m.ServiceSwaggerDefinition.BasePath
		}
		filePython, err1 := ioutil.TempFile("", "swagger-transform*.py")
		if err1 != nil {
			log.Fatal(err1)
			return err1
		}
		fileWithUpdatedDefinition, err1 := ioutil.TempFile("", "updated-definition*.json")
		filePython.Write([]byte(`import json
import sys
def replaceServiceDefinition():
	updated_definition = json.loads(open(sys.argv[2], 'r').read())
	file_with_swagger_json = json.loads(open(sys.argv[1], 'r').read())
	file_with_swagger_json["consumes"] = updated_definition["consumes"]
	file_with_swagger_json["produces"] = updated_definition["produces"]
	file_with_swagger_json["schemes"] = updated_definition["schemes"]
	file_with_swagger_json["swagger"] = updated_definition["swagger"]
	file_with_swagger_json["info"] = updated_definition["info"]
	file_with_swagger_json["host"] = updated_definition["host"]
	file_with_swagger_json["basePath"] = updated_definition["basePath"]
	open(sys.argv[1], 'w').write(json.dumps(file_with_swagger_json))
if __name__ == '__main__':
	replaceServiceDefinition()
`))
		commandToExecute = exec.Command(
			"python", filePython.Name(), file.Name(), fileWithUpdatedDefinition.Name(),
		)
		generatedJSON, _ := json.Marshal(currentServiceDefinition)
		fileWithUpdatedDefinition.Write(generatedJSON)
		commandToExecute.Stderr = os.Stderr
		if err1 = commandToExecute.Start(); err1 != nil {
			log.Fatal(err1)
			return err1
		}
		if err1 = commandToExecute.Wait(); err1 != nil {
			log.Fatal(err1)
			return err1
		}
		spew.Dump("dsadasdas", file.Name())
	}
	commandToExecute = exec.Command(
		"swagger", "serve", "--flavor=swagger", "--no-open",
		fmt.Sprintf("--port=%d", m.SwaggerPort), file.Name(),
	)
	stderr, err = commandToExecute.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return err
	}

	fmt.Printf("Please open following url in browser http://localhost:%d/docs\n", m.SwaggerPort)
	if err := commandToExecute.Start(); err != nil {
		log.Fatal(err)
		return err
	}

	slurp, _ = io.ReadAll(stderr)
	fmt.Printf("%s\n", slurp)

	if err := commandToExecute.Wait(); err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}
