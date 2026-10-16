package services

import (
	"context"
	"errors"
	"ticket-backend/app/dto"
	"ticket-backend/app/facades"
	"ticket-backend/app/models"

	"github.com/goravel/framework/support/carbon"
)

var ErrWatchNotFound = errors.New("watch not found")
var ErrInvalidTripDate = errors.New("invalid trip date")

type WatchService struct{}

func NewWatchService() *WatchService {
	return &WatchService{}
}

func (s *WatchService) Create(ctx context.Context, userID uint, input dto.CreateWatchDTO) (*models.Watch, error) {
	tripDate := carbon.ParseByFormat(
		input.TripDate,
		"Y-m-d",
	)

	if !tripDate.IsValid() {
		return nil, ErrInvalidTripDate
	}

	watch := models.Watch{
		UserID:   userID,
		From:     input.From,
		To:       input.To,
		TripDate: carbon.NewDate(tripDate),
		Filters:  input.Filters,
		IsActive: true,
	}

	if err := facades.Orm().
		WithContext(ctx).
		Query().
		Create(&watch); err != nil {
		return nil, err
	}

	return &watch, nil
}

func (s *WatchService) List(ctx context.Context, userID uint) ([]models.Watch, error) {
	watches := make([]models.Watch, 0)

	if err := facades.Orm().
		WithContext(ctx).
		Query().
		Where("user_id", userID).
		OrderByDesc("id").
		Find(&watches); err != nil {
		return nil, err
	}

	return watches, nil
}

func (s *WatchService) Delete(ctx context.Context, userID uint, watchID uint) error {
	result, err := facades.Orm().WithContext(ctx).Query().
		Where("user_id", userID).
		Where("id", watchID).
		Delete(&models.Watch{})
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ErrWatchNotFound
	}

	return nil
}

func (s *WatchService) Pause(ctx context.Context, userID uint, watchID uint) error {
	result, err := facades.Orm().
		WithContext(ctx).
		Query().
		Model(&models.Watch{}).
		Where("user_id", userID).
		Where("id", watchID).
		Update("is_active", false)
	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ErrWatchNotFound
	}

	return nil
}

func (s *WatchService) Resume(ctx context.Context, userID, watchID uint) error {
	result, err := facades.Orm().WithContext(ctx).Query().
		Model(&models.Watch{}).
		Where("user_id", userID).
		Where("id", watchID).
		Update("is_active", true)

	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return ErrWatchNotFound
	}

	return nil
}
