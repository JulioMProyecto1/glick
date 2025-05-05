package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"

	"glick/config"
	// "io"
	"net/http"
	"os"
	"strings"
	"time"
)

type TaskPayload struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	Assignees     []int  `json:"assignees"`
	DueDate       int64  `json:"due_date"`
	StartDate     int64  `json:"start_date,omitempty"`
	DueDateTime   bool   `json:"due_date_time,omitempty"`
	StartDateTime bool   `json:"start_date_time,omitempty"`
	Status        string `json:"status,omitempty"`
}

type StartTrackTimePayload struct {
	Tid string `json:"tid"`
}

type CreateTrackTimePayload struct {
	Tid      string `json:"tid"`
	Start    int64  `json:"start"`
	Duration int64  `json:"duration"`
}

func main() {

	config.LoadEnv()

	list_id := config.GetEnv("LIST_ID")
	assignee := config.GetEnv("ASSIGNEE_CLICKUP_ID")
	assignee_int, _ := strconv.Atoi(assignee)
	url := fmt.Sprintf("https://api.clickup.com/api/v2/list/%s/task", list_id)

	reader := bufio.NewReader(os.Stdin)
	today := time.Now().UnixNano() / int64(time.Millisecond)
	cu_task_name, _ := reader.ReadString('\n')
	cu_description, _ := reader.ReadString('\n')
	cu_tracktime, _ := reader.ReadString('\n')
	cu_tracktime_text := strings.TrimSpace(cu_tracktime)

	mode := ""

	if cu_tracktime_text == "y" {
		mode = "START_TRACK"
	} else if cu_tracktime_text != "" {
		mode = "WITH_TRACK"
	} else if cu_tracktime_text == "" {
		mode = "NO_TRACK"
	}

	taskData := TaskPayload{
		Name:        strings.TrimSpace(cu_task_name),
		Description: strings.TrimSpace(cu_description),
		Assignees:   []int{assignee_int},
		DueDate:     today,
	}
	switch mode {
	case "START_TRACK":

		task_response := post(url, taskData)
		start_time_url := "https://api.clickup.com/api/v2/team/529/time_entries/start"
		track_time_data := StartTrackTimePayload{
			Tid: task_response.Id,
		}

		post(start_time_url, track_time_data)
		os.Exit(0)
	case "WITH_TRACK":
		taskData.DueDateTime = true
		taskData.StartDateTime = true
		taskData.Status = "Done"
		duration, err_parse := time.ParseDuration(cu_tracktime_text)
		duration_in_millis := duration.Milliseconds()
		taskData.StartDate = taskData.DueDate - duration_in_millis
		task_response := post(url, taskData)
		if err_parse != nil {
			fmt.Println("Invalid duration", err_parse)
			os.Exit(1)
		}
		create_time_url := "https://api.clickup.com/api/v2/team/529/time_entries"
		create_time_tracked_data := CreateTrackTimePayload{
			Tid:      task_response.Id,
			Start:    taskData.StartDate,
			Duration: duration_in_millis,
		}
		post(create_time_url, create_time_tracked_data)
		os.Exit(0)
	case "NO_TRACK":
		post(url, taskData)
		os.Exit(0)
	default:
		fmt.Println("Nothing to do")
		os.Exit(0)
	}

}

type TaskResponse struct {
	Id string `json:"id"`
}

func post(url string, body any) TaskResponse {

	clickup_api_key := config.GetEnv("CLICKUP_API_KEY")
	jsonBytes, err := json.Marshal(body)
	fmt.Println("Request body:", string(jsonBytes))
	if err != nil {
		fmt.Println("Error: ", err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonBytes))

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", clickup_api_key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("HTTP request error:", err)
	}

	defer res.Body.Close()

	// bodyBytes, _ := io.ReadAll(res.Body)
	// fmt.Println("Raw response:", string(bodyBytes))
	var taskRes TaskResponse
	if err := json.NewDecoder(res.Body).Decode(&taskRes); err != nil {
		fmt.Println("Error decoding response:", err)
	}

	return taskRes

}
