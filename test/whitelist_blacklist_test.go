package test

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
	pb "github.com/vagudza/anti-brute-force/api/proto"
	"github.com/vagudza/anti-brute-force/test/suitex"
)

func TestWhiteAndBlacklistAPI(t *testing.T) {
	ctx, s := suitex.New(t)
	clearLists(ctx, t, s)

	const subnetsCount = 5
	subnetsToAdd := make([]string, 0, subnetsCount)

	for range subnetsCount {
		subnetsToAdd = append(subnetsToAdd, generateRandomSubnet(t))
	}

	t.Run("add and remove subnets", func(t *testing.T) {
		{
			t.Log("adding subnets to whitelist and blacklist")

			for _, subnet := range subnetsToAdd {
				_, err := s.AntiBruteforceClient.AddToWhitelist(ctx, &pb.IPSubnetRequest{Subnet: subnet})
				require.NoError(t, err, "Failed to add subnet %s to whitelist", subnet)

				_, err = s.AntiBruteforceClient.AddToBlacklist(ctx, &pb.IPSubnetRequest{Subnet: subnet})
				require.NoError(t, err, "Failed to add subnet %s to blacklist", subnet)
			}
		}

		{
			t.Log("check whitelist and blacklist contents")

			respWhitelist, err := s.AntiBruteforceClient.GetWhitelist(ctx, &pb.EmptyRequest{})
			require.NoError(t, err)

			respBlacklist, err := s.AntiBruteforceClient.GetBlacklist(ctx, &pb.EmptyRequest{})
			require.NoError(t, err)

			whitelistSubnets := respWhitelist.GetSubnets()
			sort.Strings(whitelistSubnets)

			blacklistSubnets := respBlacklist.GetSubnets()
			sort.Strings(blacklistSubnets)

			sort.Strings(subnetsToAdd)
			require.Equal(t, subnetsToAdd, whitelistSubnets, "whitelist should contain all added subnets")
			require.Equal(t, subnetsToAdd, blacklistSubnets, "blacklist should contain all added subnets")
		}

		{
			t.Log("removing subnet from whitelist and blacklist and check state")
			singleSubnet := subnetsToAdd[0]

			_, err := s.AntiBruteforceClient.RemoveFromWhitelist(ctx, &pb.IPSubnetRequest{Subnet: singleSubnet})
			require.NoError(t, err)

			_, err = s.AntiBruteforceClient.RemoveFromBlacklist(ctx, &pb.IPSubnetRequest{Subnet: singleSubnet})
			require.NoError(t, err)

			respWhitelist, err := s.AntiBruteforceClient.GetWhitelist(ctx, &pb.EmptyRequest{})
			require.NoError(t, err)
			require.NotContains(t, respWhitelist.GetSubnets(), singleSubnet)

			respBlacklist, err := s.AntiBruteforceClient.GetBlacklist(ctx, &pb.EmptyRequest{})
			require.NoError(t, err)
			require.NotContains(t, respBlacklist.GetSubnets(), singleSubnet)
		}
	})
}
