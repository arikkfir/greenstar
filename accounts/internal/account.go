package internal

type Account struct {
	DisplayName string            `json:"displayName" yaml:"displayName" binding:"required"`
	ParentID    string            `json:"parentID" yaml:"parentID"`
	Labels      []string          `json:"labels" yaml:"labels"`
	Annotations map[string]string `json:"annotations" yaml:"annotations"`
}
type AccountWithID struct {
	ID string `json:"id" yaml:"id"`
	Account
}
