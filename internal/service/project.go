package service

import (
	"errors"
	"pet-project/internal/repository"
	"pet-project/pkg/model"
	"time"
)

type ProjectService struct {
	Repository repository.ProjectRepository
}

func (s *ProjectService) CreateProject(name, description string, ownerID int) (*model.Project, error) {
	if err := s.validateOwner(ownerID); err != nil {
		return nil, err
	}

	if err := s.validatePName(name); err != nil {
		return nil, err
	}

	project := &model.Project{
		Name:        name,
		Description: description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	project, err := s.saveProjectToDB(project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) validateOwner(ownerID int) error {
	if ownerID == 0 {
		return errors.New("Project doesn't have owner")
	}
	return nil
}

func (s *ProjectService) saveProjectToDB(project *model.Project) (*model.Project, error) {
	err := s.Repository.CreateProject(project)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (s *ProjectService) validatePName(name string) error {
	if name == "" {
		return errors.New("Project need a title!")
	}
	return nil
}

func (s *ProjectService) GetByIDProject(projectId int, ownerID int) (*model.Project, error) {
	if err := s.validateOwner(ownerID); err != nil {
		return nil, err
	}

	project, err := s.Repository.GetByIDProject(projectId)
	if err != nil {
		return nil, err
	}

	if project.OwnerID != ownerID {
		return nil, errors.New("You don't have permission to access this project")
	}

	return project, nil
}

func (s *ProjectService) UpdateProject(project *model.Project) error {
	project.UpdatedAt = time.Now()
	err := s.Repository.UpdateProject(project)
	if err != nil {
		return err
	}
	return nil
}

func (s *ProjectService) DeleteProject(projectID, userID int) error {
	if err := s.validateOwner(userID); err != nil {
		return err
	}

	project, err := s.Repository.GetByIDProject(projectID)
	if err != nil {
		return err
	}

	if project.OwnerID != userID {
		return errors.New("You don't have permission to manage project")
	}

	err = s.Repository.DeleteProject(projectID)
	if err != nil {
		return err
	}
	return nil
}
