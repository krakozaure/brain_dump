package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	APP_NAME        = "brain_dump"
	APP_CONFIG_FILE = "$HOME/.config/brain_dump/brain_dump.json"
)

func main() {
	var err error

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

	timeVars := setTimeVars(profile.Formats)

	cliArgs, cliVars := extractCliVars()
	inpText := ""
	if len(cliArgs) > 0 {
		inpText = getInput(cliArgs)
	}
	if len(inpText) == 0 {
		openEditorFlag = true
	}
	cliVars["input"] = inpText

	context := makeContext()
	context = mergeContextMaps(context, timeVars, profile.TemplateVars, cliVars)
	context = convertContextKeys(context, profile.KeysCase)

	if profile.TemplateFile != "" {
		inpText = execFileTemplate(profile.TemplateFile, context)
	}
	if !strings.HasSuffix(inpText, "\n") {
		inpText += "\n"
	}

	outFile := expandText(profile.File)
	outFile = execTemplate(outFile, context)
	ensureParentDirs(outFile)

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

	if profile.WriteMode == "write" {
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

func getInput(cliArgs []string) string {
	var data string
	if cliArgs[0] == "-" {
		data = string(readStdin())
	} else {
		data = strings.Join(cliArgs, " ")
	}
	return data
}

func setTimeVars(formats map[string]string) map[string]string {
	timeVars := make(map[string]string)
	now := time.Now()
	for name, format := range formats {
		timeVars[name] = now.Format(format)
	}
	return timeVars
}
