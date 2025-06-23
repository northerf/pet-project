package repository

import (
	"database/sql"
	"pet-project/pkg/model"
)

type PostgresProjectRepository struct {
	DB *sql.DB
}

type ProjectRepository interface {
	CreateProject(project *model.Project) error
	UpdateProject(project *model.Project) error
	GetByIDProject(id int) (*model.Project, error)
	DeleteProject(id int) error
}

func (rp *PostgresProjectRepository) CreateProject(project *model.Project) error {
	query := `INSERT INTO projects (name, description, owner_id, created_at, updated_at)
	          VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return rp.DB.QueryRow(query, project.Name, project.Description, project.OwnerID, project.CreatedAt, project.UpdatedAt).
		Scan(&project.ID)
}

func (rp *PostgresProjectRepository) UpdateProject(project *model.Project) error {
	query := `UPDATE projects SET name = $2, description = $3, ownerId = $4, updatedAt = $5 WHERE id = $1`
	_, err := rp.DB.Exec(query, project.ID, project.Name, project.Description, project.OwnerID, project.UpdatedAt)
	if err != nil {
		return err
	}
	return nil
}

func (rp *PostgresProjectRepository) GetByIDProject(id int) (*model.Project, error) {
	project := &model.Project{}
	query := `SELECT id, name, description, ownerId, createdAt, updatedAt FROM projects WHERE id = $1`
	row := rp.DB.QueryRow(query, id)
	err := row.Scan(&project.ID, &project.Name, &project.Description, &project.OwnerID, &project.CreatedAt, &project.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return project, nil
}

func (rp *PostgresProjectRepository) DeleteProject(id int) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := rp.DB.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
