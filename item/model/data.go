package model

// Data | State 에 저장된 data
type Data struct {
	Data string `json:"data"`
}

func (i *Data) bodyModel() {
	return
}
