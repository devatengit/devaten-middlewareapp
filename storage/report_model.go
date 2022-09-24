package storage

type ReportData struct {
	AlertConfigInfoList         []AlertConfigInfoList     `json:"alertConfigInfoList"`
	AlertCriteriaModel          []AlertCriteriaModel      `json:"responseCode"`
	ApplicationId               int                       `json:"data"`
	BatchId                     int                       `json:"errorMessage"`
	Comparedtimestamp           string                    `json:"comparedtimestamp"`
	ConfigId                    int                       `json:"configId"`
	DataSourceId                int                       `json:"dataSourceId"`
	Endtimestamp                string                    `json:"endtimestamp"`
	Id                          int                       `json:"id"`
	Id_num                      int                       `json:"id_num"`
	LatestUsecaseStartTimestamp string                    `json:"latestUsecaseStartTimestamp"`
	MostExecuted                []MostExecuted            `json:"mostExecuted"`
	MostExecutedByUsecaseId     []MostExecutedByUsecaseId `json:"mostExecutedByUsecaseId"`
	OldUsecaseStartTimestamp    string                    `json:"oldUsecaseStartTimestamp"`
	ResponseTime                int                       `json:"responseTime"`
	Result                      string                    `json:"result"`
	ResultInPercentage          float64                   `json:"resultInPercentage"`
	RowNumber                   int                       `json:"rowNumber"`
	Starttimestamp              string                    `json:"starttimestamp"`
	StatusId                    int                       `json:"statusId"`
	UsecaseIdentifier           string                    `json:"usecaseIdentifier"`
	UsecaseResult               interface{}               `json:"usecaseResult"`
	ValueObjectList             []ValueObjectList         `json:"valueObjectList"`
	WrostExecuted               []WrostExecuted           `json:"wrostExecuted"`

	WrostExecutedByUsecaseId []WrostExecutedByUsecaseId `json:"wrostExecutedByUsecaseId"`
}
type WrostExecutedByUsecaseId struct {
	AppClassname      string        `json:"appClassname"`
	AppCurrentTime    string        `json:"appCurrentTime"`
	AppIpAddress      string        `json:"appIpAddress"`
	AppMethodname     string        `json:"appMethodname"`
	AppPackagename    string        `json:"appPackagename"`
	Baseline          interface{}   `json:"baseline"`
	ColumnsList       []ColumnsList `json:"columnsList"`
	Colvalues         string        `json:"colvalues"`
	DbName            interface{}   `json:"dbName"`
	DbTypeId          int           `json:"dbTypeId"`
	DigestId          interface{}   `json:"digestId"`
	Endtimestamp      interface{}   `json:"endtimestamp"`
	ExecCount         int           `json:"execCount"`
	FullScan          interface{}   `json:"fullScan"`
	Id                int           `json:"id"`
	IdNum             int           `json:"idNum"`
	IndexMissing      bool          `json:"indexMissing"`
	QueryId           string        `json:"queryId"`
	Schema_NAME       interface{}   `json:"schema_NAME"`
	SqlStatement      string        `json:"sqlStatement"`
	Starttimestamp    string        `json:"starttimestamp"`
	TableanalysisJson interface{}   `json:"tableanalysisJson"`
	Timespend         int           `json:"timespend"`
	UsecaseIdentifier string        `json:"usecaseIdentifier"`
}
type WrostExecuted struct {
	AppClassname      string      `json:"appClassname"`
	AppCurrentTime    string      `json:"appCurrentTime"`
	AppIpAddress      string      `json:"appIpAddress"`
	AppMethodname     string      `json:"appMethodname"`
	AppPackagename    string      `json:"appPackagename"`
	Baseline          interface{} `json:"baseline"`
	ColumnsList       interface{} `json:"columnsList"`
	Colvalues         string      `json:"colvalues"`
	DbName            interface{} `json:"dbName"`
	DbTypeId          int         `json:"dbTypeId"`
	DigestId          interface{} `json:"digestId"`
	Endtimestamp      interface{} `json:"endtimestamp"`
	ExecCount         int         `json:"execCount"`
	FullScan          interface{} `json:"fullScan"`
	Id                int         `json:"id"`
	IdNum             int         `json:"idNum"`
	IndexMissing      bool        `json:"indexMissing"`
	QueryId           string      `json:"queryId"`
	Schema_NAME       interface{} `json:"schema_NAME"`
	SqlStatement      string      `json:"sqlStatement"`
	Starttimestamp    string      `json:"starttimestamp"`
	TableanalysisJson interface{} `json:"tableanalysisJson"`
	Timespend         int         `json:"timespend"`
	UsecaseIdentifier string      `json:"usecaseIdentifier"`
}
type MostExecutedByUsecaseId struct {
	AppClassname      string        `json:"appClassname"`
	AppCurrentTime    string        `json:"appCurrentTime"`
	AppIpAddress      string        `json:"appIpAddress"`
	AppMethodname     string        `json:"appMethodname"`
	AppPackagename    string        `json:"appPackagename"`
	Baseline          interface{}   `json:"baseline"`
	ColumnsList       []ColumnsList `json:"columnsList"`
	Colvalues         string        `json:"colvalues"`
	DbName            interface{}   `json:"dbName"`
	DbTypeId          int           `json:"dbTypeId"`
	DigestId          interface{}   `json:"digestId"`
	Endtimestamp      interface{}   `json:"endtimestamp"`
	ExecCount         int           `json:"execCount"`
	FullScan          interface{}   `json:"fullScan"`
	Id                int           `json:"id"`
	IdNum             int           `json:"idNum"`
	IndexMissing      bool          `json:"indexMissing"`
	QueryId           string        `json:"queryId"`
	Schema_NAME       interface{}   `json:"schema_NAME"`
	SqlStatement      string        `json:"sqlStatement"`
	Starttimestamp    string        `json:"starttimestamp"`
	TableanalysisJson interface{}   `json:"tableanalysisJson"`
	Timespend         int           `json:"timespend"`
	UsecaseIdentifier string        `json:"usecaseIdentifier"`
}
type MostExecuted struct {
	AppClassname      string      `json:"appClassname"`
	AppCurrentTime    string      `json:"appCurrentTime"`
	AppIpAddress      string      `json:"appIpAddress"`
	AppMethodname     string      `json:"appMethodname"`
	AppPackagename    string      `json:"appPackagename"`
	Baseline          interface{} `json:"baseline"`
	ColumnsList       interface{} `json:"columnsList"`
	Colvalues         string      `json:"colvalues"`
	DbName            interface{} `json:"dbName"`
	DbTypeId          int         `json:"dbTypeId"`
	DigestId          interface{} `json:"digestId"`
	Endtimestamp      interface{} `json:"endtimestamp"`
	ExecCount         int         `json:"execCount"`
	FullScan          interface{} `json:"fullScan"`
	Id                int         `json:"id"`
	IdNum             int         `json:"idNum"`
	IndexMissing      bool        `json:"indexMissing"`
	QueryId           string      `json:"queryId"`
	Schema_NAME       interface{} `json:"schema_NAME"`
	SqlStatement      string      `json:"sqlStatement"`
	Starttimestamp    string      `json:"starttimestamp"`
	TableanalysisJson interface{} `json:"tableanalysisJson"`
	Timespend         int         `json:"timespend"`
	UsecaseIdentifier string      `json:"usecaseIdentifier"`
}
type ColumnsList struct {
	AddedForCompare bool        `json:"addedForCompare"`
	AlertHistoryId  int         `json:"alertHistoryId"`
	ColumnName      string      `json:"columnName"`
	ColumnTitle     interface{} `json:"columnTitle"`
	ComparedNumber  float64     `json:"comparedNumber"`
	DataSourceId    int         `json:"dataSourceId"`
	New             bool        `json:"new"`
	NewValue        float64     `json:"newValue"`
	OldValue        float64     `json:"oldValue"`
	Result          float64     `json:"result"`
	ResultMessage   interface{} `json:"resultMessage"`
	Value           float64     `json:"value"`
	ValueObjectId   float64     `json:"valueObjectId"`
}
