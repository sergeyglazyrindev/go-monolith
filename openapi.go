package gomonolith

import (
	"encoding/json"
	"fmt"
	"github.com/sergeyglazyrindev/go-monolith/core"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type OpenAPICommand struct {
}

func (c OpenAPICommand) Proceed(subaction string, args []string) error {
	var action string
	var help string
	var isCorrectActionPassed bool = false
	commandRegistry := &core.CommandRegistry{
		Actions: make(map[string]core.ICommand),
	}

	commandRegistry.AddAction("editor", &ServeOpenAPIEditorCommand{})
	if len(os.Args) > 2 {
		action = os.Args[2]
		isCorrectActionPassed = commandRegistry.IsRegisteredCommand(action)
	}
	if !isCorrectActionPassed {
		helpText := commandRegistry.MakeHelpText()
		help = fmt.Sprintf(`
Please provide what do you want to do ?
%s
`, helpText)
		fmt.Print(help)
		return nil
	}
	return commandRegistry.RunAction(subaction, "", args)
}

func (c OpenAPICommand) GetHelpText() string {
	return "Manage your open api schema"
}

type ServeOpenAPIEditorCommand struct {
}

func (command ServeOpenAPIEditorCommand) Proceed(subaction string, args []string) error {
	// Hello world, the web server

	editorHandler := func(w http.ResponseWriter, req *http.Request) {
		io.WriteString(w, `
<!DOCTYPE html><html lang=""><head><meta charset="utf-8"><meta http-equiv="X-UA-Compatible" content="IE=edge"><meta name="viewport" content="width=device-width,initial-scale=1"><link rel="icon" href="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/favicon.ico"><title>oaie-sketch</title><link href="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/css/app.eefa6853.css" rel="preload" as="style"><link href="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/css/chunk-vendors.f31ad7fb.css" rel="preload" as="style"><link href="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/js/app.c18a69dc.js" rel="preload" as="script"><link href="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/js/chunk-vendors.8dfb66d6.js" rel="preload" as="script"><link href="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/css/chunk-vendors.f31ad7fb.css" rel="stylesheet"><link href="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/css/app.eefa6853.css" rel="stylesheet"></head><body><noscript><strong>We're sorry but oaie-sketch doesn't work properly without JavaScript enabled. Please enable it to continue.</strong></noscript><div id="app"></div><script src="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/js/chunk-vendors.8dfb66d6.js"></script><script src="https://cdn.jsdelivr.net/gh/sergeyglazyrindev/oaie-sketch@0.3.8/dist/js/app.c18a69dc.js"></script></body></html>
`)
	}
	specHandler := func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case "GET":
			content, err := ioutil.ReadFile("configs/api-spec.yml")
			if err != nil {
				log.Fatal(err)
			}
			w.Write(content)
		case "POST":
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
				http.Error(w, "can't read body", http.StatusBadRequest)
				return
			}
			ioutil.WriteFile("configs/api-spec.yml", body, 0644)
			w.Header().Set("Content-Type", "application/json")
			data := map[string]string{"ok": "ok"}
			json.NewEncoder(w).Encode(data)
		}
	}

	http.HandleFunc("/", editorHandler)
	http.HandleFunc("/spec/", specHandler)
	fmt.Printf("Please open following url in browser http://localhost:8083/?saveToFileEndpoint=/spec/\n")
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", appInstance.Config.D.Swagger.APIEditorListenPort), nil))
	return nil
}

func (command ServeOpenAPIEditorCommand) GetHelpText() string {
	return "Provide a way to edit your project openapi spec"
}
