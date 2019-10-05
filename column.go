package structgenerator

type Column struct {
	Extra         string `json:"extra"`
	ColumnName    string `json:"column_name"`
	DataType      string `json:"data_type"`
	ColumnType    string `json:"column_type"`
	ColumnDefault string `json:"column_default"`
	NullAble      string `json:"null_able"`
	ColumnKey     string `json:"column_key"`
	ColumnComment string `json:"column_comment"`
}
