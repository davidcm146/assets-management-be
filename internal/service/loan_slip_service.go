package service

import (
	"context"
	"mime/multipart"

	"github.com/davidcm146/assets-management-be.git/internal/dto"
	"github.com/davidcm146/assets-management-be.git/internal/model"
	"github.com/davidcm146/assets-management-be.git/internal/policy"
	"github.com/davidcm146/assets-management-be.git/internal/repository"
	"golang.org/x/sync/errgroup"
)

type LoanSlipService struct {
	loanSlipRepo *repository.LoanSlipRepository
	uploader     Uploader
	policy       policy.LoanSlipPolicy
}

func NewLoanSlipService(loanSlipRepo *repository.LoanSlipRepository, uploader Uploader) *LoanSlipService {
	return &LoanSlipService{loanSlipRepo: loanSlipRepo, uploader: uploader}
}

func mapLoanSlipToResponse(m *model.LoanSlip) *dto.LoanSlipResponse {
	return &dto.LoanSlipResponse{
		ID:           m.ID,
		Name:         m.Name,
		BorrowerName: m.BorrowerName,
		Department:   m.Department,
		Position:     m.Position,
		Description:  m.Description,
		Status:       m.Status.String(),
		SerialNumber: m.SerialNumber,
		Images:       m.Images,
		BorrowedDate: m.BorrowedDate,
		ReturnedDate: m.ReturnedDate,
		CreatedAt:    m.CreatedAt,
	}
}

func (s *LoanSlipService) LoanSlipsListService(ctx context.Context, query *dto.LoanSlipQuery) (*dto.PagedResult[*model.LoanSlip], error) {
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.Limit <= 0 {
		query.Limit = 10
	}
	if query.Sort == "" {
		query.Sort = "created_at"
	}
	if query.Order != "ASC" {
		query.Order = "DESC"
	}
	items, err := s.loanSlipRepo.List(ctx, query)

	if err != nil {
		return nil, err
	}
	if items == nil {
		items = []*model.LoanSlip{}
	}
	responses := make([]*dto.LoanSlipResponse, 0, len(items))
	for _, item := range items {
		responses = append(responses, mapLoanSlipToResponse(item))
	}
	total, err := s.loanSlipRepo.Count(ctx, query)
	if err != nil {
		return nil, err
	}
	return dto.NewPagedResult(items, total), nil
}

func (s *LoanSlipService) uploadImages(ctx context.Context, files []*multipart.FileHeader) ([]string, error) {

	g, ctx := errgroup.WithContext(ctx)

	imageURLs := make([]string, len(files))

	for i, file := range files {
		g.Go(func() error {
			url, err := s.uploader.Upload(ctx, file)
			if err != nil {
				return err
			}

			imageURLs[i] = url
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return nil, err
	}

	return imageURLs, nil
}

func (s *LoanSlipService) CreateLoanSlipService(ctx context.Context, userID int, req *dto.CreateLoanSlipRequest) (*model.LoanSlip, error) {
	urls, err := s.uploadImages(ctx, req.Images)
	if err != nil {
		return nil, err
	}

	loan := &model.LoanSlip{
		Name:         req.Name,
		BorrowerName: req.BorrowerName,
		Department:   req.Department,
		Position:     req.Position,
		Description:  req.Description,
		SerialNumber: req.SerialNumber,
		Images:       urls,
		CreatedBy:    userID,
		BorrowedDate: req.BorrowedDate,
		ReturnedDate: req.ReturnedDate,
		Status:       model.Borrowing,
	}

	if err := s.loanSlipRepo.Create(ctx, loan); err != nil {
		return nil, err
	}

	return loan, nil
}

func applyLoanSlipUpdate(loanSlip *model.LoanSlip, updateDTO *dto.UpdateLoanSlipRequest) {
	if updateDTO.Name != nil {
		loanSlip.Name = *updateDTO.Name
	}
	if updateDTO.BorrowerName != nil {
		loanSlip.BorrowerName = *updateDTO.BorrowerName
	}
	if updateDTO.Department != nil {
		loanSlip.Department = *updateDTO.Department
	}
	if updateDTO.Position != nil {
		loanSlip.Position = *updateDTO.Position
	}
	if updateDTO.Description != nil {
		loanSlip.Description = *updateDTO.Description
	}
	if updateDTO.Status != nil {
		loanSlip.Status = model.Status(*updateDTO.Status)
	}
	if updateDTO.SerialNumber != nil {
		loanSlip.SerialNumber = *updateDTO.SerialNumber
	}
	if updateDTO.BorrowedDate != nil {
		loanSlip.BorrowedDate = updateDTO.BorrowedDate
	}
	if updateDTO.ReturnedDate != nil {
		loanSlip.ReturnedDate = updateDTO.ReturnedDate
	}
}

func (s *LoanSlipService) UpdateLoanSlipService(
	ctx context.Context,
	id int,
	updateDTO *dto.UpdateLoanSlipRequest,
) (*model.LoanSlip, error) {

	loanSlip, err := s.loanSlipRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// upload new images
	if len(updateDTO.Images) > 0 {
		urls, err := s.uploadImages(ctx, updateDTO.Images)
		if err != nil {
			return nil, err
		}
		loanSlip.Images = append(loanSlip.Images, urls...)
	}

	// apply update fields
	applyLoanSlipUpdate(loanSlip, updateDTO)

	// save
	if err := s.loanSlipRepo.Update(ctx, loanSlip); err != nil {
		return nil, err
	}

	return loanSlip, nil
}

func (s *LoanSlipService) LoanSlipDetailService(ctx context.Context, id int) (*dto.LoanSlipResponse, error) {
	loanSlip, err := s.loanSlipRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return mapLoanSlipToResponse(loanSlip), nil
}
