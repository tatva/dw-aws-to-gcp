package main

import (
	"fmt"
	"errors"
	"io"
	"bufio"
	"bytes"
	"time"
	"os"
	"os/exec"
	"strings"
	"github.com/mozilla-services/heka/message"
	. "github.com/mozilla-services/heka/pipeline"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/net/context"
	"google.golang.org/cloud/pubsub"
)

type PubsubOutputConfig struct {
	ProjectId string `toml:"project_id"`
	TopicName string `toml:"topic_name"`
	AppCredentials string `toml:"app_credentials"`
}

type PubsubOutput struct {
	config *PubsubOutputConfig
	client *pubsub.Client
	project_id string
	topic_name string
}


func (pso *PubsubOutput) publish(msg string) {
	topic := pso.topic_name
	message := msg
	msgIDs, err := client.Topic(topic).Publish(context.Background(), &pubsub.Message{
		Data: []byte(message),
	})
	if err != nil {
		log.Fatalf("Publish failed, %v", err)
	}
	fmt.Printf("Message '%s' published to topic %s and the message id is %s\n", message, topic, msgIDs[0])
}



func (pso *PubsubOutput) ConfigStruct() interface{} {
	return &PubsubOutputConfig{}
}

func (pso *PubsubOutput) Init(config interface{}) (err error) {
	pso.config = config.(*PubsubOutputConfig)
	
	pso.client = pubsub.NewClient(context.Background(), "ll-poc")
	pso.project_id = pso.config.ProjectId
	pso.topic_name = pso.config.TopicName
	return
}

func (pso *S3Output) Run(or OutputRunner, h PluginHelper) (err error) {
	inChan := or.InChan()
		
	var (
		pack    *PipelinePack
		msg     *message.Message
		ok      = true
	)

	for ok {
		select {
		case pack, ok = <- inChan:
			if !ok {
				break
			}
			msg = pack.Message
			pso.publish(msg)
			if err != nil {
				or.LogMessage(fmt.Sprintf("Warning, unable to send to PubSub: %s", err))
				err = nil
				continue
			}
			pack.Recycle(nil)
		}
	}

	or.LogMessage(fmt.Sprintf("Shutting down PubSub output runner."))
	return
}

func init() {
	RegisterPlugin("PubsubOutput", func() interface{} {
		return new(PubsubOutput)
	})
}

func main() {
	fmt.Println("In here")
}
