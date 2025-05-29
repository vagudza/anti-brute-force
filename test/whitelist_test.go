package test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"

	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/test/suitex"
)

func TestWhitelistAPI(t *testing.T) {
	ctx, s := suitex.New(t)

	// Clean blacklist before tests
	_, err := s.AntiBruteforceClient.ClearWhitelist(ctx, &pb.EmptyRequest{})
	require.NoError(t, err)

	const subnetsCount = 5
	subnetsToAdd := make([]string, 0, subnetsCount)

	for range subnetsCount {
		subnetsToAdd = append(subnetsToAdd, generateRandomSubnet(t))
	}

	t.Run("Add and remove subnets", func(t *testing.T) {
		for _, subnet := range subnetsToAdd {
			_, err := s.AntiBruteforceClient.AddToWhitelist(ctx, &pb.IPSubnetRequest{Subnet: subnet})
			require.NoError(t, err, "Failed to add subnet %s to whitelist", subnet)
		}

		resp, err := s.AntiBruteforceClient.GetWhitelist(ctx, &pb.EmptyRequest{})
		require.NoError(t, err)

		sort.Strings(subnetsToAdd)
		retrievedSubnets := resp.GetSubnets()
		sort.Strings(retrievedSubnets)

		require.Equal(t, subnetsToAdd, retrievedSubnets, "whitelist should contain all added subnets")

		singleSubnet := subnetsToAdd[0]
		_, err = s.AntiBruteforceClient.RemoveFromWhitelist(ctx, &pb.IPSubnetRequest{Subnet: singleSubnet})
		require.NoError(t, err)

		respRemove, err := s.AntiBruteforceClient.GetWhitelist(ctx, &pb.EmptyRequest{})
		require.NoError(t, err)
		require.NotContains(t, respRemove.GetSubnets(), singleSubnet)
	})
}
