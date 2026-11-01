package project

import (
	"github.com/Taya-san/portofolio-backend/internal/scheduler/jobs"
	githubproject "github.com/Taya-san/portofolio-backend/internal/usecase/github_project"
	"github.com/go-co-op/gocron/v2"
)

func GithubProjectsDailyRefresher(uc githubproject.GithubProject) jobs.Job {
	return jobs.Job{
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
