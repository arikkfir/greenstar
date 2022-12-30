package pkg

type ProcessXLSXRequest struct {
	FileName string `json:"fileName"`
	Data     []byte `json:"data"`
}
