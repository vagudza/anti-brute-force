package suitex

import (
	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/internal/config"
)

type Suite struct {
	Cfg                  *config.AppConfig
	AntiBruteforceClient pb.AntiBruteforceClient
}
