package config

import "ticket-backend/app/facades"

func init() {
	config := facades.Config()
	config.Add("telegram", map[string]any{
		"bot_token":     config.Env("TELEGRAM_BOT_TOKEN", ""),
		"init_data_ttl": config.Env("TELEGRAM_INIT_DATA_TTL", "1h"),
	})
}
