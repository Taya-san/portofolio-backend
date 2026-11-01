package jobs

import "github.com/go-co-op/gocron/v2"

type Job struct {
	JobDef gocron.JobDefinition
	Task   gocron.Task
}
