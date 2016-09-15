// Copyright 2015 Keybase, Inc. All rights reserved. Use of
// this source code is governed by the included BSD license.

package launchd

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"
)

// Env is environment for launchd
type Env struct {
	Home         string
	GoPath       string
	GithubToken  string
	SlackToken   string
	SlackChannel string
	AWSAccessKey string
	AWSSecretKey string
}

// Script is what to run
type Script struct {
	Label      string
	Path       string
	Command    string
	BucketName string
	Platform   string
	EnvVars    []EnvVar
}

// EnvVar is custom env vars
type EnvVar struct {
	Key   string
	Value string
}

type job struct {
	Env    Env
	Script Script
}

const plistTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>{{.Script.Label }}</string>
    <key>EnvironmentVariables</key>
    <dict>
        <key>GOPATH</key>
        <string>{{ .Env.GoPath }}</string>
        <key>GITHUB_TOKEN</key>
        <string>{{ .Env.GithubToken }}</string>
        <key>SLACK_TOKEN</key>
        <string>{{ .Env.SlackToken }}</string>
        <key>SLACK_CHANNEL</key>
        <string>{{ .Env.SlackChannel }}</string>
        <key>AWS_ACCESS_KEY</key>
        <string>{{ .Env.AWSAccessKey }}</string>
        <key>AWS_SECRET_KEY</key>
        <string>{{ .Env.AWSSecretKey }}</string>
        <key>PATH</key>
        <string>/sbin:/usr/sbin:/bin:/usr/bin:/usr/local/bin</string>
        <key>LOG_PATH</key>
        <string>{{ .Env.Home }}/Library/Logs/{{ .Script.Label }}.log</string>
        <key>BUCKET_NAME</key>
        <string>{{ .Script.BucketName }}</string>
        <key>SCRIPT_PATH</key>
        <string>{{ .Env.GoPath }}/src/{{ .Script.Path }}</string>
        <key>PLATFORM</key>
        <string>{{ .Script.Platform }}</string>
        <key>COMMAND</key>
        <string>{{ .Script.Command }}</string>
        {{ with .Script.EnvVars }}{{ range . }}
        <key>{{ .Key }}</key>
        <string>{{ .Value }}</string>
        {{ end }}{{ end }}
    </dict>
    <key>ProgramArguments</key>
    <array>
        <string>/bin/bash</string>
        <string>{{ .Env.GoPath }}/src/github.com/keybase/slackbot/launchd/run.sh</string>
    </array>
    <key>StandardErrorPath</key>
    <string>{{ .Env.Home }}/Library/Logs/{{ .Script.Label }}.log</string>
    <key>StandardOutPath</key>
    <string>{{ .Env.Home }}/Library/Logs/{{ .Script.Label }}.log</string>
</dict>
</plist>
`

// NewEnv creates environment
func NewEnv() Env {
	return Env{
		Home:         os.Getenv("HOME"),
		GoPath:       os.Getenv("GOPATH"),
		GithubToken:  os.Getenv("GITHUB_TOKEN"),
		SlackToken:   os.Getenv("SLACK_TOKEN"),
		SlackChannel: os.Getenv("SLACK_CHANNEL"),
		AWSAccessKey: os.Getenv("AWS_ACCESS_KEY"),
		AWSSecretKey: os.Getenv("AWS_SECRET_KEY"),
	}
}

// Plist is plist for env and args
func (e Env) Plist(script Script) ([]byte, error) {
	t := template.New("Plist template")
	j := job{Env: e, Script: script}
	t, err := t.Parse(plistTemplate)
	if err != nil {
		return nil, err
	}
	buff := bytes.NewBufferString("")
	err = t.Execute(buff, j)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

// WritePlist writes out plist
func (e Env) WritePlist(script Script) error {
	data, err := e.Plist(script)
	if err != nil {
		return err
	}
	path := fmt.Sprintf("%s/Library/Launch Agents/%s.plist", e.Home, script.Label)
	return ioutil.WriteFile(path, data, 0755)
}
