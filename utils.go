package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"
)

func prettyPrint(data interface{}) {
	prettyData := marshal(data)
	if prettyData != "" {
		fmt.Println(marshal(data))
	}
}

func marshal(v interface{}) string {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to marshal the data")
		return ""
	}
	return string(data)
}

func unmarshal(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

func expandText(text string) string {
	text = os.ExpandEnv(text)
	return expandTilde(text)
}

func expandTilde(text string) string {
	usr, _ := user.Current()
	return strings.ReplaceAll(text, "~/", usr.HomeDir+"/")
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

func fileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func readStdin() []byte {
	input, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal("Unable to read from STDIN")
	}
	return input
}

func readFile(path string) (string, error) {
	path = expandText(path)
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
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
