package service

import (
	"errors"
	"pet-project/internal/repository"
	"pet-project/pkg/model"
	"time"
)

type CommentsService struct {
	Repository repository.CommentsRepository
}

func (s *CommentsService) AddComment(com *model.Comments) error {
	if com.TaskID <= 0 {
		return errors.New("Task ID is not valid")
	}

	if com.Text == "" {
		return errors.New("Comment required text")
	}

	com.CreatedAt = time.Now()
	com.UpdatedAt = time.Now()

	err := s.Repository.AddComment(com)
	if err != nil {
		return err
	}
	return nil
}

func (s *CommentsService) DeleteComment(com_id int) error {
	err := s.Repository.DeleteComment(com_id)
	if err != nil {
		return err
	}

	return nil
}

func (s *CommentsService) GetCommentsByTask(task_id int) ([]*model.Comments, error) {
	comments, err := s.Repository.GetCommentsByTask(task_id)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentsService) GetCommentsByUser(user_id int) ([]*model.Comments, error) {
	comments, err := s.Repository.GetCommentsByUser(user_id)
	if err != nil {
		return nil, err
	}
	return comments, nil
}

func (s *CommentsService) UpdateCommentText(com_id int, new_text string) error {
	err := s.Repository.UpdateCommentText(com_id, new_text)
	if err != nil {
		return errors.New("Invalid update")
	}
	return nil
}
