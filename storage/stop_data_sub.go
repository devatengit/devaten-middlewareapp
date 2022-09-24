package storage

type Stop struct {
	UsecaseIdentifier           string                `json:"usecaseIdentifier"`
	Starttimestamp              string                `json:"starttimestamp"`
	Endtimestamp                string                `json:"endtimestamp"`
	Id_num                      int                   `json:"id_num"`
	Id                          int                   `json:"id"`
	RowNumber                   int                   `json:"rowNumber"`
	ValueObjectList             []ValueObjectList     `json:"valueObjectList"`
	ResultInPercentage          float64               `json:"resultInPercentage"`
	Result                      string                `json:"Result"`
	UsecaseResult               interface{}           `json:"usecaseResult"`
	LatestUsecaseStartTimestamp string                `json:"latestUsecaseStartTimestamp"`
	OldUsecaseStartTimestamp    string                `json:"oldUsecaseStartTimestamp"`
	BatchId                     int                   `json:"batchId"`
	DataSourceId                int                   `json:"dataSourceId"`
	ApplicationId               int                   `json:"applicationId"`
	StatusId                    int                   `json:"statusId"`
	ConfigId                    int                   `json:"configId"`
	AlertConfigInfoList         []AlertConfigInfoList `json:"alertConfigInfoList"`
	AlertCriteriaModel          AlertCriteriaModel    `json:"alertCriteriaModel"`
	ResponseTime                int                   `json:"responseTime"`
	MostExecuted                []MostExecuted        `json:"mostExecuted"`
	WrostExecuted               []WrostExecuted       `json:"wrostExecuted"`

	MostExecutedByUsecaseId  []MostExecutedByUsecaseId  `json:"mostExecutedByUsecaseId"`
	WrostExecutedByUsecaseId []WrostExecutedByUsecaseId `json:"wrostExecutedByUsecaseId"`
	Comparedtimestamp        interface{}                `json:"comparedtimestamp"`
}

type AlertCriteriaModel struct {
	StatusId         int         `json:"statusId"`
	ApplicationId    int         `json:"applicationId"`
	Failure          int         `json:"failure"`
	Warning          int         `json:"warning"`
	Success          int         `json:"success"`
	Improvements     int         `json:"improvements"`
	UpdatedTimestamp interface{} `json:"updatedTimestamp"`
	UpdatedBy        interface{} `json:"updatedBy"`
}

type AlertConfigInfoList struct {
	Id               int         `json:"id"`
	ApplicationId    int         `json:"applicationId"`
	ConfigId         int         `json:"configId"`
	ColumnName       string      `json:"columnName"`
	ColumnTitle      string      `json:"columnTitle"`
	UpdatedTimestamp interface{} `json:"updatedTimestamp"`
}

type ValueObjectList struct {
	ValueObjectId   float64 `json:"valueObjectId"`
	Value           float64 `json:"value"`
	NewValue        float64 `json:"newValue"`
	OldValue        float64 `json:"oldValue"`
	ComparedNumber  float64 `json:"comparedNumber"`
	Result          float64 `json:"result"`
	ResultMessage   string  `json:"resultMessage"`
	ColumnName      string  `json:"columnName"`
	ColumnTitle     string  `json:"columnTitle"`
	AlertHistoryId  int     `json:"alertHistoryId"`
	DataSourceId    int     `json:"dataSourceId"`
	New             bool    `json:"new"`
	AddedForCompare bool    `json:"addedForCompare"`
}
