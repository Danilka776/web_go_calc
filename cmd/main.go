package main

import application "sample-app/Desktop/service_GO/web_go_calc/app"

func main() {
	app := application.New()
	app.RunServer()
}
