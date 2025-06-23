package service

import (
	"errors"
	"pet-project/internal/repository"
	"pet-project/pkg/model"
	"time"
)

type TaskService struct {
	Repository repository.TaskRepository
}

func (s *TaskService) CreateTask(task *model.Task) error {

	if task.Title == "" {
		return errors.New("Title is required")
	}
	if task.Priority == "" {
		return errors.New("Priority is required")
	}
	if task.Status == "" {
		return errors.New("Status is required")
	}

	now := time.Now()
	task.CreatedAt = now
	task.UpdatedAt = now

	err := s.Repository.CreateTask(task)
	if err != nil {
		return err
	}
	return nil
}

func (s *TaskService) UpdateTask(task *model.Task, user_id int) error {
	if task.CreatedAt == task.UpdatedAt {
		return errors.New("The Task should update!")
	}
	task.UpdatedAt = time.Now()
	if task.Title == "" {
		return errors.New("Title is required")
	}
	if task.Priority == "" {
		return errors.New("Priority is required")
	}
	if task.Status == "" {
		return errors.New("Status is required")
	}

	if task.AssignedTo != user_id {
		return errors.New("No permission to update task")
	}

	err := s.Repository.UpdateTask(task)
	if err != nil {
		return err
	}

	return nil
}

func (s *TaskService) DeleteTask(task_id int, user_id int) error {
	task, err := s.Repository.GetByIDTask(task_id)
	if err != nil {
		return err
	}

	if task.AssignedTo != user_id {
		return errors.New("No permission to delete task")
	}

	if task.Status == "done" {
		return errors.New("can't delete complete task")
	}

	return s.Repository.DeleteTask(task_id)
}

func (s *TaskService) ListByProjectTask(project_id int) ([]*model.Task, error) {
	tasks, err := s.Repository.ListByProjectTask(project_id)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *TaskService) GetByIDTask(task_id int, user_id int) (*model.Task, error) {
	task, err := s.Repository.GetByIDTask(task_id)
	if task.Status == "done" {
		return nil, errors.New("Task just is already")
	}
	if err != nil {
		return nil, err
	}
	return task, nil

}

func (s *TaskService) MarkTaskFinished(task *model.Task) error {
	due_date := time.Now()
	task.DueDate = &due_date
	task.Status = "done"
	task.Priority = "NaN"
	task.UpdatedAt = time.Now()
	return nil
}
