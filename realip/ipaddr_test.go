package realip

import (
	"bytes"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/require"
	"net"
	"testing"
)

func TestInc(t *testing.T) {
	ip := net.ParseIP("192.168.0.1")
	res := Inc(ip)

	require.Equal(t, "192.168.0.2", res.String())

	ip = net.ParseIP("192.168.0.255")
	res = Inc(ip)

	require.Equal(t, "192.168.1.0", res.String())
}

func TestDec(t *testing.T) {
	ip := net.ParseIP("192.168.0.1")
	res := Dec(ip)

	require.Equal(t, "192.168.0.0", res.String())

	ip = net.ParseIP("192.168.4.0")
	res = Dec(ip)

	require.Equal(t, "192.168.3.255", res.String())
}

func TestSubnets(t *testing.T) {
	_, ipn, err := net.ParseCIDR("192.168.0.0/16")
	require.Nil(t, err)

	firstSub, secondSub := Subnets(ipn)
	require.Equal(t, "192.168.0.0/17", firstSub.String())
	require.Equal(t, "192.168.128.0/17", secondSub.String())
}

func TestIPNetEqual(t *testing.T) {
	_, first, err := net.ParseCIDR("192.168.0.0/16")
	require.Nil(t, err)

	_, second, err := net.ParseCIDR("192.168.0.0/16")
	require.Nil(t, err)
	require.True(t, IPNetEqual(first, second))

	_, second, err = net.ParseCIDR("192.167.0.0/16")
	require.Nil(t, err)
	require.False(t, IPNetEqual(first, second))

	_, second, err = net.ParseCIDR("192.168.0.0/15")
	require.Nil(t, err)
	require.False(t, IPNetEqual(first, second))
}

func TestIPCompareable(t *testing.T) {
	first := net.ParseIP("192.168.0.1")
	second := net.ParseIP("192.168.0.2")

	require.True(t, bytes.Compare(first, second) < 0)
	require.True(t, bytes.Compare(second, first) > 0)
}

func TestContainsSubnet(t *testing.T) {
	_, first, err := net.ParseCIDR("192.168.0.0/16")
	require.Nil(t, err)
	_, second, err := net.ParseCIDR("192.168.0.0/17")
	require.Nil(t, err)
	require.True(t, ContainsSubnet(first, second))

	_, first, err = net.ParseCIDR("192.168.0.0/16")
	require.Nil(t, err)
	_, second, err = net.ParseCIDR("192.168.128.0/17")
	require.Nil(t, err)
	require.True(t, ContainsSubnet(first, second))

	_, first, err = net.ParseCIDR("192.169.0.0/16")
	require.Nil(t, err)
	_, second, err = net.ParseCIDR("192.168.128.0/17")
	require.Nil(t, err)
	require.False(t, ContainsSubnet(first, second))

	_, first, err = net.ParseCIDR("192.168.0.0/16")
	require.Nil(t, err)
	_, second, err = net.ParseCIDR("192.168.128.0/15")
	require.Nil(t, err)
	require.False(t, ContainsSubnet(first, second))

	_, first, err = net.ParseCIDR("192.168.0.0/16")
	require.Nil(t, err)
	_, second, err = net.ParseCIDR("192.168.0.0/16")
	require.Nil(t, err)
	require.True(t, ContainsSubnet(first, second))
}

func TestExcludeSubnet(t *testing.T) {
	_, first, err := net.ParseCIDR("2000::/16")
	require.Nil(t, err)

	log.Debug().Msgf("Net: %+v %+v", first.IP, first.Mask)

	_, second, err := net.ParseCIDR("2000:1000::/32")
	require.Nil(t, err)

	res := ExcludeSubnet(first, second)

	for _, item := range res {
		fmt.Printf("%+v\n", item)
	}
}
