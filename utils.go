package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Masterminds/sprig"
)

func prettyPrint(data interface{}) {
	prettyData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Println(data)
	}
	fmt.Println(string(prettyData))
}

func readFile(path string) (string, error) {
	path = expandText(path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func marshal(v interface{}) string {
	data, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(data)
}

func unmarshalString(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func updateMaps(baseMap map[string]string, maps ...map[string]string) {
	for _, currentMap := range maps {
		for key, value := range currentMap {
			baseMap[key] = value
		}
	}
}

func expandText(text string) string {
	text = os.ExpandEnv(text)
	text = expandTilde(text)
	return text
}

func expandTilde(text string) string {
	usr, _ := user.Current()
	dir := usr.HomeDir
	return strings.ReplaceAll(text, "~/", dir+"/")
}

func makeContext(profile Profile) map[string]string {
	var err error

	context := make(map[string]string)

	context["app_name"] = APP_NAME

	now := time.Now()
	for name, format := range profile.Formats {
		context[name] = now.Format(format)
	}

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

func getTemplateFuncs() template.FuncMap {
	f := sprig.GenericFuncMap()
	f["shortenPath"] = shortenPath
	f["expandTilde"] = expandTilde
	f["expandEnv"] = os.ExpandEnv
	return f
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
		} else {
			shortenedPath += string(part[0]) + sep
		}
	}
	shortenedPath += pathParts[sepsCount]

	// fmt.Printf(
	// 	"path : %v, sepsCount: %v, pathParts : %v, len(pathParts) : %v, shortenedPath : %v\n",
	// 	path, sepsCount, pathParts, len(pathParts), shortenedPath,
	// )

	return shortenedPath
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

func ensureParentDirs(file string) error {
	parents := filepath.Dir(file)
    err := os.MkdirAll(parents, 0755)
    if err == nil {
        return nil
    }
    if os.IsExist(err) {
        info, err := os.Stat(parents)
        if err != nil {
            return err
        }
        if !info.IsDir() {
			return fmt.Errorf("Path exists but is not a directory : '%s'", parents)
        }
        return nil
    }
    return err
}

func readStdin() []byte {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Unable to read from STDIN")
	}
	return input
}

func appendFile(file string, data string) error {
	file = expandText(file)

	f, err := os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write([]byte(data))
	return err
}

func writeFile(file, data string) error {
	file = expandText(file)
	return ioutil.WriteFile(file, []byte(data), 0644)
}

func editFile(program string, arguments []string, file string) error {
	program = os.ExpandEnv(program)
	arguments = append(arguments, expandText(file))

	err := openProgram(program, arguments...)
	if err != nil {
		return err
	}
	return nil
}

func openProgram(program string, arguments ...string) error {
	cmd := exec.Command(program, arguments...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
