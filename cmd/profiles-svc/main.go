package main

import "gitlab.com/samkomarov/profiles-svc.git/internal"

func main() {
	cfg := internal.ReadConfigFromEnv()
	internal.InitializeAndStart(cfg)
}
