package suitex

import (
	"testing"

	pb "github.com/vagudza/anti-brute-force/api/proto"
)

type Suite struct {
	T                   *testing.T
	AntiBrutforceClient pb.AntiBruteforceClient
}
