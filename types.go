package main

import "go.mongodb.org/mongo-driver/bson/primitive"

type JobMeta struct {
	Key   string `bson:"key"`
	Value string `bson:"value"`
}

type NewJob struct {
	JobName    string    `bson:"job_name"`
	JobTime    uint64    `bson:"job_time"`
	ActionType int8      `bson:"action_type"`
	Meta       []JobMeta `bson:"meta"`
}

type Job struct {
	Id         primitive.ObjectID `bson:"_id"`
	JobName    string             `bson:"job_name"`
	JobTime    uint64             `bson:"job_time"`
	ActionType int8               `bson:"action_type"`
	Meta       []JobMeta          `bson:"meta"`
}
