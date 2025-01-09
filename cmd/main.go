package main

import (
	"os"
	"time"
	_ "time/tzdata"

	"github.com/mikalai2006/kingwood-api/internal/app"
)

func main() {
	// base path for config: default = ./ (for test ../)
	const configPath = "./"
	// go func() {
	// 	s, _ := gocron.NewScheduler()
	// 	defer func() { _ = s.Shutdown() }()

	// 	_, _ = s.NewJob(
	// 		gocron.DurationJob(
	// 			time.Second*5,
	// 		),
	// 		gocron.NewTask(
	// 			func() {
	// 				fmt.Println("Cron run!")
	// 			},
	// 		),
	// 	)
	// }()
	loc, err := time.LoadLocation("Europe/Minsk")
	if err != nil {
		panic("Wrong timezone")
	}
	time.Local = loc
	os.Setenv("TZ", "Europe/Minsk")

	app.Run(configPath)

}
