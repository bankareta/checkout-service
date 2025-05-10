package model

type AddEtalaseRequest struct {
	Username     string `json:"username,omitempty" validate:"required"`
	EtalaseName  string `json:"etalase_name,omitempty" validate:"required"`
	ReferenceId  string `json:"reference_id,omitempty"  validate:"required"`
	PhotoEtalase string `json:"photo_etalase,omitempty"`
}
