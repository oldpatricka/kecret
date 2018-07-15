package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"encoding/base64"
	"os/exec"
	"strings"
)

func usage() string {
	return fmt.Sprintf("usage: %s FILENAME", os.Args[0])
}

type Secret struct {
	ApiVersion string `yaml:"apiVersion"`
	Data       map[string]string
	Kind       string
	Metadata   map[string]string
	Type       string
}

func exitWithError(msg string) {
	exitWithMessage(msg, 1)
}

func exitWithMessage(msg string, code int) {
	if !strings.HasSuffix(msg, "\n") {
		msg = fmt.Sprintf("%s\n", msg)
	}
	fmt.Fprintf(os.Stderr, msg)
	os.Exit(code)
}

func encodeSecret(es Secret) {
	for k, d := range es.Data {
		v := base64.StdEncoding.EncodeToString([]byte(d))
		es.Data[k] = v
	}
}

func saveBytesToTempFile(decoded []byte) string {
	tempFile, err := ioutil.TempFile("", "kecret")
	if err != nil {
		msg := fmt.Sprintf("Couldn't create temporary file: %s", err)
		exitWithError(msg)
	}
	tempFileName := tempFile.Name()
	tempFile.Write(decoded)
	tempFile.Close()
	return tempFileName
}

func decodeSecretFile(secretFileName string) (Secret, error) {

	s := Secret{}

	rawSecret, err := ioutil.ReadFile(secretFileName)
	if err != nil {
		return s, fmt.Errorf("couldn't read secret file: %s", err)
	}

	err = yaml.Unmarshal(rawSecret, &s)
	if err != nil {
		return s, fmt.Errorf("couldn't unmarshal secret file: %s", err)
	}

	for k, v := range s.Data {
		d, err := base64.StdEncoding.DecodeString(v)
		if err != nil {
			return s, fmt.Errorf("couldn't base64 decode secret %s: %s", k, err)
		}
		s.Data[k] = string(d)
	}
	return s, nil
}

func editFile(tempFileName string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // same thing kubectl does
	}
	cmd := exec.Command(editor, tempFileName)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("couldn't start editor %s: %s", editor, err)
	}
	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("error while editing: %s", err)
	}
	return nil
}

func main() {


	if len(os.Args) != 2 {
		exitWithError(usage())
	}

	secretFileName := os.Args[1]
	s, err := decodeSecretFile(secretFileName)
	if err != nil {
		exitWithError(err.Error())
	}

	decoded, err := yaml.Marshal(&s)
	if err != nil {
		msg := fmt.Sprintf("Couldn't marshal yaml: %s", err)
		exitWithError(msg)
	}

	tempFileName := saveBytesToTempFile(decoded)
	if err != nil {
		exitWithError(err.Error())
	}

	err = editFile(tempFileName)
	if err != nil {
		exitWithError(err.Error())
	}

	editedRawSecret, err := ioutil.ReadFile(tempFileName)
	if err != nil {
		msg := fmt.Sprintf("Couldn't read edited secret: %s", err)
		exitWithError(msg)
		return
	}

	if string(editedRawSecret) == string(decoded) {
		exitWithMessage("No change to secret, aborting...", 0)
	}

	es := Secret{}
	yaml.Unmarshal(editedRawSecret, &es)
	if err != nil {
		msg := fmt.Sprintf("Couldn't unmarshal edited secret: %s", err)
		exitWithError(msg)
	}

	encodeSecret(es)

	encoded, err := yaml.Marshal(&es)
	if err != nil {
		msg := fmt.Sprintf("Couldn't marshal edited secret: %s", err)
		exitWithError(msg)
	}

	err = ioutil.WriteFile(secretFileName, encoded, 0)
	if err != nil {
		msg := fmt.Sprintf("Couldn't write edited secret: %s", err)
		exitWithError(msg)
	}
}