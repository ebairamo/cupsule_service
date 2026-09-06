package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

	"capsule_service/internal/domain"
	"capsule_service/internal/repository/postgres"
	service "capsule_service/internal/usecase"
)

func main() {
	// 1. Подключаемся к PostgreSQL
	connStr := "postgres://postgres:12345@localhost:5432/cupsule_service?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Ошибка подключения к БД: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("База недоступна: %v", err)
	}
	fmt.Println("Успешно подключились к PostgreSQL!")

	// 2. Инициализируем репозиторий и сервис
	repo := postgres.NewPostgresRepository(db)
	capsuleSvc := service.NewCapsuleService(repo)
	ctx := context.Background()

	ownerID := "11111111-1111-1111-1111-111111111111"

	// --- STEP 1: CREATE ---
	fmt.Println("\n--- 1. Создание капсулы ---")
	createdCapsule, err := capsuleSvc.CreateCapsule(
		ctx,
		ownerID,
		"Моя первая капсула",
		"Привет из прошлого!",
		domain.VisibilityPrivate,
		time.Now().Add(24*time.Hour),
	)
	if err != nil {
		log.Fatalf("Ошибка создания: %v", err)
	}
	fmt.Printf("Создана! ID: %s, Status: %s\n", createdCapsule.ID, createdCapsule.Status)

	// --- STEP 2: GET BY ID ---
	fmt.Println("\n--- 2. Получение по ID ---")
	fetched, err := capsuleSvc.GetByID(ctx, createdCapsule.ID)
	if err != nil {
		log.Fatalf("Ошибка получения: %v", err)
	}
	fmt.Printf("Найдена! Title: %s, Message: %s\n", fetched.Title, fetched.Message)

	// --- STEP 3: LIST ---
	fmt.Println("\n--- 3. Список капсул владельца ---")
	list, err := repo.List(ctx, ownerID)
	if err != nil {
		log.Fatalf("Ошибка получения списка: %v", err)
	}
	fmt.Printf("Всего капсул у пользователя %s: %d шт.\n", ownerID, len(list))

	// --- STEP 4: UPDATE ---
	fmt.Println("\n--- 4. Обновление капсулы ---")
	fetched.Title = "Обновленный заголовок"
	fetched.UpdatedAt = time.Now()
	if err := repo.Update(ctx, fetched); err != nil {
		log.Fatalf("Ошибка обновления: %v", err)
	}

	updated, _ := capsuleSvc.GetByID(ctx, fetched.ID)
	fmt.Printf("Обновлено! Новый Title: %s\n", updated.Title)

	// --- STEP 5: DELETE ---
	fmt.Println("\n--- 5. Удаление капсулы ---")
	if err := repo.Delete(ctx, fetched.ID); err != nil {
		log.Fatalf("Ошибка удаления: %v", err)
	}
	fmt.Println("Капсула успешно удалена!")

	// Проверяем, что действительно удалилась
	_, err = capsuleSvc.GetByID(ctx, fetched.ID)
	if err != nil {
		fmt.Println("Проверка успешна: капсула больше не найдена в БД.")
	}
}
