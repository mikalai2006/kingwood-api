package main

import (
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

	app.Run(configPath)

}
