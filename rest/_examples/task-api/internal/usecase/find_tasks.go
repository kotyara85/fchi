package usecase

import (
	"context"
	"github.com/swaggest/rest/_examples/task-api/internal/domain/task"
	"github.com/swaggest/usecase"
)

func FindTasks(deps interface {
	TaskFinder() task.Finder
}) usecase.Interactor {
	u := struct {
		usecase.Interactor
		usecase.Info
		usecase.WithInput
		usecase.WithOutput
	}{}

	u.SetTitle("Find Tasks")
	u.SetDescription("Find all tasks.")
	u.Output = new([]task.Entity)
	u.SetTags("Tasks")

	u.Interactor = usecase.Interact(func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*[]task.Entity)
		)

		*out = deps.TaskFinder().Find(ctx)

		return nil
	})

	return u
}
