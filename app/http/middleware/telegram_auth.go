package middleware

import (
	"encoding/json"
	"net/url"
	"strconv"
	"strings"
	"ticket-backend/app/models"
	"time"
	"unicode/utf8"

	"ticket-backend/app/dto"
	"ticket-backend/app/facades"

	"github.com/goravel/framework/contracts/http"
	initdata "github.com/telegram-mini-apps/init-data-golang"
)

type TelegramAuth struct{}

type telegramIdentityKey struct{}

func (m *TelegramAuth) Signature() string {
	return "TelegramAuth"
}

func (m *TelegramAuth) Handle(ctx http.Context) {
	config := facades.Config()

	if config.GetString("app.env") == "local" &&
		config.EnvBool("TELEGRAM_AUTH_MOCK", false) &&
		ctx.Request().Header("Authorization") == "" {
		m.authenticateLocalUser(ctx)
		return
	}

	botToken := strings.TrimSpace(config.GetString("telegram.bot_token"))
	maxAge, err := time.ParseDuration(config.GetString("telegram.init_data_ttl", "1h"))

	if err != nil || botToken == "" || maxAge <= 0 {
		_ = ctx.Response().Json(http.StatusServiceUnavailable, http.Json{
			"error": http.Json{
				"success": false,
				"message": "Telegram auth not configured",
			},
		}).Abort()
		return
	}

	header := ctx.Request().Header("Authorization")
	scheme, raw, found := strings.Cut(header, " ")

	if !found ||
		!strings.EqualFold(scheme, "tma") ||
		raw == "" ||
		len(raw) > 8192 ||
		len(ctx.Request().Headers().Values("Authorization")) != 1 {
		telegramUnauthorized(ctx)
		return
	}

	fields, err := url.ParseQuery(raw)
	if err != nil {
		telegramUnauthorized(ctx)
		return
	}

	for _, values := range fields {
		if len(values) != 1 {
			telegramUnauthorized(ctx)
			return
		}
	}

	if err := initdata.Validate(raw, botToken, maxAge); err != nil {
		facades.Log().Errorf("Telegram initData validation failed: %v", err)

		telegramUnauthorized(ctx)
		return
	}

	authDate, err := strconv.ParseInt(fields.Get("auth_date"), 10, 64)
	if err != nil ||
		authDate <= 0 ||
		time.Unix(authDate, 0).After(time.Now().Add(30*time.Second)) {
		telegramUnauthorized(ctx)
		return
	}

	var telegramUser struct {
		ID        int64  `json:"id"`
		Username  string `json:"username"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := json.Unmarshal([]byte(fields.Get("user")), &telegramUser); err != nil || telegramUser.ID <= 0 {
		telegramUnauthorized(ctx)
		return
	}

	displayName := strings.TrimSpace(telegramUser.FirstName + " " + telegramUser.LastName)
	if utf8.RuneCountInString(telegramUser.Username) > 100 ||
		utf8.RuneCountInString(displayName) > 200 {
		telegramUnauthorized(ctx)
		return
	}

	ctx.WithValue(telegramIdentityKey{}, dto.TelegramIdentity{
		TelegramID:  telegramUser.ID,
		Username:    telegramUser.Username,
		DisplayName: displayName,
	})

	ctx.Request().Next()
}

func TelegramIdentityFromContext(ctx http.Context) (dto.TelegramIdentity, bool) {
	identity, ok := ctx.Value(telegramIdentityKey{}).(dto.TelegramIdentity)
	return identity, ok
}

func telegramUnauthorized(ctx http.Context) {
	_ = ctx.Response().Json(http.StatusUnauthorized, http.Json{
		"error": http.Json{
			"success": false,
			"message": "Valid initData is required",
		},
	}).Abort()
}

func (m *TelegramAuth) authenticateLocalUser(ctx http.Context) {
	var user models.User

	err := facades.Orm().
		WithContext(ctx.Context()).
		Query().
		OrderBy("id").
		First(&user)

	if err != nil {
		facades.Log().WithContext(ctx.Context()).
			Errorf("Load local auth user: %v", err)

		_ = ctx.Response().Json(
			http.StatusInternalServerError,
			http.Json{
				"error": http.Json{
					"code":    "internal_error",
					"message": "Could not load local user.",
				},
			},
		).Abort()
		return
	}

	if user.ID == 0 || user.TelegramID <= 0 {
		_ = ctx.Response().Json(
			http.StatusServiceUnavailable,
			http.Json{
				"error": http.Json{
					"code": "auth_unavailable",
					"message": "First local user is missing " +
						"or has no valid Telegram ID.",
				},
			},
		).Abort()
		return
	}

	username := ""
	if user.Username != nil {
		username = *user.Username
	}

	ctx.WithValue(telegramIdentityKey{}, dto.TelegramIdentity{
		TelegramID:  user.TelegramID,
		Username:    username,
		DisplayName: user.DisplayName,
	})

	ctx.Request().Next()
}
