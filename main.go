package main

func main() {

	loadEnv()
	initDB()
	run()
	setupFiber()
	runner()
}
