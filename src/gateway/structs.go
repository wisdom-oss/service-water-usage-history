package gateway

type TargetInformation struct {
	Id        string  `json:"id"`
	CreatedAt float64 `json:"created_at"`
	Upstream  struct {
		Id string `json:"id"`
	} `json:"upstream"`
	Address string   `json:"target"`
	Weight  int      `json:"weight"`
	Tags    []string `json:"tags"`
}
