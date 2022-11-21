package plugin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/drone/drone-go/plugin/webhook"
	"io"
	"net/http"
	"os"
	"strconv"
)

func New(param1, param2 string) webhook.Plugin {
	return &plugin{
		param1: param1,
		param2: param2,
	}
}

type plugin struct {
	param1 string
	param2 string
}
type requestBody struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Icon  string `json:"icon"`
	Group string `json:"group"`
}

func sendCard(status string, repoName string, repoLink string, commit string, build string) {
	webhookUrl := os.Getenv("PLUGIN_WEBHOOK")
	contentType := "application/json"
	reqBody := requestBody{repoName + "Build" + status, "", "https://cdn.orluma.ltd/midway/drone.png", "Drone"}
	req, _ := json.Marshal(reqBody)
	resp, err := http.Post(webhookUrl, contentType, bytes.NewBuffer(req))
	if err != nil {
		panic(err)
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)
	if resp.StatusCode == http.StatusCreated {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		jsonStr := string(body)
		fmt.Println("Response: ", jsonStr)
	} else {
		body, _ := io.ReadAll(resp.Body)
		fmt.Println("Get failed with error: ", resp.Status, string(body))
	}
}
func (p *plugin) Deliver(ctx context.Context, req *webhook.Request) error {
	if req.Event == webhook.EventBuild {
		link := req.Repo.Link
		slug := req.Repo.Slug
		commit := req.Build.Link
		action := req.Action
		build := slug + "/" + strconv.Itoa(int(req.Build.Number))
		if action == webhook.ActionCreated {
			sendCard(action, slug, link, commit, build)
		}
		if action == webhook.ActionUpdated {
			if req.Build.Status == "success" {
				fmt.Println("success")
				sendCard(req.Build.Status, slug, link, commit, build)
			}
			if req.Build.Status == "failure" {
				sendCard(req.Build.Status, slug, link, commit, build)
			}
		}
	}
	return nil
}
