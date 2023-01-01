package model

type Account struct {
	ID          string         `json:"id"`
	DisplayName string         `json:"displayName"`
	Labels      []*KeyAndValue `json:"labels"`
	Annotations []*KeyAndValue `json:"annotations"`
	Children    []*Account     `json:"children"`
	Parent      *Account       `json:"parent"`
}
