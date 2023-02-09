package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
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
	Tokenresponse storage.Token
	appurl        = ""
	/*
		Closes the goroutine that scrapes the recording.
		The goroutine is started when the user starts the recording
	*/
	quit = make(chan bool)
)

func main() {
	err := godotenv.Load("middleware.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	apphost := os.Getenv("APP_HOST")
	appurl = apphost
	fmt.Println(appurl)
	go monitoring.Monitor()
	docs.SwaggerInfo.BasePath = ""
	router := gin.Default()

	// getting env variables SITE_TITLE and
	// The API calls
	router.GET("/Login/:Username/:Password", getAuthToken)
	router.GET("/Start/:Usecase/:Appiden", startRecording)
	router.GET("/Stop/:Usecase/:Appiden", stopRecording)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	// Starts the program
	err1 := router.Run(":8999")
	if err1 != nil {
		return
	}
}

// @BasePath /Start/{Usecase}/{Appiden}

// PingExample godoc
// @Summary Start a recording
// @Schemes
// @Description This endpoint is to stop a recording and needs a usecase and a applicationIdentifier as parameters.
// @Tags example
// @Param Usecase path string true ":Usecase"
// @Param Appiden path string true ":Appiden"
// @Accept json
// @Produce json
// @Success 200
// @Router /Start/{Usecase}/{Appiden} [get]
func startRecording(c *gin.Context) {
	// Creates the command structure by taking information from the URL call
	var command storage.StartAndStopCommand
	if err := c.ShouldBindUri(&command); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	fmt.Println(command.ApplicationIdentifier)

	var res = Operation(command.Usecase, "start", command.ApplicationIdentifier)

	fmt.Println(res.StatusCode)
	if res.StatusCode == 200 {
		PrepareStopMetrics(command.ApplicationIdentifier)
		c.JSON(res.StatusCode, gin.H{"Control": "A recording has now started"})
		go scrapeWithInterval(command)
	} else {
		var error1 = res.Proto
		c.JSON(res.StatusCode, gin.H{"Control": error1})
	}

	// Starts the scraping on a seperat thread

}

// @BasePath /Stop/{Usecase}/{Appiden}

// PingExample godoc
// @Summary Stop a recording
// @Schemes
// @Description This endpoint is to stop a recording and needs a usecase and a applicationIdentifier as parameters.
// @Tags example
// @Param Usecase path string true ":Usecase"
// @Param Appiden path string true ":Appiden"
// @Accept json
// @Produce json
// @Success 200
// @Router /Stop/{Usecase}/{Appiden} [get]
func stopRecording(c *gin.Context) {
	// Creates the command structure by taking information from the URL call
	var command storage.StartAndStopCommand
	if err := c.ShouldBindUri(&command); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}
	// Sends true through the quit channel to the goroutine that is scraping the recording
	quit <- true

	var res = StopRecordingdata(command.Usecase, command.ApplicationIdentifier)
	fmt.Println(res.StatusCode)
	if res.StatusCode == 200 {
		c.JSON(res.StatusCode, gin.H{"Control": "A recording has now ended"})
	} else {
		var error1 = res.Proto
		c.JSON(res.StatusCode, gin.H{"Control": error1})
	}

}

// @BasePath /Login/{Username}/{Password}

// PingExample godoc
// @Summary Send middleware user information
// @Schemes
// @Description this is a request to give the middleware user information. this will allow the middleware to set up the authentication token need to start and stop the recording.
// @Tags example
// @Param Username path string true ":Username"
// @Param Password path string true ":Password"
// @Produce json
// @Success 200
// @Router /Login/{Username}/{Password} [get]
func getAuthToken(c *gin.Context) {
	var url = appurl + "/oauth/token"
	method := "POST"

	// Creates the command structure by taking information from the URL call
	var command storage.LoginCommand
	if err := c.ShouldBindUri(&command); err != nil {
		c.JSON(400, gin.H{"msg": err})
		return
	}

	// Generates the user info string
	payload := strings.NewReader(generateUserInfo(command.Username, command.Password))

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", authToken)

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = json.Unmarshal(body, &Tokenresponse)
	if err != nil {
		return
	}
	if res.StatusCode == 200 {
		c.JSON(res.StatusCode, gin.H{"Control": "Login successfully."})
	} else {
		c.JSON(res.StatusCode, gin.H{"Control": "There is some problem in Login."})
	}

	fmt.Println("******************************************** Auth Token ********************************************")
	//fmt.Printf("%s : %s\n", Tokenresponse.Type, Tokenresponse.AccessToken)
}

func Operation(usecase string, action string, applicationIdentifier string) *http.Response {
	url := appurl + "/devaten/data/operation?usecaseIdentifier=" + usecase + "&action=" + action
	method := "GET"
	// applicationIdentifier1 := applicationIdentifier
	// applicationIdentifier1 = strings.Replace(applicationIdentifier1, "\n", "", -1)

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
		monitoring.ParseBody(body, action)
	} else {
		res.StatusCode = 500
		res.Proto = gjson.Get(string(body), "errorMessage").String()
	}

	return res

}
func StopRecordingdata(usecase string, applicationIdentifier string) *http.Response {
	url := appurl + "/devaten/data/stopRecording?usecaseIdentifier=" + usecase + "&inputSource=application&frocefullyStop=false"
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
	// applicationIdentifier1 := applicationIdentifier
	// applicationIdentifier1 = strings.Replace(applicationIdentifier1, "\n", "", -1)

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
	//fmt.Println(string(body))
	//responsecode := gjson.Get(string(body), "responseCode").Int()
	//if responsecode == 200 {
	monitoring.RecordReport(body)
	//}

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
	//value10 :=gjson.Get(res,"data.#.columnName")
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

	//var resultdatagot []string

	var resultdata []string
	value10 := gjson.Get(string(body), "data.#.columnName").Array()
	//value10 := gjson.Parse(string(body)).Get("data").Array()
	for _, v := range value10 {
		resultdata = append(resultdata, v.Str)
	}

	// var results []map[string]interface{}
	// json.Unmarshal([]byte(body), &results)
	// for key, result := range results {
	// 	data := result["data"].(map[string]interface{})
	// 	fmt.Println("Reading Value for Key :", key, data)
	// 	//Reading each value by its key
	// 	//for key1, result1 := range data {

	// 	//fmt.Println(data["columnName"])
	// 	//}
	// }

	monitoring.CreateStopMetrics(resultdata)

	return res

}

/*
The function that is called when the user starts the recording
Will every 5 seconds do the run operation, which returns some information about the current recording
*/
func scrapeWithInterval(command storage.StartAndStopCommand) {
	for {
		select {
		case <-quit:
			return
		default:
			Operation(command.Usecase, "run", command.ApplicationIdentifier)
		}
		time.Sleep(5 * time.Second)

	}
}

// Takes a username and a password and generates the string that is needed to login
func generateUserInfo(username string, password string) string {
	var userInfo = "username=" + username + "&password=" + password + "&grant_type=password"

	return userInfo
}
