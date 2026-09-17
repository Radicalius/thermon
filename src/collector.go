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

var freeReading = regexp.MustCompile(`SwapUse\s+([0-9]+)\s+CachUse\s+([0-9]+)\s+MemUse\s+([0-9]+)\s+MemFree\s+([0-9]+)`)
var freeAdapters []string = []string{"swap", "ram", "ram", "ram"}
var freeMeasurements []string = []string{"swapUse", "cacheUse", "memUse", "memFree"}

func CollectTempData() {
	config := GetConfig()
	db := GetDatabase()

	for {
		ts := time.Now().UnixMilli()

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

		free := exec.Command("free", "-L")
		out, err = free.Output()
		if err == nil {
			matches := freeReading.FindStringSubmatch(string(out))

			if len(matches) < 5 {
				fmt.Printf("error parsing free output: not enough submatches\n")
			}

			for i := 0; i < 4; i++ {
				val, err := strconv.ParseInt(matches[i+1], 10, 64)
				if err != nil {
					fmt.Printf("error parsing %s from free: %s\n", freeMeasurements[i], err.Error())
				}

				err = db.Write(SensorReading{
					Timestamp:   ts,
					Device:      "memory",
					Adapter:     freeAdapters[i],
					Measurement: freeMeasurements[i],
					Value:       float64(val),
				})
				if err != nil {
					fmt.Printf("error writing swap data to db: %s\n", err.Error())
				}
			}
		} else {
			fmt.Printf("error when running free command: %s\n", err.Error())
		}

		err = db.ClearExpiredRecords()
		if err != nil {
			fmt.Printf("warning: error when clearing expired records: %s", err.Error())
		}

		time.Sleep(time.Duration(config.SampleInterval) * time.Second)
	}
}
