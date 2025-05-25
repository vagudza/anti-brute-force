package test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/test/suitex"
)

func TestBlacklistAPI(t *testing.T) {
	ctx, s := suitex.New(t)

	subnetsToAdd := []string{
		"192.168.1.0/24",
		"10.0.0.0/8",
		"172.16.0.0/12",
	}

	t.Run("Add and remove subnets", func(t *testing.T) {
		for _, subnet := range subnetsToAdd {
			_, err := s.AntiBruteforceClient.AddToBlacklist(ctx, &pb.IPSubnetRequest{Subnet: subnet})
			require.NoError(t, err, "Failed to add subnet %s to blacklist", subnet)
		}

		resp, err := s.AntiBruteforceClient.GetBlacklist(ctx, &pb.EmptyRequest{})
		require.NoError(t, err)

		sort.Strings(subnetsToAdd)
		retrievedSubnets := resp.GetSubnets()
		sort.Strings(retrievedSubnets)

		require.Equal(t, subnetsToAdd, retrievedSubnets, "Blacklist should contain all added subnets")

		singleSubnet := subnetsToAdd[0]
		_, err = s.AntiBruteforceClient.RemoveFromBlacklist(ctx, &pb.IPSubnetRequest{Subnet: singleSubnet})
		require.NoError(t, err)

		respRemove, err := s.AntiBruteforceClient.GetBlacklist(ctx, &pb.EmptyRequest{})
		require.NoError(t, err)
		require.NotContains(t, respRemove.GetSubnets(), singleSubnet)
	})
}
