package osrs

import (
	"fmt"
	"github.com/ctfloyd/hazelmere-commons/pkg/hz_client"
	"net/url"
)

type Skill struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Level int    `json:"level"`
	Xp    int    `json:"xp"`
}

type Activity struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Rank  int    `json:"rank"`
	Score int    `json:"score"`
}

type Hiscore struct {
	Skills     []Skill    `json:"skills"`
	Activities []Activity `json:"activities"`
}

type HiscoreClient struct {
	httpClient *hz_client.HttpClient
}

func NewHiscoreClient(httpClient *hz_client.HttpClient) *HiscoreClient {
	return &HiscoreClient{httpClient}
}

func (hc *HiscoreClient) GetHiscore(username string) (Hiscore, error) {
	path := fmt.Sprintf("%s?player=%s", hc.httpClient.GetHost(), url.PathEscape(username))
	var hiscore Hiscore
	err := hc.httpClient.Get(path, &hiscore)
	if err != nil {
		return Hiscore{}, err
	}
	return hiscore, nil
}
