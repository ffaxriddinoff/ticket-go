package config

import "ticket-backend/app/facades"

func init() {
	config := facades.Config()

	config.Add("eticket", map[string]any{
		"base_url": config.Env("ETICKET_BASE_URL", ""),
		"timeout":  config.Env("ETICKET_TIMEOUT", 10),
	})
}
