// Original template https://gist.github.com/MakoTano/624fe3fdea914b262e2c
package main

import (
	"encoding/base64"
	"fmt"
	"os"
	"regexp"
	"time"

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

	interval := c.Uint("interval")

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
	subscriptionsService := pubsubService.Projects.Subscriptions

	pullRequest := &pubsub.PullRequest{
		ReturnImmediately: false,
		MaxMessages:       1,
	}
	for {
		res, err := subscriptionsService.Pull(fqn, pullRequest).Do()
		if err != nil {
			fmt.Printf("Failed to pull from %v cause of %v\n", fqn, err)
			return err
		}

		for _, recvMsg := range res.ReceivedMessages {
			m := recvMsg.Message
			var decodedData string
			decoded, err := base64.StdEncoding.DecodeString(m.Data)
			if err != nil {
				decodedData = fmt.Sprintf("Failed to decode data by base64 because of %v", err)
			} else {
				decodedData = string(decoded)
			}
			fmt.Printf("%v %s: %v %s\n", m.PublishTime, m.MessageId, m.Attributes, decodedData)
			ackRequest := &pubsub.AcknowledgeRequest{
				AckIds: []string{recvMsg.AckId},
			}
			_, err = subscriptionsService.Acknowledge(fqn, ackRequest).Do()
			if err != nil {
				fmt.Printf("Failed to Acknowledge to %v cause of %v\n", fqn, err)
				return err
			}
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}
}
