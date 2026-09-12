package postgres

import (
	"context"
	"database/sql"

	"capsule_service/internal/domain"
	"capsule_service/internal/repository"
)

var _ repository.CapsuleRepository = (*PostgresRepository)(nil)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, capsule *domain.Cupsule) error {
	query := `
		INSERT INTO cupsule (id, owner_id, title, message, status, visibility, open_at, created_at, update_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.ExecContext(
		ctx, query,
		capsule.ID, capsule.OwnerID, capsule.Title, capsule.Message,
		string(capsule.Status), string(capsule.Visibility),
		capsule.OpenAt, capsule.CreatedAt, capsule.UpdatedAt,
	)
	return err
}

func (r *PostgresRepository) GetByID(ctx context.Context, id string) (*domain.Cupsule, error) {
	query := `
	SELECT id,owner_id, title, message, status, visibility, open_at, created_at, update_at
	FROM cupsule
	WHERE id = $1
	`
	var cupsule domain.Cupsule
	err := r.db.QueryRowContext(ctx, query, id).Scan(&cupsule.ID, &cupsule.OwnerID, &cupsule.Title, &cupsule.Message, &cupsule.Status, &cupsule.Visibility, &cupsule.OpenAt, &cupsule.CreatedAt, &cupsule.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return &cupsule, nil
}

func (r *PostgresRepository) List(ctx context.Context, ownerID string) ([]*domain.Cupsule, error) {
	query := `
	SELECT id, title, message, status, visibility, open_at, created_at, update_at
	FROM cupsule
	WHERE owner_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var cupsules []*domain.Cupsule
	for rows.Next() {
		var cupsule domain.Cupsule
		err := rows.Scan(&cupsule.ID, &cupsule.Title, &cupsule.Message, &cupsule.Status, &cupsule.Visibility, &cupsule.OpenAt, &cupsule.CreatedAt, &cupsule.UpdatedAt)
		if err != nil {
			return nil, err
		}
		cupsules = append(cupsules, &cupsule)
	}
	return cupsules, nil
}

func (r *PostgresRepository) Update(ctx context.Context, capsule *domain.Cupsule) error {
	query := `
	UPDATE cupsule
	SET 
		title = $1,
		message = $2,
		status = $3,
		update_at = NOW()
	WHERE id = $4
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		capsule.Title,   // $1
		capsule.Message, // $2
		capsule.Status,  // $3
		capsule.ID,      // $4 (теперь UUID ровно на своем месте!)
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *PostgresRepository) Delete(ctx context.Context, id string) error {
	query := `
	DELETE FROM cupsule
	WHERE id = $1
	`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}
	return nil
}
func (r *PostgresRepository) Seal(ctx context.Context, id string) error {
	query := `
	UPDATE cupsule
	SET 
		status = $1
	WHERE id = $2
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		domain.StatusSealed,
		id,
	)
	if err != nil {
		return err
	}
	return nil
}
func (r *PostgresRepository) Open(ctx context.Context, id string) error {
	query := `
	UPDATE cupsule
	SET 
		status = $1
	WHERE id = $2
	
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		domain.StatusOpened,
		id)
	if err != nil {
		return err
	}
	return nil
}
