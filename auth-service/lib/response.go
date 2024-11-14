package lib

type Paginate struct {
	Page      int `json:"page"`
	PerPage   int `json:"per_page"`
	Total     int `json:"total"`
	TotalPage int `json:"total_page"`
}

type ResponseParams struct {
	StatusCode int
	Message    string
	Paginate   *Paginate
	Data       any
}

type FilterParams struct {
	Page   int
	Limit  int
	Offset int
	Search string
}

type ResponseData struct {
	Code     int       `json:"code"`
	Status   bool      `json:"status"`
	Message  string    `json:"message"`
	Paginate *Paginate `json:"paginate,omitempty"`
	Data     any       `json:"data"`
}

type ResponseNoPaginate struct {
	Code    int    `json:"code"`
	Status  bool   `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func Response(params ResponseParams) any {

	var response any
	var status bool

	if params.StatusCode >= 200 && params.StatusCode <= 299 {
		status = true
	} else {
		status = false
	}

	if params.Data != nil {
		response = &ResponseData{
			Code:     params.StatusCode,
			Status:   status,
			Message:  params.Message,
			Paginate: params.Paginate,
			Data:     params.Data,
		}
	} else if params.Paginate == nil {
		response = &ResponseNoPaginate{
			Code:    params.StatusCode,
			Status:  status,
			Message: params.Message,
			Data:    params.Data,
		}
	}

	return response
}
