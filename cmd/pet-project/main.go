package main

import (
	"database/sql"
	"log"
	"net/http"
	"pet-project/internal/handler"
	"pet-project/internal/middleware"
	"pet-project/internal/repository"
	"pet-project/internal/service"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=localhost port=5432 user=petuser password=petpassword dbname=petprojectdb sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	userRepo := &repository.PostgresUserRepository{DB: db}
	projectRepo := &repository.PostgresProjectRepository{DB: db}
	taskRepo := &repository.PostgresTaskRepository{DB: db}

	authService := &service.AuthService{
		Repository: userRepo,
		JwtSecret:  []byte("supersecretkey"),
	}
	projectService := &service.ProjectService{
		Repository: projectRepo,
	}
	taskService := &service.TaskService{
		Repository: taskRepo,
	}

	authHandler := &handler.AuthHandler{AuthService: authService}
	projectHandler := &handler.ProjectHandler{ProjectService: projectService}
	taskHandler := &handler.TaskHandler{TaskService: taskService}

	r := chi.NewRouter()

	r.Post("/login", authHandler.Login)
	r.Post("/register", authHandler.Register)

	r.Route("/projects", func(pr chi.Router) {
		pr.Use(middleware.AuthMiddleware([]byte("supersecretkey")))

		pr.Post("/", projectHandler.CreateProject)              // POST /projects — создание проекта
		pr.Get("/{projectID}", projectHandler.GetProjectInfo)   // GET /projects/{id} — получение информации о проекте
		pr.Put("/{projectID}", projectHandler.UpdateProject)    // PUT /projects/{id} — обновление проекта
		pr.Delete("/{projectID}", projectHandler.DeleteProject) // DELETE /projects/{id} — удаление проекта
	})

	r.Route("/tasks", func(tr chi.Router) {
		tr.Use(middleware.AuthMiddleware([]byte("supersecretkey")))
		tr.Post("/", taskHandler.CreateTaskRequest)
		tr.Put("/{taskID}", taskHandler.UpdateProjectRequest)
		tr.Get("/{taskID}", taskHandler.GetByIDTaskRequest)
		tr.Delete("/{taskID}", taskHandler.DeleteTaskRequest)
	})

	r.With(middleware.AuthMiddleware([]byte("supersecretkey"))).
		Get("/projects/{projectID}/tasks", taskHandler.ListByProjectTaskRequest)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
