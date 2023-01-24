package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/google/go-github/v49/github"
	"github.com/mitchellh/colorstring"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
)

func main() {
	if err := NewGhDeployCommand().Execute(); err != nil {
		os.Exit(1)
	}
}

type GhDeployRunner struct {
	Debug       bool
	Repository  string
	Environment string
	Tag         string
}

func NewGhDeployCommand() *cobra.Command {
	gdr := GhDeployRunner{}
	root := &cobra.Command{
		Use:     "gh-deploy",
		Short:   "Create deployments in Github",
		Long:    "Create deployments in Github",
		Version: "VERSION",
		PreRun:  gdr.preRun,
		RunE:    gdr.RunE,
	}
	root.Flags().StringVarP(&gdr.Repository, "repo", "r", "", "repository name")
	_ = root.MarkFlagRequired("repo")
	root.Flags().StringVarP(&gdr.Environment, "env", "e", "", "environment name")
	_ = root.MarkFlagRequired("env")
	root.Flags().StringVarP(&gdr.Tag, "tag", "t", "", "tag name")
	_ = root.MarkFlagRequired("tag")
	return root
}

func (gdr *GhDeployRunner) preRun(c *cobra.Command, args []string) {
	logrus.SetOutput(os.Stderr)
	if gdr.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	if gdr.Environment == "production" {
		colorstring.Print("[yellow]CAUTION: You will create a new deployment in production environment. Do you want to continue? [y/N]: ")
		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		t := s.Text()
		if t != "y" && t != "Y" {
			colorstring.Println("[red]This deployment was canceled.")
			os.Exit(0)
		}
	}
}

func (gdr *GhDeployRunner) RunE(c *cobra.Command, args []string) error {
	ownerRepo := strings.Split(gdr.Repository, "/")
	if len(ownerRepo) != 2 {
		return fmt.Errorf("can't parse an option repo (expect owner/repo)")
	}
	owner := ownerRepo[0]
	repo := ownerRepo[1]

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)
	cli := github.NewClient(tc)
	autoMerge := false
	req := github.DeploymentRequest{
		Ref:              &gdr.Tag,
		Environment:      &gdr.Environment,
		RequiredContexts: &[]string{},
		AutoMerge:        &autoMerge,
	}
	logrus.Infoln("Create deployment")
	deploy, _, err := cli.Repositories.CreateDeployment(context.Background(), owner, repo, &req)
	if err != nil {
		return err
	}
	logrus.Infof("DeployID: %d, DeployStatusesURL: %s", deploy.GetID(), deploy.GetStatusesURL())
	return nil
}
