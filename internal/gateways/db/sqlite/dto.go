package sqlite

type HistoryDTO struct {
	ID          uint64 `json:"id,omitempty" db:"id"`
	Date        string `json:"date,omitempty" db:"date"`
	Type        string `json:"type,omitempty" db:"type"`
	Status      string `json:"status,omitempty" db:"status"`
	Error       string `json:"error,omitempty" db:"error"`
	Description string `json:"description,omitempty" db:"description"`
	DataID      uint64 `json:"data_id,omitempty" db:"data_id"`
}

type DataDTO struct {
	ID          uint64 `json:"id,omitempty" db:"id"`
	DataType    string `json:"data_type,omitempty" db:"data_type"`
	DataContent string `json:"data_content,omitempty" db:"data_content"`
}

type HistoryDataDTO struct {
	HistoryDTO
	DataDTO
}

func NewDataDTO() *DataDTO {
	return &DataDTO{}
}

func NewHistoryDTO() *HistoryDTO {
	return &HistoryDTO{}
}
