/*
Copyright The Helm Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package chart

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"sigs.k8s.io/yaml"
)

// chartName is a regular expression for testing the supplied name of a chart.
// This regular expression is probably stricter than it needs to be. We can relax it
// somewhat. Newline characters, as well as $, quotes, +, parens, and % are known to be
// problematic.
var chartName = regexp.MustCompile("^[a-zA-Z0-9._-]+$")

const (
	// ChartfileName is the default Chart file name.
	ChartfileName = "Chart.yaml"
	// ValuesfileName is the default values file name.
	ValuesfileName = "values.yaml"
	// SchemafileName is the default values schema file name.
	SchemafileName = "values.schema.json"
	// TemplatesDir is the relative directory name for templates.
	TemplatesDir = "templates"
	// ChartsDir is the relative directory name for charts dependencies.
	ChartsDir = "charts"
	// TemplatesTestsDir is the relative directory name for tests.
	TemplatesTestsDir = TemplatesDir + sep + "tests"
	// IgnorefileName is the name of the Helm ignore file.
	IgnorefileName = ".helmignore"
	// IngressFileName is the name of the example ingress file.
	IngressFileName = TemplatesDir + sep + "ingress.yaml"
	// DeploymentName is the name of the example deployment file.
	DeploymentName = TemplatesDir + sep + "deployment.yaml"
	// ServiceName is the name of the example service file.
	ServiceName = TemplatesDir + sep + "service.yaml"
	// ServiceAccountName is the name of the example serviceaccount file.
	ServiceAccountName = TemplatesDir + sep + "serviceaccount.yaml"
	// HorizontalPodAutoscalerName is the name of the example hpa file.
	HorizontalPodAutoscalerName = TemplatesDir + sep + "hpa.yaml"
	// NotesName is the name of the example NOTES.txt file.
	NotesName = TemplatesDir + sep + "NOTES.txt"
	// HelpersName is the name of the example helpers file.
	HelpersName = string(filepath.Separator) + "_helpers.tpl"
	// TestConnectionName is the name of the example test file.
	TestConnectionName = TemplatesTestsDir + sep + "test-connection.yaml"
	APPNAME            = "<APPNAME>"
	CHARTNAME          = "<CHARTNAME>"
)

// maxChartNameLength is lower than the limits we know of with certain file systems,
// and with certain Kubernetes fields.
const maxChartNameLength = 250

const sep = string(filepath.Separator) + "%s_"

// compatibility.
var Stderr io.Writer = os.Stderr

// CreateFrom creates a new chart, but scaffolds it from the src chart.
func CreateFrom(chartfile *chart.Metadata, dest, src string) error {
	schart, err := loader.Load(src)
	if err != nil {
		return errors.Wrapf(err, "could not load %s", src)
	}

	schart.Metadata = chartfile

	var updatedTemplates []*chart.File

	for _, template := range schart.Templates {
		newData := transform(string(template.Data), schart.Name())
		updatedTemplates = append(updatedTemplates, &chart.File{Name: template.Name, Data: newData})
	}

	schart.Templates = updatedTemplates
	b, err := yaml.Marshal(schart.Values)
	if err != nil {
		return errors.Wrap(err, "reading values file")
	}

	var m map[string]interface{}
	if err := yaml.Unmarshal(transform(string(b), schart.Name()), &m); err != nil {
		return errors.Wrap(err, "transforming values file")
	}
	schart.Values = m

	// SaveDir looks for the file values.yaml when saving rather than the values
	// key in order to preserve the comments in the YAML. The name placeholder
	// needs to be replaced on that file.
	for _, f := range schart.Raw {
		if f.Name == ValuesfileName {
			f.Data = transform(string(f.Data), schart.Name())
		}
	}

	return chartutil.SaveDir(schart, dest)
}

// transform performs a string replacement of the specified source for
// a given key with the replacement string
func transform(src, replacement string) []byte {
	return []byte(strings.ReplaceAll(src, "<CHARTNAME>", replacement))
}

func writeFile(name string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	return ioutil.WriteFile(name, content, 0644)
}

func writeFileAppend(name string, content []byte) error {
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		return err
	}
	file, err := os.OpenFile(name, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("文件打开失败", err)
		return err
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	write.Write(content)
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	return err
}

func validateChartName(name string) error {
	if name == "" || len(name) > maxChartNameLength {
		return fmt.Errorf("chart name must be between 1 and %d characters", maxChartNameLength)
	}
	if !chartName.MatchString(name) {
		return fmt.Errorf("chart name must match the regular expression %q", chartName.String())
	}
	return nil
}

type TemplateModel struct {
	FileName string // deployment_%s.yaml
	Content  string // yaml
	Type     string // deployment/svc
	GetName  func() string
}

type Model map[string]TemplateModel

var model = Model{
	"deployment": TemplateModel{
		FileName: "deployment_%s.yaml",
		Content:  defaultDeployment,
		Type:     "deployment",
	},
	"svc": TemplateModel{
		FileName: "svc_%s.yaml",
		Content:  defaultService,
		Type:     "svc",
	},
	"service": TemplateModel{
		FileName: "service_%s.yaml",
		Content:  defaultService,
		Type:     "service",
	},
	"pv":  TemplateModel{},
	"pvc": TemplateModel{},
	"set": TemplateModel{},
}

// 单个应用
type App struct {
	Name   string
	Types  []string
	Values map[string]interface{}
}

// 组合应用
type Apps struct {
	Name    string
	Path    string
	Sets    []*App
	Version string
}

// 构建单应用的部署文件
func WriteTplFile(apps *Apps, app *App, dir string) (string, error) {
	if err := validateChartName(app.Name); err != nil {
		return "", err
	}

	path, err := filepath.Abs(dir)
	if err != nil {
		return path, err
	}

	if fi, err := os.Stat(path); err != nil {
		return path, err
	} else if !fi.IsDir() {
		return path, errors.Errorf("no such directory %s", path)
	}

	for _, t := range app.Types {
		m, ok := model[t]
		if !ok {
			continue
		}
		path := filepath.Join(path, fmt.Sprintf(m.FileName, app.Name))
		if _, err := os.Stat(path); err == nil {
			// There is no handle to a preferred output stream here.
			fmt.Fprintf(Stderr, "WARNING: File %q already exists. Overwriting.\n", path)
		}
		if strings.Contains(m.Content, APPNAME) {
			m.Content = strings.ReplaceAll(m.Content, APPNAME, app.Name)
		}
		if strings.Contains(m.Content, CHARTNAME) {
			m.Content = strings.ReplaceAll(m.Content, CHARTNAME, apps.Name)
		}
		if err := writeFile(path, []byte(m.Content)); err != nil {
			return path, err
		}
	}
	return "", err
}

// 构建value.yaml文件
func WriteValueFile(path string, apps *Apps, defaultContent []byte) error {
	var err error
	var content [][]byte
	appValue := make(map[string]interface{})
	if defaultContent != nil {
		content = append(content, defaultContent)
	}
	for _, app := range apps.Sets {
		appValue[app.Name] = app.Values
	}
	value, err := yaml.Marshal(appValue)
	content = append(content, value)
	if err := writeFile(filepath.Join(path, ValuesfileName), bytes.Join(content, []byte("\n"))); err != nil {
		fmt.Println("create Value.yaml err:", err)
	}
	return err
}

// 构建_helpers.tpl
func WriteHelperFile(path, defaultHelpers string, app *App, apps *Apps) error {
	var err error
	content := strings.ReplaceAll(defaultHelpers, APPNAME, app.Name)
	if err := writeFileAppend(path, transform(content, apps.Name)); err != nil {
		fmt.Println("create Helpers.tpl err:", err)
	}
	return err
}

// 构建多个应用的部署文件
func ChartsFile(apps *Apps) (string, error) {

	if err := validateChartName(apps.Name); err != nil {
		return "", err
	}

	path, err := filepath.Abs(apps.Path)
	if err != nil {
		return path, err
	}

	if fi, err := os.Stat(path); err != nil {
		return path, err
	} else if !fi.IsDir() {
		return path, errors.Errorf("no such directory %s", path)
	}

	cdir := filepath.Join(path, apps.Name)
	if fi, err := os.Stat(cdir); err == nil && !fi.IsDir() {
		return cdir, errors.Errorf("file %s already exists and is not app directory", cdir)
	}
	templatesDir := filepath.Join(cdir, TemplatesDir)
	// create helpers.tpl for all templates
	if err := writeFile(templatesDir+HelpersName, transform(defaultHelpers, apps.Name)); err != nil {
		fmt.Println("chartsFile Chart.yaml err:", err)
	}
	for _, app := range apps.Sets {
		// create helpers.tpl for app templates
		WriteHelperFile(templatesDir+HelpersName, defaultAppHelpers, app, apps)
		_, err := WriteTplFile(apps, app, templatesDir)
		if err != nil {
			fmt.Println("chartsFile err:", err)
			break
		}
	}

	// create value.yaml
	WriteValueFile(cdir, apps, nil)
	// Chart.yaml
	if err := writeFile(filepath.Join(cdir, ChartfileName), transform(fmt.Sprintf(defaultChartfile, apps.Name), apps.Name)); err != nil {
		fmt.Println("chartsFile Chart.yaml err:", err)
	}
	return "", err
}

func InitApps() *Apps {
	apps := new(Apps)
	apps.Name = "demo"
	apps.Version = "1.0.0"
	apps.Path = "."
	var n map[string]interface{}
	const defaultValuesJson = `
{
  "appname": "appname",
  "version": "version",
  "value": {
    "replicaCount": 1,
    "image": {
      "repository": "nginx",
      "pullPolicy": "IfNotPresent",
      "tag": ""
    },
    "imagePullSecrets": [ ],
    "nameOverride": "",
    "fullnameOverride": "",
    "podAnnotations": { },
    "podSecurityContext": { },
    "securityContext": { },
    "service": {
      "type": "ClusterIP",
      "port": 80
    },
    "resources": { },
    "nodeSelector": { },
    "affinity": { },
    "secret": "",
    "configMap": "",
    "volumeMounts": [ ],
    "volumes": [ ],
    "env": [
      {
        "name": "APPNAME",
        "value": "dragon-claw"
      },
      {
        "name": "APP_EXT_PARAM"
      },
      {
        "name": "APP_PORT",
        "value": "8088"
      },
      {
        "name": "APP_RUN_MODE",
        "value": "fg"
      }
    ]
  }
}
	`
	json.Unmarshal([]byte(defaultValuesJson), &n)

	app1 := &App{
		Name:   "app1",
		Types:  []string{"deployment", "svc"},
		Values: n,
	}
	app2 := &App{
		Name:   "app2",
		Types:  []string{"deployment", "svc"},
		Values: n,
	}
	app3 := &App{
		Name:   "app3",
		Types:  []string{"deployment", "svc"},
		Values: n,
	}
	apps.Sets = []*App{app1, app2, app3}
	return apps

}
