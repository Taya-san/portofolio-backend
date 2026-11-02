package project

import (
	"github.com/Taya-san/portofolio-backend/internal/scheduler/job"
	githubproject "github.com/Taya-san/portofolio-backend/internal/usecase/github_project"
	"github.com/go-co-op/gocron/v2"
)

func GithubProjectsDailyRefresher(uc githubproject.GithubProject) job.Job {
	return job.Job{
		JobDef: gocron.DailyJob(
			1,
			gocron.NewAtTimes(
				gocron.NewAtTime(0, 0, 0),
			),
		),
		Task: gocron.NewTask(
			uc.RefreshMirror,
		),
	}
}
