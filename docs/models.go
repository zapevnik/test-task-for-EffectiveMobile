package docs

type SubscriptionRequestDTO struct {
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int     `json:"price" example:"499"`
	UserID      string  `json:"user_id" example:"e5c7c66b-4a3e-4728-84d9-b6c6b46ef1a6"`
	StartDate   string  `json:"start_date" example:"07-2024"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2024"`
}

type SubscriptionResponseDTO struct {
	ID          string  `json:"id" example:"696c530f-b6c5-467f-ab70-45916e72daa7"`
	ServiceName string  `json:"service_name" example:"Netflix"`
	Price       int     `json:"price" example:"499"`
	UserID      string  `json:"user_id" example:"e5c7c66b-4a3e-4728-84d9-b6c6b46ef1a6"`
	StartDate   string  `json:"start_date" example:"07-2024"`
	EndDate     *string `json:"end_date,omitempty" example:"12-2024"`
}