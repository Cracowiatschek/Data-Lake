package gios

type LinksDTO struct {
	Next  string `json:"next"`
	Prev  string `json:"prev"`
	First string `json:"first"`
	Last  string `json:"last"`
}
