package core

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

type Microservice struct {
	Name        string
	Prefix      string
	AuthBackend string
	URLPrefix   string
	Port        int
	SwaggerPort int
	IncludeTags []string
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
