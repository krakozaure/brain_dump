package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
)

const (
	APP_NAME        = "brain_dump"
	APP_CONFIG_FILE = "$HOME/.config/brain_dump/brain_dump.json"
)

func main() {
	var (
		err     error
		inpText string
		context map[string]string
	)

	parseFlags()

	if printConfigFlag {
		prettyPrint(getDefaultConfig())
		return
	}

	config := getConfig()
	profile := config.getProfile(profileNameFlag)

	if profile.File == "" {
		log.Fatal("The 'file' configuration key must be set !")
		return
	}
	outFile := expandText(profile.File)

	cliArgs, cliVars := extractCliVars()

	inpText = ""
	if len(cliArgs) > 0 {
		inpText = getInput(cliArgs)
	}
	if len(inpText) == 0 {
		openEditorFlag = true
	}

	cliVars["input"] = inpText
	context = makeContext(profile)
	context = mergeContextMaps(context, profile.TemplateVars, cliVars)
	context = convertContextKeys(context, profile.KeysCase)

	outFile = execTemplate(outFile, context)

	if profile.TemplateFile != "" {
		inpText = execFileTemplate(profile.TemplateFile, context)
	}

	if !strings.HasSuffix(inpText, "\n") {
		inpText += "\n"
	}

	debugLevel, err := strconv.Atoi(os.Getenv("BRAINDUMP_DEBUG"))
	if err == nil && !openEditorFlag && debugLevel > 1 {
		fmt.Println("Profile :")
		prettyPrint(profile)
		fmt.Println("Context :")
		prettyPrint(context)
	}
	if err == nil && !openEditorFlag && debugLevel > 0 {
		fmt.Println("Output file :", outFile)
		fmt.Println(strings.Repeat("-", 55))
		fmt.Println(inpText)
	}

	if profile.FileMode == "write" {
		err = writeFile(outFile, inpText)
	} else {
		err = appendFile(outFile, inpText)
	}
	if err != nil {
		log.Fatal(err)
	}

	if !openEditorFlag {
		return
	}

	if profile.Editor == "" {
		log.Fatal("The 'editor' configuration key must be set to edit files")
		return
	}
	err = editFile(profile.Editor, profile.EditorArgs, outFile)
	if err != nil {
		log.Fatal(err)
	}
}

func getConfig() Config {
	config, err := getUserConfig()
	if err != nil {
		log.Println("Unable to load the configuration file. Default configuration is used.")
		config = getDefaultConfig()
	}

	configString := expandText(config.String())
	err = unmarshalString(configString, &config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func getInput(cliArgs []string) string {
	var data string
	if cliArgs[0] == "-" {
		data = string(readStdin())
	} else {
		data = strings.Join(cliArgs, " ")
	}
	return data
}

func mergeContextMaps(context map[string]string, maps ...map[string]string) map[string]string {
	for _, currentMap := range maps {
		for key, value := range currentMap {
			context[key] = value
		}
	}
	return context
}

func convertContextKeys(context map[string]string, keys_case string) map[string]string {
	converter := strcase.ToSnake
	if keys_case == "CamelCase" {
		converter = strcase.ToCamel
	}

	newContext := make(map[string]string)
	for key, value := range context {
		newContext[converter(key)] = value
	}
	return newContext
}
