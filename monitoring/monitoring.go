package monitoring

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/team7mysupermon/devaten_middlewareapp/storage"
	"github.com/tidwall/gjson"
)

var (
	dbinstancemetrics    = make(map[string]*prometheus.GaugeVec)
	stopmetrics          = make(map[string]*prometheus.GaugeVec)
	runmetrics           = make(map[string]*prometheus.GaugeVec)
	mostexecutedmetrics  = make(map[string]*prometheus.GaugeVec)
	worstexecutedmetrics = make(map[string]*prometheus.GaugeVec)
	tableanalysismetrics = make(map[string]*prometheus.GaugeVec)
	usecaseidentifiers   string
	appclassname         string
	appipaddress         string
	appmethodname        string
	databasetype         string
	databasename         string
	starttimestamp       string
	idnum                string
	usecaseId            string
	usecasestopmetrics   = make(map[string]interface{})
	stopdetails          []storage.Stop
	reportdata           []storage.ReportData
)

func Monitor() {
	go http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":9091", nil)
	if err != nil {
		log.Fatalln("Failed to serve metrics on port 9091 ")
	}
	log.Fatal(http.ListenAndServe(":9091", nil))
}

func ParseBody(body []byte, action string) {
	justString := GetPrometheusRegisteredMetrics()
	if action == "start" {
		idnumv := gjson.Get(string(body), "data.idNum").Int()
		idnum = strconv.FormatInt(int64(idnumv), 10)
		usecase := gjson.Get(string(body), "data.usecaseIdentifier").String()
		usecaseId = usecase
		starttimestampdata := gjson.Get(string(body), "data.starttimestamp").String()
		starttimestamp = starttimestampdata
		fmt.Println(idnum)
		startdataresponse := gjson.Get(string(body), "data.dataSourceList.#.databaseType").Array()
		databasetype = strings.ToUpper(startdataresponse[0].String())
		startdatabaseName := gjson.Get(string(body), "data.dataSourceList.#.databaseName").Array()
		databasename = startdatabaseName[0].String()
		instanceinfo := gjson.Get(string(body), "data.dataSourceList.#.data").Array()

		go func() {
			for key, val := range instanceinfo[0].Map() {
				//fmt.Println(key, val)
				registered := strings.Contains(justString, "DBINSTANCE_"+strings.ToUpper(key))
				if !registered {
					dbinstancemetrics["DBINSTANCE_"+strings.ToUpper(key)] = prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "DBINSTANCE_" + strings.ToUpper(key),
							Help: "",
						}, []string{
							"database",
							"databaseName",
							"idNum",
							"usecase",
							"starttimestamp",
						},
					)

					prometheus.MustRegister(
						dbinstancemetrics["DBINSTANCE_"+strings.ToUpper(key)],
					)
				}
				dbinstancemetrics["DBINSTANCE_"+strings.ToUpper(key)].With(prometheus.Labels{"database": strings.ToUpper(databasetype), "databaseName": databasename, "idNum": idnum, "usecase": usecaseId, "starttimestamp": starttimestamp}).Set(val.Float())
			}
		}()
	}

	if action == "run" {
		runinfo := gjson.Get(string(body), "data.runSituationResult.#.data").Array()
		//go func() {
		//starttimestampdata := runinfo[0].Map()["starttimestamp"].String()
		//starttimestamp = starttimestampdata
		for _, run := range runinfo {
			for key, val := range run.Map() {
				fmt.Print(" ", val)
				registered := strings.Contains(justString, "RUN_"+strings.ToUpper(key)+"_"+databasetype)
				if !registered {
					runmetrics["RUN_"+strings.ToUpper(key)+"_"+databasetype] = prometheus.NewGaugeVec(
						prometheus.GaugeOpts{
							Name: "RUN_" + strings.ToUpper(key) + "_" + databasetype,
							Help: "",
						}, []string{
							"database",
							"starttimestamp",
							"databaseName",
							"idNum",
							"usecase",
						},
					)
					prometheus.MustRegister(
						runmetrics["RUN_"+strings.ToUpper(key)+"_"+databasetype],
					)
				}
			}
			fmt.Println()
		}
		for _, run := range runinfo {
			for key, val := range run.Map() {
				runmetrics["RUN_"+strings.ToUpper(key)+"_"+databasetype].With(prometheus.Labels{"database": strings.ToUpper(databasetype), "starttimestamp": starttimestamp, "databaseName": databasename, "idNum": idnum, "usecase": usecaseId}).Set(val.Float())

			}
		}
		//}()
	}

}
func CreateStopMetrics(arr []string) {
	justString := GetPrometheusRegisteredMetrics()
	go func() {
		for x := 0; x < len(arr); x++ {

			registered2 := strings.Contains(justString, "MOSTEXECUTE_"+strings.ToUpper(arr[x])+"_"+databasetype)

			if !registered2 {
				mostexecutedmetrics["MOSTEXECUTE_"+strings.ToUpper(arr[x])+"_"+databasetype] = prometheus.NewGaugeVec(
					prometheus.GaugeOpts{
						Name: "MOSTEXECUTE_" + strings.ToUpper(arr[x]) + "_" + databasetype,
						Help: "",
					}, []string{
						"database",
						"usecase",
						"queryid",
						"appIpAddress",
						"appClassname",
						"appMethodname",
						"databaseName",
						"idNum",
						"starttimestamp",
					},
				)
				prometheus.MustRegister(
					mostexecutedmetrics["MOSTEXECUTE_"+strings.ToUpper(arr[x])+"_"+databasetype],
				)
				mostexecutedmetrics["MOSTEXECUTE_"+strings.ToUpper(arr[x])+"_"+databasetype].WithLabelValues("database", "usecase", "queryid", "appIpAddress", "appClassname", "appMethodname", "databaseName", "idNum", "starttimestamp").Set(0)

			}

			registered := strings.Contains(justString, "WORSTEXECUTE_"+strings.ToUpper(arr[x])+"_"+databasetype)

			if !registered {
				worstexecutedmetrics["WORSTEXECUTE_"+strings.ToUpper(arr[x])+"_"+databasetype] = prometheus.NewGaugeVec(
					prometheus.GaugeOpts{
						Name: "WORSTEXECUTE_" + strings.ToUpper(arr[x]) + "_" + databasetype,
						Help: "",
					}, []string{
						"database",
						"usecase",
						"queryid",
						"appIpAddress",
						"appClassname",
						"appMethodname",
						"databaseName",
						"idNum",
						"starttimestamp",
					},
				)
				prometheus.MustRegister(
					worstexecutedmetrics["WORSTEXECUTE_"+strings.ToUpper(arr[x])+"_"+databasetype],
				)
				worstexecutedmetrics["WORSTEXECUTE_"+strings.ToUpper(arr[x])+"_"+databasetype].WithLabelValues("database", "usecase", "queryid", "appIpAddress", "appClassname", "appMethodname", "databaseName", "idNum", "starttimestamp").Set(0)

			}
			registered1 := strings.Contains(justString, "STOP_"+strings.ToUpper(arr[x])+"_"+databasetype)

			if !registered1 {
				stopmetrics["STOP_"+strings.ToUpper(arr[x])+"_"+databasetype] = prometheus.NewGaugeVec(
					prometheus.GaugeOpts{
						Name: "STOP_" + strings.ToUpper(arr[x]) + "_" + databasetype,
						Help: "",
					}, []string{
						"database",
						"usecase",
						"databaseName",
						"idNum",
						"starttimestamp",
					},
				)
				prometheus.MustRegister(
					stopmetrics["STOP_"+strings.ToUpper(arr[x])+"_"+databasetype],
				)
				stopmetrics["STOP_"+strings.ToUpper(arr[x])+"_"+databasetype].WithLabelValues("database", "usecase", "databaseName", "idNum", "starttimestamp").Set(0)
			}

		}
	}()

}
func RecordStopMetrics(body []byte) {
	//stopmetrics := GetStopMetricsMap()
	value10 := gjson.Get(string(body), "data").Array()
	for _, v := range value10 {
		for key, val := range v.Map() {
			err := json.Unmarshal([]byte(val.Raw), &stopdetails)
			if err != nil {
				panic(err)
			}
			var stopcolumnsmetrics = make(map[string]float64)
			for x := 0; x < len(stopdetails); x++ {
				stopData := stopdetails[x].ValueObjectList

				for y := 0; y < len(stopData); y++ {
					stopcolumnsmetrics["STOP_"+strings.ToUpper(stopData[y].ColumnName)+"_"+databasetype] = stopData[y].NewValue
				}
			}
			usecasestopmetrics[key] = stopcolumnsmetrics
		}
	}

	for key, element := range usecasestopmetrics {
		myMap := element.(map[string]float64)
		for columnname, value := range myMap {
			stopmetrics[columnname].WithLabelValues(strings.ToUpper(databasetype), key, databasename, idnum, starttimestamp).Set(value)
		}
	}

}

// func GetStopMetricsMap() map[string]*prometheus.GaugeVec {

//		stopmetrics["STOP_SQL_PER_SEC_GAUGE"] = prometheus.NewGaugeVec(
//			prometheus.GaugeOpts{
//				Name: "STOP_SQL_PER_SEC_GAUGE",
//				Help: "",
//			}, []string{
//				"databse",
//				"usecase",
//				"starttimestamp",
//			},
//		)
//		return stopmetrics
//	}
func GetPrometheusRegisteredMetrics() string {
	scientists := []string{
		"Einstein",
	}
	//mfs, err := promethe
	mfs, err := prometheus.DefaultGatherer.Gather()
	if err != nil {
		panic(err)
	}

	for _, mf := range mfs {
		scientists = append(scientists, mf.GetName())
	}
	justString := strings.Join(scientists, " ")
	return justString
}
func RecordReport(body []byte) {

	report := gjson.Get(string(body), "list")
	err := json.Unmarshal([]byte(report.Raw), &reportdata)
	if err != nil {
		panic(err)
	}
	for x := 0; x < len(reportdata); x++ {
		mostExecuteddata := reportdata[x].MostExecuted
		for y := 0; y < len(mostExecuteddata); y++ {
			queryid := mostExecuteddata[y].QueryId
			//fmt.Println(queryid)
			if mostExecuteddata[y].AppIpAddress == "" {
				usecaseidentifiers = mostExecuteddata[y].UsecaseIdentifier
				appclassname = ""
				appipaddress = ""
				appmethodname = ""
			} else {
				usecaseidentifiers = ""
				appclassname = mostExecuteddata[y].AppClassname
				appipaddress = mostExecuteddata[y].AppIpAddress
				appmethodname = mostExecuteddata[y].AppMethodname
			}
			res := strings.Split(mostExecuteddata[y].Colvalues, ",")
			for j := 0; j < len(res); j++ {
				medata := strings.Split(res[j], "|")
				mcolname := "MOSTEXECUTE_" + strings.ToUpper(medata[0]) + "_" + databasetype
				mcolval := medata[1]
				if s, err := strconv.ParseFloat(mcolval, 64); err == nil {
					mostexecutedmetrics[mcolname].WithLabelValues(strings.ToUpper(databasetype), usecaseidentifiers, queryid, appipaddress, appclassname, appmethodname, databasename, idnum, starttimestamp).Set(s)
				}
			}
		}
		wrostExecuteddata := reportdata[x].WrostExecuted
		for i := 0; i < len(wrostExecuteddata); i++ {
			queryid := wrostExecuteddata[i].QueryId
			if wrostExecuteddata[i].AppIpAddress == "" {
				usecaseidentifiers = wrostExecuteddata[i].UsecaseIdentifier
				appclassname = ""
				appipaddress = ""
				appmethodname = ""
			} else {
				usecaseidentifiers = ""
				appclassname = wrostExecuteddata[i].AppClassname
				appipaddress = wrostExecuteddata[i].AppIpAddress
				appmethodname = wrostExecuteddata[i].AppMethodname
			}
			res1 := strings.Split(wrostExecuteddata[i].Colvalues, ",")
			for k := 0; k < len(res1); k++ {
				wedata := strings.Split(res1[k], "|")
				wcolname := "WORSTEXECUTE_" + strings.ToUpper(wedata[0]) + "_" + databasetype
				wcolval := wedata[1]
				if s, err := strconv.ParseFloat(wcolval, 64); err == nil {
					worstexecutedmetrics[wcolname].WithLabelValues(strings.ToUpper(databasetype), usecaseidentifiers, queryid, appipaddress, appclassname, appmethodname, databasename, idnum, starttimestamp).Set(s)
				}
			}
		}
	}
}
func TableanalysisReportReg(body []byte) {
	tablereport := gjson.Get(string(body), "data").Array()
	for key := range tablereport[0].Map() {
		tablecolumn := key
		colname := "TABLEANALYSISDATA_" + strings.ToUpper(tablecolumn) + "_" + databasetype
		if key != "TABLE_NAME" {
			//fmt.Println(val)
			justString := GetPrometheusRegisteredMetrics()
			registered := strings.Contains(justString, colname)
			if !registered {
				tableanalysismetrics[colname] = prometheus.NewGaugeVec(
					prometheus.GaugeOpts{
						Name: colname,
						Help: "",
					}, []string{
						"database",
						"tablename",
						"usecase",
						"starttimestamp",
					},
				)
				prometheus.MustRegister(
					tableanalysismetrics[colname],
				)
			}
		}
	}

}
func TableanalysisReport(body []byte) {
	tablereport := gjson.Get(string(body), "data").Array()
	for _, v := range tablereport {
		reportval := v
		for key, val := range reportval.Map() {
			tablecolumn := key
			colname := "TABLEANALYSISDATA_" + strings.ToUpper(tablecolumn) + "_" + databasetype
			if key != "TABLE_NAME" {
				tableval := val.Float()
				tableanalysismetrics[colname].WithLabelValues(strings.ToUpper(databasetype), (reportval.Map()["TABLE_NAME"].String()), usecaseId, starttimestamp).Set(tableval)
			}
		}
	}
}
