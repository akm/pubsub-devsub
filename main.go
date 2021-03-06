// Original template https://gist.github.com/MakoTano/624fe3fdea914b262e2c
package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"

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

	baseFlags := []cli.Flag{
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
		cli.BoolFlag{
			Name:  "return-immediately,r",
			Usage: "Return result immediately on pull",
		},
		cli.UintFlag{
			Name:  "max-messages,m",
			Value: 10,
			Usage: "Max messages per pull",
		},
		cli.BoolFlag{
			Name:  "verbose,V",
			Usage: "Show debug logs",
		},
	}
	sort.Sort(cli.FlagsByName(baseFlags))

	app.Flags = append(baseFlags,
		cli.BoolFlag{
			Name:  "follow,f",
			Usage: "Keep subscribing",
		},
		cli.BoolFlag{
			Name:  "ack,A",
			Usage: "Send ACK for received message",
		},
	)

	app.Commands = []cli.Command{
		{
			Name:    "inspect",
			Aliases: []string{"i"},
			Usage:   "Pull and show messages without ACK",
			Flags:   baseFlags,
			Action: func(c *cli.Context) error {
				puller := buildPuller(c)
				puller.Ack = false
				puller.Follow = false
				return puller.Run()
			},
		},
		{
			Name:    "subscribe",
			Aliases: []string{"s"},
			Usage:   "Subscribe and show message with ACK",
			Flags:   baseFlags,
			Action: func(c *cli.Context) error {
				puller := buildPuller(c)
				puller.Ack = true
				puller.Follow = true
				return puller.Run()
			},
		},
	}

	app.Action = func(c *cli.Context) error {
		puller := buildPuller(c)
		return puller.Run()
	}

	sort.Sort(cli.FlagsByName(app.Flags))
	sort.Sort(cli.CommandsByName(app.Commands))

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

func buildPuller(c *cli.Context) *Puller {
	fqn := buildFqn(c)

	ctx := context.Background()

	// https://github.com/google/google-api-go-client#application-default-credentials-example
	client, err := google.DefaultClient(ctx, pubsub.PubsubScope)
	if err != nil {
		fmt.Printf("Failed to google.DefaultClient with scope %v cause of %v\n", pubsub.PubsubScope, err)
		os.Exit(1)
	}

	// Creates a pubsubService
	pubsubService, err := pubsub.New(client)
	if err != nil {
		fmt.Printf("Failed to create pubsub.Service with client %v cause of %v\n", client, err)
		os.Exit(1)
	}

	puller := &Puller{
		SubscriptionsService: pubsubService.Projects.Subscriptions,
		Ack:                  c.Bool("ack"),
		Follow:               c.Bool("follow"),
		Fqn:                  fqn,
		Interval:             int(c.Uint("interval")),
		MaxMessages:          int64(c.Uint("max-messages")),
		ReturnImmediately:    c.Bool("return-immediately"),
		Verbose:              c.Bool("verbose"),
	}
	puller.Setup()
	return puller
}
