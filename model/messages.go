package model

//RequestList the pagination and filter of a job listing
type RequestList struct {
	Limit  int    `json:"limit"`
	Page   int    `json:"page"`
	Status string `json:"status"`
}

//ResponseList the list of jobs matching the request criteria
type ResponseList struct {
	Jobs []Job `json:"jobs"`
}

//RequestCreate the data about a job to be created or replaced
type RequestCreate struct {
	Job `json:"job"`
}

//ResponseCreate result status of the create request
type ResponseCreate struct {
	Error  error `json:"error"`
	Status int64 `json:"status"`
}

//RequestActivate the information about a job to change to active
type RequestActivate struct {
	PartnerID int64 `json:"partner_id"`
}

//ResponseActivate the activation status response
type ResponseActivate struct {
	Error  error `json:"error"`
	Status int64 `json:"status"`
}

//RequestPercentage the category to totalize
type RequestPercentage struct {
	CategoryID int64 `json:"category_id"`
}

//ResponsePercentage totalized category count response
type ResponsePercentage struct {
	CategoryPercentage
}
