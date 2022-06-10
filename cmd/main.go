package main

import (
	"fmt"
	"helm-maker/chart"
)

func main() {
	apps := chart.InitApps()
	result, err := chart.ChartsFile(apps)
	fmt.Println(result, err)
}
