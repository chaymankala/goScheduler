package main

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func setupFiber() {
	app := fiber.New()
	app.Use(cors.New())
	app.Use(recover.New())
	// app.Use(jsonFilter)
	setupRoutes(app)
	log.Fatal(app.Listen(":3009"))
}

func setupRoutes(app *fiber.App) {
	app.Get("/ping", ping)
	app.Post("/jobs/create", createJob)
	app.Post("/jobs/remove", removeJobApi)
	app.Post("/jobs/update", updateAJob)
	app.Get("/jobs/list", listJobs)
}

func ping(c *fiber.Ctx) error {
	c.SendString("Bathike Unna!")
	return nil
}

func jsonFilter(c *fiber.Ctx) error {
	if c.Is("json") {
		return c.Next()
	}
	return send400Response(c)
}

func createJob(c *fiber.Ctx) error {
	type PMeta struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	payload := struct {
		JobName    string  `json:"jobName"`
		JobTime    string  `json:"time"`
		ActionType int8    `json:"actionType"`
		Meta       []PMeta `json:"meta"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return send400Response(c)
	}
	if len(payload.JobName) < 1 || len(payload.JobTime) < 1 || payload.ActionType < 1 {
		return send400Response(c)
	}
	jobTime, err := time.Parse(time.RFC3339, payload.JobTime)
	if err != nil {
		return send400Response(c)
	}

	var newJob NewJob
	newJob.JobName = payload.JobName
	newJob.JobTime = uint64(jobTime.Unix())
	newJob.ActionType = payload.ActionType
	for i := 0; i < len(payload.Meta); i++ {
		var newMeta JobMeta
		newMeta.Key = payload.Meta[i].Key
		newMeta.Value = payload.Meta[i].Value
		newJob.Meta = append(newJob.Meta, newMeta)
	}

	collectionJobs := getCollection(getJobsCollectionName())
	_, insertErr := collectionJobs.InsertOne(context.TODO(), newJob)
	if insertErr != nil {
		return send500Response(c)
	}
	c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status": "success",
	})
	return nil
}

func removeJobApi(c *fiber.Ctx) error {
	payload := struct {
		JobId string `json:"jobId"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return send400Response(c)
	}
	if len(payload.JobId) < 1 {
		return send400Response(c)
	}

	objId, err := primitive.ObjectIDFromHex(payload.JobId)

	if err != nil {
		return send400Response(c)
	}

	collectionJobs := getCollection(getJobsCollectionName())
	_, removeErr := collectionJobs.DeleteOne(context.TODO(), bson.M{"_id": objId})

	if removeErr != nil {
		return send500Response(c)
	}
	c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status": "success",
	})
	return nil
}

func listJobs(c *fiber.Ctx) error {
	collectionJobs := getCollection(getJobsCollectionName())
	findOptions := options.Find()
	skip, _ := strconv.ParseInt(c.Params("skip", "0"), 10, 64)
	limit, _ := strconv.ParseInt(c.Params("limit", "20"), 10, 64)
	findOptions.SetSkip(skip)
	findOptions.SetLimit(limit)
	cur, findErr := collectionJobs.Find(context.TODO(), bson.D{{}})
	if findErr != nil {
		return send500Response(c)
	}
	var jobs []*Job
	// Iterate through the cursor
	for cur.Next(context.TODO()) {
		var elem Job
		err := cur.Decode(&elem)
		if err != nil {
			continue
		}
		jobs = append(jobs, &elem)
	}
	c.JSON(fiber.Map{
		"jobs": jobs,
	})
	return nil
}

func updateAJob(c *fiber.Ctx) error {
	type PMeta struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}

	payload := struct {
		JobId      string  `json:"jobId"`
		JobName    string  `json:"jobName"`
		JobTime    string  `json:"time"`
		ActionType int8    `json:"actionType"`
		Meta       []PMeta `json:"meta"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return send400Response(c)
	}
	if len(payload.JobId) < 1 {
		return send400Response(c)
	}

	collectionJobs := getCollection(getJobsCollectionName())
	id, _ := primitive.ObjectIDFromHex(payload.JobId)

	updateQinner := bson.M{
		"job_name":    payload.JobName,
		"action_type": payload.ActionType,
		"job_time":    payload.JobTime,
		"meta":        payload.Meta,
	}

	if len(payload.JobName) < 1 {
		delete(updateQinner, "job_name")
	}

	if payload.ActionType < 1 {
		delete(updateQinner, "action_type")
	}

	if len(payload.JobTime) < 1 {
		delete(updateQinner, "job_time")
	}

	if len(payload.Meta) < 1 {
		delete(updateQinner, "meta")
	}

	updateQ := bson.M{
		"$set": updateQinner,
	}

	_, err := collectionJobs.UpdateOne(context.TODO(), bson.M{"_id": id}, updateQ)

	if err != nil {
		return send500Response(c)
	}
	c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"status": "success",
	})
	return nil
}
