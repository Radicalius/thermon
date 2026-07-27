package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var tempReading = regexp.MustCompile(`^([^:]+):\s+\+?-?([0-9.]+).*$`)

func CollectTempData() {
	config := GetConfig()
	db := GetDatabase()

	for {
		sensors := exec.Command("sensors")
		out, err := sensors.Output()
		if err == nil {
			component := ""
			adapter := ""
			lines := strings.Split(string(out), "\n")
			for _, line := range lines {
				if line == "" {
					continue
				}

				if strings.HasPrefix(line, "Adapter") {
					adapter = strings.Replace(line, "Adapter: ", "", 1)
					continue
				}

				if !strings.Contains(line, ":") && !strings.Contains(line, "(") {
					component = line
					continue
				}

				if !tempReading.MatchString(line) {
					continue
				}

				matches := tempReading.FindStringSubmatch(line)
				if len(matches) < 3 {
					continue
				}

				value, err := strconv.ParseFloat(matches[2], 64)
				if err != nil {
					continue
				}

				ts := time.Now().UnixMilli()
				err = db.Write(SensorReading{
					Timestamp:   ts,
					Device:      component,
					Adapter:     adapter,
					Measurement: matches[1],
					Value:       value,
				})
				if err != nil {
					fmt.Printf("warning: error writing to database: %s", err.Error())
				}
			}
		} else {
			fmt.Printf("error when running sensors command: %s\n", err.Error())
		}

		err = db.ClearExpiredRecords()
		if err != nil {
			fmt.Printf("warning: error when clearing expired records: %s", err.Error())
		}

		time.Sleep(time.Duration(config.SampleInterval) * time.Second)
	}
}
