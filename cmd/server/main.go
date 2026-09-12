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
	fmt.Println("\n--- 5. Удаление первой капсулы ---")
	if err := repo.Delete(ctx, fetched.ID); err != nil {
		log.Fatalf("Ошибка удаления: %v", err)
	}
	fmt.Println("Капсула успешно удалена!")

	// --- STEP 6: SEAL (Запечатывание) ---
	fmt.Println("\n--- 6. Создание и запечатывание новой капсулы ---")
	targetCapsule, err := capsuleSvc.CreateCapsule(
		ctx,
		ownerID,
		"Капсула для проверки Seal и Open",
		"Секретный текст",
		domain.VisibilityPrivate,
		time.Now().Add(1*time.Hour), // Валидное время в будущем
	)
	if err != nil {
		log.Fatalf("Ошибка создания капсулы для теста: %v", err)
	}

	if err := capsuleSvc.Seal(ctx, targetCapsule.ID); err != nil {
		log.Fatalf("Ошибка при запечатывании: %v", err)
	}

	sealedCapsule, err := capsuleSvc.GetByID(ctx, targetCapsule.ID)
	if err != nil {
		log.Fatalf("Ошибка получения капсулы после Seal: %v", err)
	}
	fmt.Printf("Успешно! Капсула %s запечатана. Статус: %v\n", sealedCapsule.ID, sealedCapsule.Status)

	// --- STEP 7: OPEN (Открытие) ---
	fmt.Println("\n--- 7. Открытие запечатанной капсулы ---")

	// Временно сдвигаем open_at в прошлую дату через SQL, чтобы заставить time.Now().After(openAt) пройти
	_, err = db.ExecContext(ctx, "UPDATE cupsule SET open_at = $1 WHERE id = $2", time.Now().Add(-1*time.Hour), targetCapsule.ID)
	if err != nil {
		log.Fatalf("Ошибка сдвига времени в БД: %v", err)
	}

	if err := capsuleSvc.Open(ctx, targetCapsule.ID); err != nil {
		log.Fatalf("Ошибка при открытии капсулы: %v", err)
	}

	openedCapsule, err := capsuleSvc.GetByID(ctx, targetCapsule.ID)
	if err != nil {
		log.Fatalf("Ошибка получения капсулы после Open: %v", err)
	}
	fmt.Printf("Успешно! Капсула %s открыта. Итоговый статус: %v\n", openedCapsule.ID, openedCapsule.Status)
}
