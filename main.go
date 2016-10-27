package main

import (
	"fmt"
)

func main() {
	//
	initEvent()

}

func initEvent() {

	eventList := os.Getenv("EVENT_LIST")

	RUN_ID = os.Getenv("RUN_ID")
	POD_NAME = os.Getenv("POD_NAME")
	SERVICE_ADDR = os.Getenv("SERVICE_ADDR")
	REGISTER_URL = os.Getenv("REGISTER_URL")

	for _, eventInfo := range strings.Split(eventList, ";") {
		if len(strings.Split(eventInfo, ",")) > 1 {
			eventKey := strings.Split(eventInfo, ",")[0]
			eventId := strings.Split(eventInfo, ",")[1]
			if os.Getenv(eventKey) != "" {
				event := new(EventInfo)
				event.id = eventId
				event.url = os.Getenv(eventKey)
				eventMap[eventKey] = *event
			}
		}
	}
}

func notifyEvent(eventKey, bodyType string, body io.Reader) {
	if event, ok := eventMap[eventKey]; ok {
		eventId := event.id
		notifyUrl := event.url
		if strings.Contains(notifyUrl, "?") {
			notifyUrl += "&runId=" + RUN_ID
		} else {
			notifyUrl += "?runId=" + RUN_ID
		}
		notifyUrl += "&event=" + eventKey + "&eventId=" + eventId

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodPut, notifyUrl, body)
		if err != nil {
			fmt.Println("error when create a notify request:" + err.Error())
		}

		resp, err := client.Do(req)

		if err != nil {
			fmt.Println("error when notify event", eventKey, body)
		}

		respBody, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Println("notify resp is", string(respBody))
	}
}

func waitForData(serviceInfo []string) {
	// send a message to pipelint to notify that current component is ready to receive data
	registerUrl := eventMap["REGISTER_URL"].url
	if strings.Contains(registerUrl, "?") {
		registerUrl += "&runId=" + RUN_ID
	} else {
		registerUrl += "?runId=" + RUN_ID
	}

	serviceAddr := ":" + serviceInfo[1]

	registerUrl += "&podName=" + POD_NAME
	registerUrl += "&receiveUrl=" + url.QueryEscape(serviceAddr+"/receivedata")

	client := &http.Client{}
	req, _ := http.NewRequest("PUT", registerUrl, nil)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(resp)

	http.HandleFunc("/receivedata", receiveDataHandler)
	http.ListenAndServe(":"+serviceInfo[2], nil)

}

func receiveDataHandler(w http.ResponseWriter, r *http.Request) {
	result, _ := json.Marshal(map[string]string{"message": "ok"})

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println("error when get data body:" + err.Error())
	}

	codePathMap := make(map[string]string)
	json.Unmarshal([]byte(body), &codePathMap)
	dataChan <- codePathMap["path"]

	w.Write(result)
}
