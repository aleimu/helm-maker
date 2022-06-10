package chart

import (
	"fmt"
	"github.com/pkg/errors"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/chartutil"
	"helm.sh/helm/v3/pkg/helmpath"
	"path/filepath"
)

type createOptions struct {
	starter    string // --starter
	name       string
	starterDir string
}

func (c *Helm) GenLocalChart(chartName string) (*chart.Chart, error) {
	o := &createOptions{}
	o.name = chartName
	o.starterDir = helmpath.DataPath("starters")
	chartname := filepath.Base(o.name)
	chartpath := filepath.Dir(o.name)
	cfile := &chart.Metadata{
		Name:        chartname,
		Description: "A Helm chart for Kubernetes",
		Type:        "application",
		Version:     "0.1.0",
		AppVersion:  "0.1.0",
		APIVersion:  chart.APIVersionV2,
	}
	fmt.Println("----------1")

	if o.starter != "" {
		// Create from the starter
		lstarter := filepath.Join(o.starterDir, o.starter)
		// If path is absolute, we don't want to prefix it with helm starters folder
		if filepath.IsAbs(o.starter) {
			lstarter = o.starter
		}
		err := chartutil.CreateFrom(cfile, chartpath, lstarter)
		if err != nil {
			return nil, err
		}
	}
	//创建默认
	_, err := chartutil.Create(chartname, chartpath)
	if err != nil {
		return nil, fmt.Errorf("Create err:%s", err)
	}
	fmt.Println("----------2", chartpath)
	// 加载

	helmChart, err := loader.Load(chartname)
	if err != nil {
		//panic(err)
		return nil, fmt.Errorf("Load err:%s", err)
	}

	if helmChart.Metadata.Deprecated {
		return nil, errors.New("deprecated chart")
	}
	fmt.Println("----------3")
	// 修改数据
	helmChart.Metadata.Version = "0.2.0"
	//helmChart.Metadata.Name = "just-test"
	//fmt.Printf("--------%+v", helmChart.Values)


	fmt.Println("----------4")
	// 保存为tgz压缩包
	//_, err = chartutil.Save(helmChart, chartpath)
	//if err != nil {
	//	return nil, fmt.Errorf("Save err:%s", err)
	//}
	// 保存为目录
	err = chartutil.SaveDir(helmChart, ".")
	if err != nil {
		return nil, fmt.Errorf("SaveDir err:%s", err)
	}

	return helmChart, err
}

func warpChartCreate() {

}
