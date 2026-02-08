package repository

import (
	"context"

	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoanSlipRepository struct {
	db *pgxpool.Pool
}

func NewLoanSlipRepository(db *pgxpool.Pool) *LoanSlipRepository {
	return &LoanSlipRepository{db: db}
}

func (r *LoanSlipRepository) Create(ctx context.Context, loanSlip *model.LoanSlip) error {
	status := model.Borrowing
	_, err := r.db.Exec(ctx,
		"INSERT INTO loan_slips (name, borrower_name, department, position, description, status, serial_number, images, created_by, borrowed_date, returned_date, updated_at, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, NOW(), NOW())",
		loanSlip.Name, loanSlip.BorrowerName, loanSlip.Department, loanSlip.Position, loanSlip.Description, status, loanSlip.SerialNumber, loanSlip.Images, loanSlip.CreatedBy, loanSlip.BorrowedDate, loanSlip.ReturnedDate,
	)
	return err
}

func (r *LoanSlipRepository) Update(ctx context.Context, loanSlip *model.LoanSlip) error {
	_, err := r.db.Exec(ctx,
		"UPDATE loan_slips SET name=$1, borrower_name=$2, department=$3, position=$4, description=$5, status=$6, serial_number=$7, images=$8, borrowed_date=$9, returned_date=$10, updated_at=NOW() WHERE id=$11",
		loanSlip.Name, loanSlip.BorrowerName, loanSlip.Department, loanSlip.Position, loanSlip.Description, loanSlip.Status, loanSlip.SerialNumber, loanSlip.Images, loanSlip.BorrowedDate, loanSlip.ReturnedDate, loanSlip.ID,
	)
	return err
}

func (r *LoanSlipRepository) FindByID(ctx context.Context, id int) (*model.LoanSlip, error) {
	row := r.db.QueryRow(ctx,
		`SELECT id, name, borrower_name, department, position, description, status, serial_number, images, created_by, borrowed_date, returned_date, updated_at, created_at
		 FROM loan_slips WHERE id = $1`, id)

	var loanSlip model.LoanSlip
	err := row.Scan(
		&loanSlip.ID,
		&loanSlip.Name,
		&loanSlip.BorrowerName,
		&loanSlip.Department,
		&loanSlip.Position,
		&loanSlip.Description,
		&loanSlip.Status,
		&loanSlip.SerialNumber,
		&loanSlip.Images,
		&loanSlip.CreatedBy,
		&loanSlip.BorrowedDate,
		&loanSlip.ReturnedDate,
		&loanSlip.UpdatedAt,
		&loanSlip.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &loanSlip, nil
}
