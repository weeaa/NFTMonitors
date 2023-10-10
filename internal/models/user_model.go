package models

type (
	TraderAddBody struct {
		BaseAddress     string `json:"base_address"`
		Status          string `json:"status"`
		TwitterUsername string `json:"twitter_username"`
		TwitterName     string `json:"twitter_name"`
		TwitterURL      string `json:"twitter_url"`
		UserID          int    `json:"user_id"`
		AddedBy         string `json:"added_by"`
	}

	TraderRemoveBody struct {
		BaseAddress string `json:"baseAddress"`
	}
)
