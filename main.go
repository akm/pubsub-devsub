// Original template https://gist.github.com/MakoTano/624fe3fdea914b262e2c
package main

import (
	"fmt"
	"os"
	"regexp"

	"github.com/urfave/cli"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"

	pubsub "google.golang.org/api/pubsub/v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "pubsub-devsub"
	app.Usage = "development tool for Google Cloud Pubsub"
	app.Version = Version

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "project",
			Usage:  "GCS Project ID",
			EnvVar: "GCP_PROJECT,PROJECT",
		},
		cli.BoolFlag{
			Name:  "follow,f",
			Usage: "Keep subscribing",
		},
		cli.UintFlag{
			Name:  "interval",
			Value: 10,
			Usage: "Interval to pull",
		},
	}

	app.Action = executeCommand

	app.Run(os.Args)
}

func buildFqn(c *cli.Context) string {
	if !c.Args().Present() {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}

	re := regexp.MustCompile("^projects/.+/subscriptions/.+$")
	// => [projects/proj1/subscriptions/sub1 proj1 sub1] for "projects/proj1/subscriptions/sub1"
	if re.MatchString(c.Args().First()) {
		return c.Args().First()
	} else {
		proj := c.String("project")
		if proj == "" {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}
		return fmt.Sprintf("projects/%s/subscriptions/%s", proj, c.Args().First())
	}
}

func executeCommand(c *cli.Context) error {
	fqn := buildFqn(c)

	ctx := context.Background()

	// https://github.com/google/google-api-go-client#application-default-credentials-example
	client, err := google.DefaultClient(ctx, pubsub.PubsubScope)
	if err != nil {
		fmt.Printf("Failed to google.DefaultClient with scope %v cause of %v\n", pubsub.PubsubScope, err)
		return err
	}

	// Creates a pubsubService
	pubsubService, err := pubsub.New(client)
	if err != nil {
		fmt.Printf("Failed to create pubsub.Service with client %v cause of %v\n", client, err)
		return err
	}

	puller := &Puller{
		SubscriptionsService: pubsubService.Projects.Subscriptions,
		Fqn:                  fqn,
		Interval:             int(c.Uint("interval")),
	}
	puller.Setup()
	if c.Bool("follow") {
		return puller.Follow()
	} else {
		return puller.Execute()
	}
}
