package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type SearchTrainRequest struct {
	Name string `form:"name" json:"name"`
}

func (r *SearchTrainRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *SearchTrainRequest) Filters(ctx http.Context) map[string]any {
	return map[string]any{}
}

func (r *SearchTrainRequest) Rules(ctx http.Context) map[string]any {
	return map[string]any{}
}

func (r *SearchTrainRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *SearchTrainRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *SearchTrainRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
