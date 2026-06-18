package jobs

type Metrics struct {
	Pending    int `json:"pending"`
	Processing int `json:"processing"`
	Completed  int `json:"completed"`
	Failed     int `json:"failed"`
	Retrying   int `json:"retrying"`
	Total      int `json:"total"`
}
