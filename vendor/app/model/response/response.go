package response

type ResponseData struct {
	Recommends []Post `json:"recommend,omitempty"`
	Algorithm  int      `json:"alg,omitempty"`
}

type Post struct {
	ID string `json:"id,omitempty"`
}
