package test

import (
	"fmt"
	"helm-maker/chart"
	"testing"
)

func TestGenLocalChart(t *testing.T) {
	h, err := chart.NewHelm()
	if err != nil {
		fmt.Println(err)
		return
	}
	c, err := h.GenLocalChart("test")
	fmt.Printf("########%+v , %s", c, err)
}

func TestGenChartsFile(t *testing.T) {
	apps := chart.InitApps()
	result, err := chart.ChartsFile(apps)
	fmt.Println(result, err)
}
