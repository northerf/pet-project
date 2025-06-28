package repository

import (
	"database/sql"
	"pet-project/pkg/model"
)

type PostgresCommentsRepository struct {
	DB *sql.DB
}

type CommentsRepository interface {
	AddComment(com *model.Comments) error
	DeleteComment(com_id int) error
	GetCommentsByTask(task_id int) ([]*model.Comments, error)
	GetCommentsByUser(user_id int) ([]*model.Comments, error)
	UpdateCommentText(com_id int, new_text string) error
}

func (r *PostgresCommentsRepository) AddComment(com *model.Comments) error {
	query := `INSERT INTO comments (task_id, user_id, text, created_at)
				VALUES ($1, $2, $3, $4, $5) RETURNING id`
	return r.DB.QueryRow(query, com.TaskID, com.UserID, com.Text, com.CreatedAt).Scan(&com.ID)
}

func (r *PostgresCommentsRepository) DeleteComment(com_id int) error {
	query := `DELETE FROM comments WHERE id = $1`
	_, err := r.DB.Exec(query, com_id)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresCommentsRepository) GetCommentsByTask(task_id int) ([]*model.Comments, error) {
	query := `SELECT task_id, text, created_at, user_id, id FROM comments WHERE task_id = $1`
	rows, err := r.DB.Query(query, task_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	comments := []*model.Comments{}

	for rows.Next() {
		var comment model.Comments
		err := rows.Scan(&comment.TaskID, &comment.Text, &comment.CreatedAt, &comment.UserID, &comment.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *PostgresCommentsRepository) GetCommentsByUser(user_id int) ([]*model.Comments, error) {
	query := `SELECT id, text, created_at, task_id, user_id FROM comments WHERE user_id = $1`
	rows, err := r.DB.Query(query, user_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []*model.Comments{}

	for rows.Next() {
		var comment model.Comments
		if err := rows.Scan(&comment.ID, &comment.Text, &comment.CreatedAt, &comment.TaskID, &comment.UserID); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return comments, nil
}

func (r *PostgresCommentsRepository) UpdateCommentText(com_id int, new_text string) error {
	comment := &model.Comments{}
	query := `UPDATE comments SET text = $1 WHERE id = $2`
	_, err := r.DB.Exec(query, comment.Text, comment.ID)
	if err != nil {
		return err
	}
	return nil
}
