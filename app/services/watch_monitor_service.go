package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"ticket-backend/app/dto"
	"ticket-backend/app/facades"
	"ticket-backend/app/models"

	"github.com/goravel/framework/support/carbon"
)

type WatchMonitorService struct {
	trainService *TrainService
}

func NewWatchMonitorService() *WatchMonitorService {
	return &WatchMonitorService{
		trainService: NewTrainService(),
	}
}

func (s *WatchMonitorService) CheckActiveWatches(ctx context.Context) error {
	today := carbon.Now().ToDateString()

	var watches []models.Watch

	if err := facades.Orm().
		WithContext(ctx).
		Query().
		Select("id").
		Where("is_active", true).
		Where("trip_date >= ?", today).
		Find(&watches); err != nil {
		return fmt.Errorf("load active watches: %w", err)
	}

	facades.Log().WithContext(ctx).Infof(
		"active watches loaded: count=%d",
		len(watches),
	)

	return s.checkWatchesConcurrently(ctx, watches)
}

func (s *WatchMonitorService) checkWatchesConcurrently(
	ctx context.Context,
	watches []models.Watch,
) error {
	const workers = 4

	jobs := make(chan uint)

	var wg sync.WaitGroup
	var mu sync.Mutex
	var checkErrors []error

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for watchID := range jobs {
				if ctx.Err() != nil {
					return
				}

				if err := s.CheckWatch(ctx, watchID); err != nil {
					facades.Log().WithContext(ctx).Errorf(
						"watch check watch id %d failed: %v",
						watchID,
						err,
					)

					mu.Lock()
					checkErrors = append(checkErrors, err)
					mu.Unlock()
				}
			}
		}()
	}

sendLoop:
	for _, watch := range watches {
		select {
		case <-ctx.Done():
			break sendLoop
		case jobs <- watch.ID:
		}
	}

	close(jobs)
	wg.Wait()

	if err := ctx.Err(); err != nil {
		checkErrors = append(checkErrors, err)
	}

	return errors.Join(checkErrors...)
}

func (s *WatchMonitorService) CheckWatch(ctx context.Context, watchId uint) error {
	today := carbon.Now().ToDateString()

	var watch models.Watch

	if err := facades.Orm().
		WithContext(ctx).
		Query().
		Where("id", watchId).
		Where("is_active", true).
		Where("trip_date >= ?", today).
		First(&watch); err != nil {
		return fmt.Errorf("watch %d: load: %w", watchId, err)
	}

	if watch.ID == 0 {
		return nil
	}

	if watch.TripDate == nil {
		return fmt.Errorf("watch %d: %w: trip_date is nil", watchId, ErrInvalidTripDate)
	}

	facades.Log().WithContext(ctx).Infof(
		"watch check started: watch_id=%d",
		watch.ID,
	)

	trains, err := s.trainService.Search(ctx, dto.SearchTrainsDTO{
		From:     watch.From,
		To:       watch.To,
		TripDate: watch.TripDate.ToDateString(),
		Lang:     "uz",
	})
	if err != nil {
		return fmt.Errorf("watch %d: search trains: %w", watchId, err)
	}

	result, err := facades.Orm().
		WithContext(ctx).
		Query().
		Model(&models.Watch{}).
		Where("id", watchId).
		Where("is_active", true).
		Update(
			"last_checked_at",
			carbon.NewDateTime(carbon.Now()),
		)
	if err != nil {
		return fmt.Errorf("watch %d: update checked time: %w", watchId, err)
	}

	// So'rov davomida watch pauza qilingan yoki o'chirilgan
	if result.RowsAffected == 0 {
		return nil
	}

	facades.Log().WithContext(ctx).Infof(
		"watch check finished: watch_id=%d trains=%d",
		watch.ID,
		len(trains),
	)

	return nil
}
