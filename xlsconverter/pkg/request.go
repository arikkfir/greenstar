package pkg

const ChannelName = "xlsconverter"

type ConvertXLSFileToXLSXRequest struct {
	FileName string `json:"fileName"`
	Data     []byte `json:"data"`
	ReplyTo  string `json:"replyTo"`
}

type ConvertXLSFileToXLSXResponse struct {
	FileName string `json:"fileName"`
	Data     []byte `json:"data"`
}
