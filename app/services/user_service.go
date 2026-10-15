package services

import (
	"context"
	"errors"
	"ticket-backend/app/dto"
	"ticket-backend/app/facades"
	"ticket-backend/app/models"
	"unicode/utf8"
)

var ErrInvalidIdentity = errors.New("invalid identity")

type UserService struct{}

func NewUserService() *UserService {
	return &UserService{}
}

func (s *UserService) Profile(
	ctx context.Context,
	identity dto.TelegramIdentity,
) (*models.User, error) {
	if identity.TelegramID <= 0 ||
		utf8.RuneCountInString(identity.Username) > 100 ||
		utf8.RuneCountInString(identity.DisplayName) > 200 {
		return nil, ErrInvalidIdentity
	}

	attributes := map[string]any{
		"telegram_id": identity.TelegramID,
	}

	values := map[string]any{
		"username":     nil,
		"display_name": identity.DisplayName,
	}

	if identity.Username != "" {
		values["username"] = identity.Username
	}

	orm := facades.Orm().WithContext(ctx)
	var user models.User

	if err := orm.Query().FirstOrCreate(
		&user,
		attributes,
		values,
	); err != nil {
		// Parallel so‘rov foydalanuvchini oldinroq yaratgan bo‘lishi mumkin.
		var postgresError interface {
			SQLState() string
		}

		if !errors.As(err, &postgresError) ||
			postgresError.SQLState() != "23505" {
			return nil, err
		}
	}

	if _, err := orm.Query().
		Model(&models.User{}).
		Where("telegram_id", identity.TelegramID).
		Update(values); err != nil {
		return nil, err
	}

	user = models.User{}

	if err := orm.Query().
		Where("telegram_id", identity.TelegramID).
		FirstOrFail(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
