package scheduled_task

import (
	"goflet/config"
	"log"
	"reflect"
)

var scheduleToCheck = map[string]func(){
	"DeleteEmptyFolder": DeleteEmptyFolder,
	"CleanOutdatedFile": CleanOutdatedFile,
}

// runOneTask runs the task
func runOneTask(name string, task func()) {
	log.Printf("Start running task [%s].", name)
	task()
	log.Printf("Task [%s] completed.", name)
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
			log.Printf("Task [%s] is not scheduled.", name)
			continue
		}
		if value.Int() <= 0 {
			log.Printf("Task [%s] is not scheduled.", name)
			continue
		}

		// Run the task in a goroutine
		log.Printf("Task [%s] is scheduled every %d seconds.", name, value.Int())
		go runOneTask(name, task)
	}
}
