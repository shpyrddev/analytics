package github

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/go-github/github"
)

type GitHubHandler struct {
	db driver.Conn
}

func New() *GitHubHandler {

	conn, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{fmt.Sprintf("%s:%d", "localhost", 9000)},
		Auth: clickhouse.Auth{
			Database: "analytics",
			Username: "default",
		},
		Settings: clickhouse.Settings{
			"flatten_nested": 0,
		},
	})
	if err != nil {

		fmt.Printf("Error: Unable to connect to clickhouse, exiting. err=%s", err)
		os.Exit(1)
	}

	fmt.Println("Connected to clickhouse")
	gh := &GitHubHandler{
		db: conn,
	}

	http.HandleFunc("/github", gh.Handle)

	return gh

}

func (g *GitHubHandler) Handle(w http.ResponseWriter, req *http.Request) {
	payload, err := io.ReadAll(req.Body)
	if err != nil {
		fmt.Printf("Error: Unable to read request body: err=%s\n", err)
	}
	defer req.Body.Close()

	ctx := context.Background()
	batch, err := g.db.PrepareBatch(ctx, "INSERT INTO github_events")
	if err != nil {
		fmt.Printf("Error: Unable to prepare batch: err=%s\n", err)
		return
	}

	event, err := github.ParseWebHook(github.WebHookType(req), payload)
	if err != nil {
		fmt.Printf("Error: Unable to parse: err=%s\n", err)
		return
	}

	switch event := event.(type) {
	case *github.PushEvent:

		err = batch.AppendStruct(event)
		if err != nil {
			fmt.Printf("Error: Unable to append batch: err=%s\n", err)
			return
		}
	}

	err = batch.Send()
	if err != nil {
		fmt.Printf("Error: Unable to send batch: err=%s\n", err)
		return
	}

	fmt.Println(batch.IsSent())

	fmt.Println("PUSH INSERTED")
}
