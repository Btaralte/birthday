package main

import (
	"birthdayreminder/api"
	"birthdayreminder/config"
	"birthdayreminder/worker"
	"log"

	"github.com/robfig/cron/v3"
)

func main() {
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	server, err := api.NewServer(config)
	if err != nil {
		log.Fatal(err)
	}
	// istLoc, _ := time.LoadLocation("Asia/Kolkata") // IST timezone
	// c := cron.New(cron.WithLocation(istLoc))             // Scheduler uses IST
	//_, err = c.AddFunc("0 9 * * *", worker.RunDailyTask) // At 9:00 AM IST

	c := cron.New(cron.WithSeconds())
	_, err = c.AddFunc("*/5 * * * * *", worker.RunDailyTask) // At 9:00 AM IST
	if err != nil {
		log.Fatal("Error scheduling the daily task:", err)
	}
	c.Start()
	server.Start()
	defer c.Stop()
}
