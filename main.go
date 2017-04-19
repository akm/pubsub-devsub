// Original template https://gist.github.com/MakoTano/624fe3fdea914b262e2c
package main

import (
	"fmt"
	"os"
	"time"

	"github.com/urfave/cli"

	"golang.org/x/net/context"
	"google.golang.org/api/iterator"
	"cloud.google.com/go/pubsub"
)

func main() {
	app := cli.NewApp()
	app.Name = "pubsub-devsub"
	app.Usage = "github.com/groovenauts/pubsub-devsub"
	app.Version = Version

  app.Flags = []cli.Flag {
    cli.StringFlag{
      Name: "project",
      Usage: "GCS Project ID",
			EnvVar: "GCP_PROJECT,PROJECT",
    },
    cli.StringFlag{
      Name: "subscription",
      Usage: "Subscription",
			EnvVar: "SUBSCRIPTION",
    },
    cli.UintFlag{
      Name: "interval",
      Value: 10,
      Usage: "Interval to pull",
    },
  }

	app.Action = executeCommand

	app.Run(os.Args)
}

func executeCommand(c *cli.Context) error {
	proj := c.String("project")
	subscription := c.String("subscription")
	if proj == "" || subscription == "" {
		cli.ShowAppHelp(c)
		os.Exit(1)
	}
	interval := c.Uint("interval")

	ctx := context.Background()
	client, err := pubsub.NewClient(ctx, proj)
	if err != nil {
		fmt.Println("Failed to get new pubsub client for ", proj, " cause of ", err)
		os.Exit(1)
	}
	sub := client.Subscription(subscription)
	for ;; {
		it, err := sub.Pull(ctx)
		if err != nil {
			fmt.Println("Failed to pull from ", subscription, " cause of ", err)
			os.Exit(1)
		}
		// Ensure that the iterator is closed down cleanly.
		defer it.Stop()

		for ;; {
			m, err := it.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Println("Failed to get pulled message from ", subscription, " cause of ", err)
				break
			}
			fmt.Printf("%v %s: %v %s\n", m.PublishTime, m.ID, m.Attributes, m.Data)
			m.Done(true)
		}

		time.Sleep(time.Duration(interval) * time.Second)
	}

	return nil
}
