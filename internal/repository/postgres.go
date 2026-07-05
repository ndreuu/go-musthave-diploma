package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"go-musthave-diploma/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(db *sql.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateUser(ctx context.Context, login string, passwordHash string) (*model.User, error) {
	const query = `
		INSERT INTO users (login, password_hash)
		VALUES ($1, $2)
		RETURNING id, login, password_hash
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, login, passwordHash).
		Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, ErrAlreadyExists
		}

		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) GetUserByLogin(ctx context.Context, login string) (*model.User, error) {
	const query = `
		SELECT id, login, password_hash
		FROM users
		WHERE login = $1
	`

	var user model.User
	err := r.db.QueryRowContext(ctx, query, login).
		Scan(&user.ID, &user.Login, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *PostgresRepository) CreateOrder(ctx context.Context, userID int64, number string) error {
	const query = `
		INSERT INTO orders (number, user_id, status, uploaded_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(
		ctx,
		query,
		number,
		userID,
		model.OrderStatusNew,
		time.Now(),
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return ErrAlreadyExists
		}

		return err
	}

	return nil
}

func (r *PostgresRepository) GetOrderByNumber(ctx context.Context, number string) (*model.Order, error) {
	const query = `
		SELECT number, user_id, status, accrual, uploaded_at
		FROM orders
		WHERE number = $1
	`

	var order model.Order
	err := r.db.QueryRowContext(ctx, query, number).
		Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return &order, nil
}

func (r *PostgresRepository) GetOrdersByUserID(ctx context.Context, userID int64) ([]model.Order, error) {
	const query = `
		SELECT number, user_id, status, accrual, uploaded_at
		FROM orders
		WHERE user_id = $1
		ORDER BY uploaded_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var order model.Order
		if err := rows.Scan(&order.Number, &order.UserID, &order.Status, &order.Accrual, &order.UploadedAt); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PostgresRepository) CreateWithdrawal(ctx context.Context, userID int64, order string, sum float64) error {
	const query = `
		INSERT INTO withdrawals (user_id, order_number, sum, processed_at)
		VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query, userID, order, sum, time.Now())
	return err
}

func (r *PostgresRepository) GetWithdrawalsByUserID(ctx context.Context, userID int64) ([]model.Withdrawal, error) {
	const query = `
		SELECT order_number, user_id, sum, processed_at
		FROM withdrawals
		WHERE user_id = $1
		ORDER BY processed_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var withdrawals []model.Withdrawal

	for rows.Next() {
		var withdrawal model.Withdrawal
		if err := rows.Scan(
			&withdrawal.Order,
			&withdrawal.UserID,
			&withdrawal.Sum,
			&withdrawal.ProcessedAt,
		); err != nil {
			return nil, err
		}

		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, nil
}

func (r *PostgresRepository) GetOrdersForAccrual(ctx context.Context, limit int) ([]model.Order, error) {
	const query = `
		SELECT number, user_id, status, accrual, uploaded_at
		FROM orders
		WHERE status IN ($1, $2)
		ORDER BY uploaded_at ASC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		model.OrderStatusNew,
		model.OrderStatusProcessing,
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []model.Order

	for rows.Next() {
		var order model.Order
		if err := rows.Scan(
			&order.Number,
			&order.UserID,
			&order.Status,
			&order.Accrual,
			&order.UploadedAt,
		); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PostgresRepository) UpdateOrderAccrual(ctx context.Context, number string, status model.OrderStatus, accrual *float64) error {
	const query = `
		UPDATE orders
		SET status = $2,
		    accrual = $3
		WHERE number = $1
	`

	_, err := r.db.ExecContext(ctx, query, number, status, accrual)
	return err
}

func (r *PostgresRepository) GetAccrualSumByUserID(ctx context.Context, userID int64) (float64, error) {
	const query = `
		SELECT COALESCE(SUM(accrual), 0)
		FROM orders
		WHERE user_id = $1
	`

	var sum float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}

func (r *PostgresRepository) GetWithdrawalSumByUserID(ctx context.Context, userID int64) (float64, error) {
	const query = `
		SELECT COALESCE(SUM(sum), 0)
		FROM withdrawals
		WHERE user_id = $1
	`

	var sum float64
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&sum)
	if err != nil {
		return 0, err
	}

	return sum, nil
}
