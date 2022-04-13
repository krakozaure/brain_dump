package main

import (
	"bytes"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/iancoleman/strcase"
)

type Context map[string]string

func makeContext() Context {
	var err error

	context := make(Context)

	context["app_name"] = APP_NAME

	username, err := user.Current()
	if err == nil {
		context["home"] = username.HomeDir
		context["user"] = username.Username
	}

	cwd, err := os.Getwd()
	if err == nil {
		context["cwd"] = cwd
		context["short_cwd"] = shortenPath(cwd)
	}

	return context
}

func mergeContextMaps(context Context, maps ...Context) Context {
	for _, currentMap := range maps {
		for key, value := range currentMap {
			context[key] = value
		}
	}
	return context
}

func convertContextKeys(context Context, keys_case string) Context {
	converter := strcase.ToSnake
	if keys_case == "CamelCase" {
		converter = strcase.ToCamel
	}

	newContext := make(Context)
	for key, value := range context {
		newContext[converter(key)] = value
	}
	return newContext
}

func getTemplateFuncs() template.FuncMap {
	f := sprig.GenericFuncMap()
	f["shortenPath"] = shortenPath
	f["expandTilde"] = expandTilde
	f["expandEnv"] = os.ExpandEnv
	return f
}

func execTemplate(text string, context interface{}) string {
	writer := bytes.NewBufferString("")

	tpl := template.New("").Option("missingkey=zero")
	tpl.Funcs(getTemplateFuncs())

	tpl, err := tpl.Parse(text)
	if err != nil {
		log.Println(err)
		return text
	}

	err = tpl.Execute(writer, context)
	if err != nil {
		log.Println(err)
		return text
	}

	return writer.String()
}

func execFileTemplate(path string, context interface{}) string {
	path = expandText(path)
	text, err := readFile(path)
	if err != nil {
		log.Println(err)
		return ""
	}
	return execTemplate(text, context)
}

func shortenPath(path string) string {
	username, err := user.Current()
	if err == nil {
		path = strings.Replace(path, username.HomeDir, "~", 1)
	}

	sep := string(filepath.Separator)
	sepsCount := strings.Count(path, sep)
	pathParts := strings.Split(path, sep)

	shortenedPath := ""
	for _, part := range pathParts[:sepsCount] {
		if part == "" && string(path[0]) == sep {
			shortenedPath = sep
		} else if string(part[0]) == "." {
			shortenedPath += string(part[0:2]) + sep
		} else {
			shortenedPath += string(part[0]) + sep
		}
	}
	shortenedPath += pathParts[sepsCount]

	return shortenedPath
}
