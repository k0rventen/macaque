package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"time"

	cron "github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	slack "github.com/slack-go/slack"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	kubernetes "k8s.io/client-go/kubernetes"
	rest "k8s.io/client-go/rest"
)

// configuration for the macaque
type MacaqueConfig struct {
	selector      string
	namespace     string
	crontab       string
	podname       string
	slack_token   string
	slack_channel string
	timezone      string
}

// returns true if slack variables are not empty
func (c *MacaqueConfig) HasSlack() bool {
	if c.slack_channel != "" && c.slack_token != "" {
		return true
	}
	return false
}

func (c *MacaqueConfig) IsValid() bool {
	_, cronNotValid := parseCron(c.crontab)
	if cronNotValid != nil || c.namespace == "" {
		return false
	}
	return true
}

// parse the config from env var and returns a MacaqueConfig struct
func parseConfig() (MacaqueConfig, error) {
	crontab := os.Getenv("MACAQUE_CRONTAB")
	selector := os.Getenv("MACAQUE_SELECTOR")
	namespace := os.Getenv("MACAQUE_NAMESPACE")
	slack_token := os.Getenv("MACAQUE_SLACK_TOKEN")
	slack_channel := os.Getenv("MACAQUE_SLACK_CHANNEL")
	timezone := os.Getenv("MACAQUE_TIMEZONE")
	pod_name := os.Getenv("HOSTNAME") // this one is set by kube

	c := MacaqueConfig{
		crontab:       crontab,
		namespace:     namespace,
		selector:      selector,
		podname:       pod_name,
		timezone:      timezone,
		slack_token:   slack_token,
		slack_channel: slack_channel,
	}
	return c, nil
}

// parse the cron spec, returns a cron.Schedule or err if it fails
func parseCron(cronSpec string) (cron.Schedule, error) {
	cronParser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
	recc, err := cronParser.Parse(cronSpec)
	if err != nil {
		log.Error("cron format is invalid: ", err.Error())
		return nil, err
	}
	return recc, nil
}

// list the pods corresponding to the criterias (ns, selector)
func listPods(c MacaqueConfig, k *kubernetes.Clientset) ([]v1.Pod, error) {
	pods, err := k.CoreV1().Pods(c.namespace).List(context.TODO(), metav1.ListOptions{LabelSelector: c.selector})
	if err != nil {
		return nil, err
	}
	var podsList []v1.Pod
	//filter our own pod out & pods not running
	for _, v := range pods.Items {
		if v.Name != c.podname && v.Status.Phase == "Running" {
			podsList = append(podsList, v)
		}
	}
	return podsList, nil
}

// goroutine that handles pods listing and killing.
func podKiller(conf MacaqueConfig, ch chan bool, slack_ch chan string) {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Error("Unable to retrieve in-cluster configuration", err.Error())
		slack_ch <- err.Error()

		// this is mandatory, so crash now
		os.Exit(1)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Error("Unable to connect to the cluster", err.Error())
		slack_ch <- err.Error()

		// this is mandatory, so crash now
		os.Exit(1)
	}
	log.Info("Waiting for a termination order..")
	for {
		// wait for the signal
		<-ch
		podsList, err := listPods(conf, clientset)
		if err != nil {
			// the error might be a RBAC related problem, but that can be changed without restarting the pod
			log.Error(err.Error())
			slack_ch <- err.Error()
		} else {
			if len(podsList) == 0 {
				log.Warn("No pods were found..")
				slack_ch <- "No pods were found, therefore none were killed.."
			} else {
				log.Debug("Selecting a random pod..")
				seed := rand.NewSource(time.Now().UnixNano())
				rng := rand.New(seed)
				index := rng.Intn(len(podsList))
				choosen_pod := podsList[index]

				err := clientset.CoreV1().Pods(conf.namespace).Delete(context.TODO(), choosen_pod.Name, metav1.DeleteOptions{})
				log.Info("Pod ", choosen_pod.Name, " has been selected. !")
				if err != nil {
					log.Error(err.Error())
				}
				slack_ch <- "Pod " + choosen_pod.Name + " has been terminated."
			}
		}
	}
}

// sends message from the channel to slack
func slackSender(conf MacaqueConfig, ch chan string) {
	for {
		msg := <-ch
		if conf.HasSlack() {
			api := slack.New(conf.slack_token)
			_, _, _, err := api.SendMessage(conf.slack_channel, slack.MsgOptionText(msg, false))
			if err != nil {
				log.Error(err.Error())
			}
		}
	}
}

// computes the next occurence for the cron spec and waits until then
func sleepUntilNextCron(conf MacaqueConfig, sch cron.Schedule) {
	loc, err := time.LoadLocation(conf.timezone)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
	t := time.Now().In(loc)
	next_t := sch.Next(t)
	clean_next_t := next_t.Round(time.Second)
	delta := next_t.Sub(t)
	clean_delta := delta.Round(time.Second)
	log.Info("next occurence at ", clean_next_t, ", sleeping for ", clean_delta)
	time.Sleep(delta)
}

func setLogging() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
	})
}

func main() {
	fmt.Print("\no(..)o  Starting macaque v0.2 !\n  (-) _/\n\n")
	// init everything
	setLogging()

	config, err := parseConfig()
	if err != nil {
		log.Fatal("Unable to load config: ", err.Error())
	}

	schedule, err := parseCron(config.crontab)
	if err != nil {
		log.Fatal("Unable to parse cron format: ", err.Error())
	}
	log.Info("cron schedule is '", config.crontab, "'")

	kill_channel := make(chan bool)
	message_channel := make(chan string)
	go podKiller(config, kill_channel, message_channel)
	go slackSender(config, message_channel)

	for {
		// wait for the next cron occurence, then notify the killer routine
		sleepUntilNextCron(config, schedule)
		kill_channel <- true
	}
}
