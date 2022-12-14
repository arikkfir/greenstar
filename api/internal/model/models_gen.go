// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AccountUpdate struct {
	DisplayName *string             `json:"displayName"`
	Labels      []*KeyAndValueInput `json:"labels"`
	Annotations []*KeyAndValueInput `json:"annotations"`
	ParentID    *string             `json:"parentID"`
}

type KeyAndValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type KeyAndValueInput struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NewAccount struct {
	ID          *string             `json:"id"`
	DisplayName string              `json:"displayName"`
	Labels      []*KeyAndValueInput `json:"labels"`
	Annotations []*KeyAndValueInput `json:"annotations"`
	ParentID    *string             `json:"parentID"`
}
