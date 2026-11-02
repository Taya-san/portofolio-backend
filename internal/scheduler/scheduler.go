package scheduler

import (
	"context"
	"fmt"

	"github.com/Taya-san/portofolio-backend/internal/scheduler/job"
	"github.com/Taya-san/portofolio-backend/internal/scheduler/job/project"
	githubproject "github.com/Taya-san/portofolio-backend/internal/usecase/github_project"
	"github.com/go-co-op/gocron/v2"
)

func NewScheduler(
	ctx context.Context,
	l gocron.Logger,
	uc githubproject.GithubProject,
) (gocron.Scheduler, error) {
	var jobs []job.Job

	scheduler, err := gocron.NewScheduler(
		gocron.WithGlobalJobOptions(
			gocron.WithContext(ctx),
		),
		gocron.WithLogger(
			l,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("Scheduler - NewScheduler - gocron.NewScheduler: %w", err)
	}

	jobs = append(jobs, project.GithubProjectsDailyRefresher(uc))

	for _, job := range jobs {
		scheduler.NewJob(
			job.JobDef,
			job.Task,
		)
	}

	return scheduler, nil
}
