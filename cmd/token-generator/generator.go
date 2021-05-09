package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/weijunji/go-lottery/pkgs/auth"
	"github.com/weijunji/go-lottery/pkgs/utils"
	"os"
	"time"
)

func main() {
	f, err := os.OpenFile("test_token.csv", os.O_CREATE | os.O_WRONLY, 0600)
	if err != nil {
		log.Info(err)
	}

	for i := 1; i < 20000; i++ {
		token, _ := utils.GenerateToken(uint64(i), auth.RoleNormal, time.Hour * 24000)
		f.WriteString(token + "\n")
	}
}
