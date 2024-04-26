package web


type WebResponse struct {
	Code int `json:"code"`
	Status string `json:"status"`
	Message string `json:"message"`
	Data any `json:"data"`
}

func (w WebResponse)Error()string  {
	return w.Message
}