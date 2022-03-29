package main

import (
	"bufio"
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Abhiraman-V/Prbot/src"
	"github.com/bradleyfalzon/ghinstallation"
	"github.com/gin-gonic/gin"
	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
)

var (
	client *github.Client
	err    error
)

func init() {
	godotenv.Load()

}

func NewClient() (*github.Client, error) {
	appId, _ := strconv.Atoi(os.Getenv("APP_ID"))
	installId, _ := strconv.Atoi(os.Getenv("INSTALL_ID"))
	privateKey := os.Getenv("PRIVATE_KEY")

	trans := http.DefaultTransport

	gh, err := ghinstallation.NewKeyFromFile(trans, int64(appId), int64(installId), privateKey)

	if err != nil {
		return &github.Client{}, err
	}

	return github.NewClient(&http.Client{Transport: gh}), nil

}

func main() {
	client, err = NewClient()

	gin.SetMode(gin.ReleaseMode)

	if err != nil {
		panic(err)
	}

	app := gin.Default()
	// app := gin.New()
	app.POST("/webhook", handleEvent)

	app.Run("127.0.0.1:8080")
	fmt.Println("7")

}

func handleEvent(c *gin.Context) {
	event := c.GetHeader("X-GitHub-Event")
	payload, _ := ioutil.ReadAll(c.Request.Body)

	// fmt.Println(verifySign(payload, c.GetHeader("X-Hub-Signature-256")))

	if verifySign(payload, c.GetHeader("X-Hub-Signature-256")) && len(payload) > 0 {

		switch event {
		case "pull_request":
			fmt.Println("Pull Request")
			handlePushEvent(payload)
		}
	}

}

func handlePushEvent(payload []byte) {

	pr := PullRequest{}

	if err := json.Unmarshal(payload, &pr); err != nil {
		panic(err)
	}

	if pr.Action == "opened" || pr.Action == "reopened" {
		fmt.Println("Opened PR")
		files := pr.ChangedFilesFromPullRequest(client)

		for _, f := range files {
			if filepath.Ext(f.GetFilename()) == ".yaml" && f.GetStatus() == "modified" {
				if val, n := isValid(f); !val {

					body := fmt.Sprintf("@%s policy does not validate with the base",
						pr.PullRequestData.User.Login)

					comment := &github.PullRequestComment{
						Body:     github.String(body),
						Path:     f.Filename,
						CommitID: github.String(pr.PullRequestData.Head.Sha),
						Position: github.Int(n),
					}
					pr.CommentOnPullRequest(client, comment)
					return
				}
			}
		}
		pr.MergePullRequest(client)
	}

	// pr := github.PullRequestReviewCommentEvent{}

	// acti := *pr.Action

	// fmt.Printf("%s", acti)

}

func verifySign(payload []byte, sign string) bool {
	key := hmac.New(sha256.New, []byte(os.Getenv("WEBHOOK_SECRET")))
	key.Write([]byte(string(payload)))
	computedSign := "sha256=" + hex.EncodeToString(key.Sum(nil))

	fmt.Println("computed=", computedSign)
	fmt.Println("sign=", sign)

	return computedSign == sign
}

func isValid(file *github.CommitFile) (bool, int) {

	buff := bytes.NewBuffer([]byte(file.GetPatch()))
	scanner := bufio.NewScanner(buff)
	count := 0
	for scanner.Scan() {
		line := scanner.Text()
		count++
		if strings.HasPrefix(line, "-") {
			return false, count
		}
	}
	return true, count
}

/*Get all Changed files in github pull request*/
func (pr *PullRequest) ChangedFilesFromPullRequest(client *github.Client) []*github.CommitFile {

	files, _, err := client.PullRequests.ListFiles(context.Background(), pr.PullRequestData.Head.User.Login,
		pr.PullRequestData.Head.Repo.Name, int(pr.Number), &github.ListOptions{})

	if err != nil {
		log.Println(err.Error())
		return []*github.CommitFile{}
	}
	return files
}

/*Merge the pull request*/
func (pr *PullRequest) MergePullRequest(client *github.Client) (*github.PullRequestMergeResult, error) {

	log.Printf("Merging pull request %d on %s",
		int(pr.Number),
		pr.PullRequestData.Base.Repo.Name)

	result, _, err := client.PullRequests.Merge(context.Background(),
		pr.PullRequestData.Base.User.Login,
		pr.PullRequestData.Base.Repo.Name, int(pr.Number),
		fmt.Sprintf("Merging pull request %d", pr.Number),
		&github.PullRequestOptions{})

	return result, err

}

/*closes the pull request*/
func (pr *PullRequest) ClosePullRequest(client *github.Client) {

	log.Printf("Closing pull request %d on %s",
		int(pr.Number),
		pr.PullRequestData.Base.Repo.Name)

	result, res, err := client.PullRequests.Edit(context.Background(), pr.PullRequestData.Base.User.Login,
		pr.PullRequestData.Base.Repo.Name, int(pr.Number), &github.PullRequest{
			State: github.String("closed"),
		})
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(res.StatusCode)
	log.Println(*result.State)
}

func (pr PullRequest) CommentOnPullRequest(client *github.Client, comment *github.PullRequestComment) {

	_, res, err := client.PullRequests.CreateComment(context.Background(),
		pr.PullRequestData.Base.User.Login,
		pr.PullRequestData.Base.Repo.Name, int(pr.Number),
		comment)

	if err != nil {
		log.Println(err.Error())
	}
	log.Println(res.Status)

}
