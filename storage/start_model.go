package storage

type StartAutoGenerated struct {
	Status        string        `json:"status"`
	ResponseCode  int           `json:"responseCode"`
	StartMetaData StartMetaData `json:"data"`
	ErrorMessage  interface{}   `json:"errorMessage"`
	ErrorCode     interface{}   `json:"errorCode"`
	ReportLink    interface{}   `json:"reportLink"`
}

//	type StartData struct {
//		Statements              float64 `json:"STATEMENTS"`
//		StatementLatencyInS     float64 `json:"STATEMENT_LATENCY_IN_S"`
//		FileIoLatencyInS        float64 `json:"FILE_IO_LATENCY_IN_S"`
//		CurrentConnections      float64 `json:"CURRENT_CONNECTIONS"`
//		DatabaseSizeInMb        float64 `json:"DATABASE_SIZE_IN_MB"`
//		StatementAvgLatencyInMs float64 `json:"STATEMENT_AVG_LATENCY_IN_MS"`
//		ApplicationID           float64 `json:"APPLICATION_ID"`
//		FileIos                 float64 `json:"FILE_IOS"`
//		TableScans              float64 `json:"TABLE_SCANS"`
//		DataSourceID            float64 `json:"DATA_SOURCE_ID"`
//		UsecaseIdentifier       float64 `json:"USECASE_IDENTIFIER"`
//		UniqueUsers             float64 `json:"UNIQUE_USERS"`
//	}
type StartDataSourceList struct {
	DataSourceID    int         `json:"dataSourceId"`
	DatabaseType    string      `json:"databaseType"`
	DatabaseName    string      `json:"databaseName"`
	SchemaName      string      `json:"schemaName"`
	HostURL         string      `json:"hostUrl"`
	StartData       interface{} `json:"data"`
	ValueObjectList interface{} `json:"valueObjectList"`
}
type StartMetaData struct {
	IDNum                 int                   `json:"idNum"`
	UsecaseIdentifier     string                `json:"usecaseIdentifier"`
	ApplicationID         int                   `json:"applicationId"`
	ApplicationName       string                `json:"applicationName"`
	ApplicationIdentifier string                `json:"applicationIdentifier"`
	Starttimestamp        string                `json:"starttimestamp"`
	StartDataSourceList   []StartDataSourceList `json:"dataSourceList"`
}
