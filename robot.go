package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/opensourceways/community-robot-lib/config"
	framework "github.com/opensourceways/community-robot-lib/robot-gitee-framework"
	sdk "github.com/opensourceways/go-gitee/gitee"
	"github.com/sirupsen/logrus"
)

const (
	botName      = "keeper-approve"
	mergeCommand = "MERGE"
)

type Owners struct {
	Approvers    []string `yaml:"approvers"`
	Reviewers    []string `yaml:"reviewers"`
	BranchKeeper []string `yaml:"branch_keeper"`
}

var commandReg = regexp.MustCompile(`(?m)^/([^\s]+)[\t ]*([^\n\r]*)`)

type iClient interface {
	GetPRLabels(org, repo string, number int32) ([]sdk.Label, error)
	GetBot() (sdk.User, error)
	AddPRLabel(org, repo string, number int32, label string) error
	RemovePRLabel(org, repo string, number int32, label string) error
}

func newRobot(cli iClient, botName string, token string) *robot {
	return &robot{cli: ghclient{cli}, botName: botName, token: token}
}

type robot struct {
	token   string
	cli     ghclient
	botName string
}

func (bot *robot) NewConfig() config.Config {
	return &configuration{}
}

func (bot *robot) RegisterEventHandler(f framework.HandlerRegitster) {
	f.RegisterNoteEventHandler(bot.handleNoteEvent)
}

func (bot *robot) handleNoteEvent(e *sdk.NoteEvent, cnf config.Config, log *logrus.Entry) error {
	if !e.IsPullRequest() {
		log.Info("Event is not a creation of a comment on a PR, skipping.")
		return nil
	}

	if !e.IsCreatingCommentEvent() {
		log.Info("Event is not a creation of a comment on an open PR, skipping.")
		return nil
	}

	if !isMergeCommand(e.Comment.Body) {
		log.Info("Event is not a merge comment, skipping.")
		return nil
	}

	c, ok := cnf.(*configuration)
	if !ok {
		return fmt.Errorf("can't convert to configuration")
	}

	tagName := c.TagName
	if tagName == "" {
		tagName = "keeper_approved"
	}

	org, repo := e.GetOrgRepo()
	number := e.GetPRNumber()

	var keepers []string
	token := bot.token
	err := loadOwnersInfo(org, repo, token, &keepers)
	if err != nil {
		return nil
	}

	commenter := e.GetCommenter()

	if !isBranchKeeper(commenter, keepers) {
		log.Info("Commenter is not a branch keeper, skipping.")
		return nil
	} else {
		return bot.cli.AddLabel(org, repo, number, tagName)
	}
}

func loadOwnersInfo(org, repo, token string, keeper *[]string) error {
	url := fmt.Sprintf("https://gitee.com/api/v5/repos/%s/%s/raw/OWNERS?access_token=%s&ref=br_develop_mindie", org, repo, token)
	resp, err := http.Get(url)
	if err != nil {
		logrus.Info("request owners file failure...")
		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logrus.Info("Error reading request body")
		return err
	}

	var owners Owners

	defer resp.Body.Close()
	err = yaml.Unmarshal(body, &owners)

	if err != nil {
		logrus.Info("Error unmarshalling body")
		return err
	}
	*keeper = owners.BranchKeeper
	return nil
}

func isBranchKeeper(commenter string, keeps []string) bool {
	for _, ele := range keeps {
		if ele == commenter {
			return true
		}
	}
	return false
}

func isMergeCommand(comment string) bool {
	for _, match := range commandReg.FindAllStringSubmatch(comment, -1) {
		cmd := strings.ToUpper(match[1])

		if cmd == mergeCommand {
			return true
		}
	}

	return false
}
