package main

import (
	"database/sql"
	"log"
	"net/http"
	"pet-project/internal/handler"
	"pet-project/internal/middleware"
	"pet-project/internal/realtime"
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
	comRepo := &repository.PostgresCommentsRepository{DB: db}
	notRepo := &repository.PostgresNotificationRepository{DB: db}

	clientManager := realtime.NewClientManager()

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
	comService := &service.CommentsService{
		Repository: comRepo,
	}
	notService := &service.NotificationService{
		Repository:    notRepo,
		ClientManager: clientManager,
	}

	authHandler := &handler.AuthHandler{AuthService: authService}
	projectHandler := &handler.ProjectHandler{ProjectService: projectService}
	taskHandler := &handler.TaskHandler{TaskService: taskService}
	commentsHandler := &handler.CommentsHandler{CommentsService: comService}
	notificationHandler := &handler.NotificationHandler{NotificationService: notService}
	notificationWSHandler := &handler.NotificationWSHandler{
		ClientManager: clientManager,
		JwtSecret:     []byte("supersecretkey"),
		AuthService:   authService,
	}

	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			if req.Method == "OPTIONS" {
				w.WriteHeader(http.StatusOK)
				return
			}
			next.ServeHTTP(w, req)
		})
	})

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

	r.Route("/comments", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware([]byte("supersecretkey")))
		r.Post("/", commentsHandler.AddCommentRequest)
		r.Delete("/{comID}", commentsHandler.DeleteCommentRequest)
		r.Get("/task/{taskID}", commentsHandler.GetCommentsByTaskRequest)
		r.Get("/user/{userID}", commentsHandler.GetCommentsByUserRequest)
		r.Put("/{comID}", commentsHandler.UpdateCommentTextRequest)
	})

	r.Route("/notification", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware([]byte("supersecretkey")))
		r.Post("/", notificationHandler.CreateNotification)
		r.Get("/", notificationHandler.GetNotifications)
		r.Post("/mark-read", notificationHandler.MarkAsRead)
		r.Get("/unread-count", notificationHandler.CountUnread)
	})

	r.Get("/ws/notifications", notificationWSHandler.WSNotifications)

	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
