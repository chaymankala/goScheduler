package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var currentJobs []*Job

func run() {
	scrapeNow()
	stopped := make(chan bool, 1)
	ticker := time.NewTicker(getScrapeInterval())

	go func() {
		for {
			select {
			case <-ticker.C:
				scrapeNow()
			case <-stopped:
				ticker.Stop()
				return
			}
		}
	}()
}

func getScrapeInterval() time.Duration {
	i, err := strconv.ParseInt(getEnvVariable(ENV_SCRAPE_INTERVAL), 10, 64)
	if err != nil {
		log.Fatal("Scrape Interval is not valid")
	}
	return time.Duration(i) * time.Second
}

func getScrapeIntervalRaw() int64 {
	i, err := strconv.ParseInt(getEnvVariable(ENV_SCRAPE_INTERVAL), 10, 64)
	if err != nil {
		log.Fatal("Scrape Interval is not valid")
	}
	return i
}

func scrapeNow() {
	fmt.Println("scraping")
	collection := getCollection(getJobsCollectionName())
	findOptions := options.Find()
	timeMax := time.Now().Unix() + getScrapeIntervalRaw()
	findQuery := bson.M{
		"job_time": bson.M{
			"$gte": time.Now().Unix(),
			"$lt":  timeMax,
		},
	}
	findOptions.SetLimit(100)

	cur, err := collection.Find(context.TODO(), findQuery, findOptions)
	if err != nil {
		// TODO move to failed jobs
		fmt.Println(err)
	}

	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var job Job
		err := cur.Decode(&job)
		if err != nil {
			// TODO move to failed jobs
			fmt.Println(err)
		}
		currentJobs = append(currentJobs, &job)
	}

	if err := cur.Err(); err != nil {
		// log.Fatal(err)
		// TODO move to failed jobs
		fmt.Println(err)
	}

	cur.Close(context.TODO())
}

func runner() {
	for {
		if len(currentJobs) > 0 {
			for i := 0; i < len(currentJobs); i++ {
				if currentJobs[i].ActionType == 1 {
					doActionType(1, currentJobs[i])
				}
				//  else {
				// 	//TODO move to failedjobs
				// }
			}
		}
	}
}

func removeJob(id primitive.ObjectID) {
	r := -1
	for i := 0; i < len(currentJobs); i++ {
		if primitive.ObjectID.String(currentJobs[i].Id) == primitive.ObjectID.String(id) {
			r = i
			break
		}
	}
	if r != -1 {
		currentJobs = append(currentJobs[:r], currentJobs[r+1:]...)
	}
	fmt.Println("Removed job")
}
