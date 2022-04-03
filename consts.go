package main

import (
	"time"
)

const dbConnectTimeout = 10 * time.Second
const ENV_SCRAPE_INTERVAL = "SCRAPE_INTERVAL"
const ENV_DB_URI = "DB_URI"
const ENV_DB_NAME = "DB_NAME"
const ENV_COLLECTION_NAME = "COLLECTION_NAME"
