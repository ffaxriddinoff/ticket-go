package requests

import (
	"fmt"
	"strings"
	"ticket-backend/app/dto"

	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
	"github.com/goravel/framework/support/carbon"
)

type CreateWatchRequest struct {
	From     string              `form:"from" json:"from"`
	To       string              `form:"to" json:"to"`
	TripDate string              `form:"trip_date" json:"trip_date"`
	Filters  dto.WatchFiltersDTO `form:"filters" json:"filters"`
}

func (r *CreateWatchRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *CreateWatchRequest) Rules(ctx http.Context) map[string]any {
	today := carbon.Now().ToDateString()

	return map[string]any{
		"from": "required|string|max:30",
		"to":   "required|string|max:30|different:from",

		"trip_date": "required|string|date_format:2006-01-02|after_or_equal:" + today,

		"filters":          "required|map",
		"filters.position": "required|string|in:all,upper,lower",
		"filters.zone":     "required|string|in:all,no_side",

		"filters.trains":   "list",
		"filters.trains.*": "required|string|max:100",

		"filters.classes":   "list",
		"filters.classes.*": "string|max:100",

		"filters.wagon_types":   "sometimes|list",
		"filters.wagon_types.*": "sometimes|string|in:seated,platzkart,coupe,lux,soft,sv,general",
	}
}
func (r *CreateWatchRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateWatchRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *CreateWatchRequest) PrepareForValidation(
	ctx http.Context,
	data validation.Data,
) error {
	for _, field := range []string{"from", "to", "trip_date"} {
		value, _ := data.Get(field)
		if text, ok := value.(string); ok {
			if err := data.Set(field, strings.TrimSpace(text)); err != nil {
				return err
			}
		}
	}

	filterValue, exists := data.Get("filters")
	if !exists || filterValue == nil {
		filterValue = map[string]any{}
	}

	if filters, ok := filterValue.(map[string]any); ok {
		for _, field := range []string{"position", "zone"} {
			value, exists := filters[field]
			if !exists || value == nil {
				filters[field] = "all"
				continue
			}

			if text, ok := value.(string); ok {
				text = strings.TrimSpace(text)
				if text == "" {
					text = "all"
				}
				filters[field] = text
			}
		}

		for _, field := range []string{"trains", "wagon_types", "classes"} {
			filters[field] = normalizeWatchList(filters[field])
		}

		if err := data.Set("filters", filters); err != nil {
			return err
		}
	}

	return nil
}

func normalizeWatchList(value any) any {
	var values []any

	switch list := value.(type) {
	case nil:
		return []any{}
	case []any:
		values = list
	case []string:
		values = make([]any, len(list))
		for i, item := range list {
			values[i] = item
		}
	default:
		return value
	}

	result := make([]any, 0, len(values))
	seen := make(map[string]struct{}, len(values))

	for _, item := range values {
		if text, ok := item.(string); ok {
			item = strings.TrimSpace(text)
		}

		key := fmt.Sprintf("%T:%v", item, item)
		if _, exists := seen[key]; exists {
			continue
		}

		seen[key] = struct{}{}
		result = append(result, item)
	}

	return result
}
