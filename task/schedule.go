// Package task provides functions for running scheduled tasks
package task

import (
	"reflect"
	"time"

	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/util/log"
)

var scheduleToCheck = map[string]func(){
	"DeleteEmptyFolder": DeleteEmptyFolder,
	"CleanOutdatedFile": CleanOutdatedFile,
}

// runTask runs the task
func runTask(name string, task func(), interval int) {
	for {
		log.Infof("Start running task [%s].", name)
		task()
		log.Infof("Task [%s] completed.", name)
		time.Sleep(time.Duration(interval) * time.Second) // Sleep for the interval
	}
}

// RunScheduledTask runs the scheduled tasks
func RunScheduledTask() {
	scheduledTaskConfig := config.GofletCfg.CronConfig
	reflectSchedule := reflect.ValueOf(scheduledTaskConfig) // Get the reflect value of the schedule

	// Check the schedule
	for name, task := range scheduleToCheck {
		if task == nil {
			continue
		}
		value := reflectSchedule.FieldByName(name)
		if !value.IsValid() {
			log.Debugf("Task [%s] is not scheduled.", name)
			continue
		}
		if value.Int() <= 0 {
			log.Debugf("Task [%s] is not scheduled.", name)
			continue
		}

		// Run the task in a goroutine
		log.Infof("Task [%s] is scheduled every %d seconds.", name, value.Int())
		go runTask(name, task, int(value.Int()))
	}
}
