package entity

type Response struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status"`
	Message string `json:"message"`
	QueryServiceResponse
	Count int `json:"count"`
}

type QueryServiceResponse struct {
	Data []QueryServiceData `json:"data"`
}

type QueryServiceData struct {
	URL            string  `json:"url"`
	Views          int     `json:"views"`
	RelevanceScore float64 `json:"relevanceScore"`
}
