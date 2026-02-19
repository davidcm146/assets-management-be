package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DashboardRepository interface {
	GetLoanMetrics(ctx context.Context, filter model.TimeFilter) (*model.LoanMetrics, error)
}

type dashboardRepository struct {
	db *pgxpool.Pool
}

func NewDashboardRepository(db *pgxpool.Pool) DashboardRepository {
	return &dashboardRepository{db: db}
}

func (r *dashboardRepository) GetLoanMetrics(ctx context.Context, filter model.TimeFilter) (*model.LoanMetrics, error) {
	var result model.LoanMetrics

	where := []string{"1=1"}
	args := []any{}
	argPos := 1

	if filter.From != nil {
		where = append(where, fmt.Sprintf("created_at >= $%d", argPos))
		args = append(args, &filter.From)
		argPos++
	}

	if filter.To != nil {
		where = append(where, fmt.Sprintf("created_at <= $%d", argPos))
		args = append(args, &filter.To)
		argPos++
	}

	borrowingPos := argPos
	returnedPos := argPos + 1
	overduePos := argPos + 2

	args = append(args,
		model.Borrowing,
		model.Returned,
		model.Overdue,
	)

	query := fmt.Sprintf(`
		SELECT
			COUNT(*) AS total,
			COALESCE(SUM(CASE WHEN status = $%d THEN 1 ELSE 0 END),0) AS borrowing,
			COALESCE(SUM(CASE WHEN status = $%d THEN 1 ELSE 0 END),0) AS returned,
			COALESCE(SUM(CASE WHEN status = $%d THEN 1 ELSE 0 END),0) AS overdue
		FROM loan_slips
		WHERE %s
	`,
		borrowingPos,
		returnedPos,
		overduePos,
		strings.Join(where, " AND "),
	)

	row := r.db.QueryRow(ctx, query, args...)

	err := row.Scan(
		&result.Total,
		&result.Borrowing,
		&result.Returned,
		&result.Overdue,
	)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
