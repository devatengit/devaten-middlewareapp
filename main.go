package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
	"github.com/team7mysupermon/devaten_middlewareapp/storage"
	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	"github.com/team7mysupermon/devaten_middlewareapp/monitoring"

	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	docs "github.com/team7mysupermon/devaten_middlewareapp/docs"
)

var (
	// The authentication token needed to be able to get the access token when logging in
	authToken = "Basic cGVyZm9ybWFuY2VEYXNoYm9hcmRDbGllbnRJZDpsamtuc3F5OXRwNjEyMw=="

	/*
		Instantiated when a user calls the login API call.
		Contains the authentication token
	*/
	Tokenresponse     storage.Token
	appurl                  = ""
	recordingmail           = ""
	explainjson             = ""
	jira                    = ""
	report                  = ""
	loginusername           = ""
	password                = ""
	loginresponse           = 0
	scrapintervaltime int64 = 5
	/*
		Closes the goroutine that scrapes the recording.
		The goroutine is started when the user starts the recording
	*/
	//quit = make(chan bool)
	stopInterval = false
)

func main() {
	err := godotenv.Load("middleware.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	apphost := os.Getenv("APP_HOST")
	recordingmaile := os.Getenv("RECORDING_MAIL")
	recordingmail = recordingmaile
	explainjsone := os.Getenv("EXPLAIN_JSON")
	explainjson = explainjsone
	jirae := os.Getenv("JIRA")
	jira = jirae
	reporte := os.Getenv("REPORT")
	report = reporte

	userdata := os.Getenv("LOGIN_USER_NAME")
	loginusername = userdata

	pass := os.Getenv("PASSWORD")
	password = pass

	scrapintervaltimee := os.Getenv("SCRAP_INTERVAL_TIME")
	//scrapintervaltime=scrapintervaltimee
	scrapintervaltime, err = strconv.ParseInt(scrapintervaltimee, 16, 64)

	fmt.Println(scrapintervaltime, err, reflect.TypeOf(scrapintervaltime))

	appurl = apphost
	fmt.Println(appurl)
	go monitoring.Monitor()
	docs.SwaggerInfo.BasePath = ""
	router := gin.Default()

	router.GET("/Start/:Usecase/:Appiden", startRecording)
	router.GET("/Stop/:Usecase/:Appiden", stopRecording)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	err1 := router.Run(":8999")
	if err1 != nil {
		return
	}
}

func startRecording(c *gin.Context) {
	// Creates the command structure by taking information from the URL call
	var command storage.StartAndStopCommand
	if err := c.ShouldBindUri(&command); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	fmt.Println(loginusername)
	fmt.Println(password)
	go scrapeWithIntervalForLogin(loginusername, password)

	fmt.Println(command.ApplicationIdentifier)
	stopInterval = false
	time.Sleep(5 * time.Second)
	if loginresponse == 200 {

		c.JSON(loginresponse, gin.H{"Control": "Login Successfully."})
		var res = Operation(command.Usecase, "start", command.ApplicationIdentifier)

		fmt.Println(res.StatusCode)
		if res.StatusCode == 200 {
			PrepareStopMetrics(command.ApplicationIdentifier)
			c.JSON(res.StatusCode, gin.H{"Control": "A recording has now started"})
			go scrapeWithInterval(command)
			go scrapeWithIntervalforactive(command)
		} else {
			var error1 = res.Proto
			c.JSON(res.StatusCode, gin.H{"Control": error1})
		}

	} else {

		c.JSON(loginresponse, gin.H{"Control": "Login fail. Please enter correct devaten dashboard username and password or host url in middleware.env fie.."})
	}

	// Starts the scraping on a seperat thread

}

func stopRecording(c *gin.Context) {
	// Creates the command structure by taking information from the URL call
	var command storage.StartAndStopCommand
	if err := c.ShouldBindUri(&command); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Sends true through the quit channel to the goroutine that is scraping the recording
	go scrapeWithIntervalForLogin(loginusername, password)
	var res = StopRecordingdata(command.Usecase, command.ApplicationIdentifier)
	fmt.Println(res.StatusCode)
	if res.StatusCode == 200 {

		c.JSON(res.StatusCode, gin.H{"Control": "A recording has now ended"})
	} else {
		var error1 = res.Proto
		c.JSON(res.StatusCode, gin.H{"Control": error1})
	}

}

func getAuthToken(loginusername string, password string) *http.Response {
	var url = appurl + "/oauth/token"
	method := "POST"

	payload := strings.NewReader(generateUserInfo(loginusername, password))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	defer res.Body.Close()
	loginresponse = res.StatusCode
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	err = json.Unmarshal(body, &Tokenresponse)
	if err != nil {
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}

	fmt.Printf("%s : %s\n", Tokenresponse.Type, Tokenresponse.AccessToken)

	fmt.Println("******************************************** Auth Token ********************************************")

	return res
}

func Operation(usecase string, action string, applicationIdentifier string) *http.Response {
	url := appurl + "/devaten/data/operation?usecaseIdentifier=" + usecase + "&action=" + action
	method := "GET"
	// applicationIdentifier1 := applicationIdentifier
	// applicationIdentifier1 = strings.Replace(applicationIdentifier1, "\n", "", -1)

	payload := strings.NewReader("")
	//fmt.Println(url)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	req.Header.Add("applicationIdentifier", applicationIdentifier)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Tokenresponse.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	defer res.Body.Close()
	//fmt.Println(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	responsecode := gjson.Get(string(body), "responseCode").Int()
	if responsecode == 200 {
		monitoring.ParseBody(body, action)
	} else {
		res.StatusCode = 500
		res.Proto = gjson.Get(string(body), "errorMessage").String()
	}

	return res

}

func OperationWhoIsActive(applicationIdentifier string) *http.Response {
	url := appurl + "/devaten/data/getwhoIsActiveInformation"
	method := "GET"
	payload := strings.NewReader("")
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	req.Header.Add("applicationIdentifier", applicationIdentifier)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Tokenresponse.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	defer res.Body.Close()
	//fmt.Println(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	responsecode := gjson.Get(string(body), "responseCode").Int()
	if responsecode == 200 {
		monitoring.WhoIsActive(body)
	} else {
		res.StatusCode = 500
		res.Proto = gjson.Get(string(body), "errorMessage").String()
	}

	return res

}
func StopRecordingdata(usecase string, applicationIdentifier string) *http.Response {
	url := appurl + "/devaten/data/stopRecording?usecaseIdentifier=" + usecase + "&recordingMail=" + recordingmail + "&explainJson=" + explainjson + "&jira=" + jira + "&report=" + report + "&inputSource=application&frocefullyStop=false"
	method := "GET"

	payload := strings.NewReader("")
	fmt.Println(Tokenresponse.AccessToken)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	req.Header.Add("applicationIdentifier", applicationIdentifier)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Tokenresponse.AccessToken)

	res, err := client.Do(req)
	fmt.Println(res.StatusCode)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}

	defer res.Body.Close()

	//fmt.Println(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	responsecode := gjson.Get(string(body), "responseCode").Int()
	if responsecode == 200 {
		monitoring.RecordStopMetrics(body)
		report := gjson.Get(string(body), "reportLink")
		//fmt.Println(report)
		lastBin := strings.LastIndex(report.String(), "view")

		//fmt.Println(report.String()[lastBin+5 : len(report.String())])
		reporturl := report.String()[lastBin+5 : len(report.String())]

		reportdata(reporturl, applicationIdentifier)
		idNum := strings.Split(reporturl, "/")

		//fmt.Println(idNum[1])

		tableanalysisdata(idNum[1], idNum[0], applicationIdentifier)
	} else {
		res.StatusCode = 500
		res.Proto = gjson.Get(string(body), "errorMessage").String()
	}
	//quit <- true
	stopInterval = true
	return res

}

func tableanalysisdata(idNum string, usecase string, applicationIdentifier string) *http.Response {
	url := appurl + "/userMgt/getTableWiseDetailsInformation?idNum=" + idNum + "&usecaseIdentifier=" + usecase
	method := "GET"

	payload := strings.NewReader("")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	req.Header.Add("applicationIdentifier", applicationIdentifier)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Tokenresponse.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	defer res.Body.Close()
	//fmt.Println(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	//fmt.Println(string(body))
	responsecode := gjson.Get(string(body), "responseCode").Int()

	if responsecode == 200 {
		monitoring.TableanalysisReportReg(body)
		monitoring.TableanalysisReport(body)
	} else {
		res.StatusCode = 500
	}

	return res

}
func reportdata(usecase string, applicationIdentifier string) *http.Response {
	url := appurl + "/userMgt/report/" + usecase
	method := "GET"
	payload := strings.NewReader("")
	fmt.Println(usecase)
	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	req.Header.Add("applicationIdentifier", applicationIdentifier)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Tokenresponse.AccessToken)
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	defer res.Body.Close()
	//fmt.Println(res.Body)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	monitoring.RecordReport(body)
	return res

}
func PrepareStopMetrics(applicationIdentifier string) *http.Response {

	url := appurl + "/devaten/data/getAlertConfigInfoByApplicationIdentifier"
	method := "GET"

	payload := strings.NewReader("")

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	req.Header.Add("applicationIdentifier", applicationIdentifier)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Bearer "+Tokenresponse.AccessToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return &http.Response{
			Status:     err.Error(),
			StatusCode: 500,
		}
	}

	var resultdata []string
	value10 := gjson.Get(string(body), "data.#.columnName").Array()
	for _, v := range value10 {
		resultdata = append(resultdata, v.Str)
	}

	monitoring.CreateStopMetrics(resultdata)

	return res

}

func scrapeWithInterval(command storage.StartAndStopCommand) {
	for {

		if stopInterval {
			return
		} else {
			Operation(command.Usecase, "run", command.ApplicationIdentifier)
		}
		time.Sleep(time.Duration(scrapintervaltime) * time.Second)

	}
}
func scrapeWithIntervalForLogin(loginusername string, password string) {
	for {

		if stopInterval {
			return
		} else {
			getAuthToken(loginusername, password)

		}

		time.Sleep(3600 * time.Second)

	}
}
func scrapeWithIntervalforactive(command storage.StartAndStopCommand) {
	for {

		if stopInterval {
			return
		} else {
			OperationWhoIsActive(command.ApplicationIdentifier)
		}

		time.Sleep(time.Duration(scrapintervaltime) * time.Second)

	}
}

// Takes a loginusername and a password and generates the string that is needed to login
func generateUserInfo(username string, password string) string {
	var userInfo = "username=" + username + "&password=" + password + "&grant_type=password"
	fmt.Println(userInfo)
	return userInfo
}
