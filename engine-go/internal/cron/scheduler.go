package cron

import (
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/app"
	"github.com/SeaweedbrainCY/jellyfin-newsletter/internal/newsletter"
	"github.com/go-co-op/gocron/v2"
	"go.uber.org/zap"
)

func CreateNewsletterScheduler(
	newsletterWorkflow newsletter.Workflow,
	app *app.ApplicationContext,
) (gocron.Scheduler, gocron.Job, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		return nil, err
	}

	job, err := scheduler.NewJob(
		gocron.CronJob(app.Config.Scheduler.CronExpr, false),
		gocron.NewTask(newsletterWorkflow.Run, app),
	)
	if err != nil {
		return nil, nil, err
	}

	return scheduler, job, nil
}

func LogNextRun(job gocron.Job, app *app.ApplicationContext) {
	nextRunDatetime, err := job.NextRun()
	if err != nil {
		app.Logger.Debug("Error while computing next run.", zap.Error(err))
		return
	}
	jobNextRunStr := nextRunDatetime.Format("2006-01-02T15:04:05Z07:00")

	app.Logger.Info(
		"Scheduler created.",
		zap.String("job_id", job.ID().String()),
		zap.String("Cron expression", app.Config.Scheduler.CronExpr),
		zap.String("Next run", jobNextRunStr),
	)
}
