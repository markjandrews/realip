package realip

import (
	"bytes"
	"github.com/rs/zerolog/log"
	"net"
)

type IPHostMask []byte

func checkIPv4(ip net.IP) net.IP {
	if x := ip.To4(); x != nil {
		return x
	}

	return ip
}

func Inc(ip net.IP) net.IP {
	ip = checkIPv4(ip)

	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]++
		if ip[i] > 0 {
			break
		}
	}

	return ip
}

func Dec(ip net.IP) net.IP {
	ip = checkIPv4(ip)

	for i := len(ip) - 1; i >= 0; i-- {
		ip[i]--
		if ip[i] < 255 {
			break
		}
	}

	return ip
}

func validateIPv4(ip net.IP) net.IP {
	if x := ip.To4(); x != nil {
		return x
	}

	return ip
}

func BroadcastAddress(network *net.IPNet) net.IP {
	ip := make(net.IP, len(validateIPv4(network.IP)))
	for i :=0; i <len(ip); i++{
		ip[i] = ^network.Mask[i] | network.IP[i]
	}

	return ip
}

func Subnets(network *net.IPNet) (first *net.IPNet, second *net.IPNet) {
	prefix, bits := network.Mask.Size()
	newMask := net.CIDRMask(prefix + 1, bits)

	first = &net.IPNet{
		IP:   network.IP,
		Mask: newMask,
	}

	secondIP := Inc(BroadcastAddress(first))
	second = &net.IPNet{
		IP:   secondIP,
		Mask: newMask,
	}
	return
}

func IPNetEqual(network *net.IPNet, other *net.IPNet) bool {
	if network.IP.Equal(other.IP) &&
		bytes.Compare(network.Mask, other.Mask) == 0{
		return true
	}

	return false
}

func ContainsSubnet(network *net.IPNet, other *net.IPNet) bool {
	if bytes.Compare(network.IP, other.IP) <= 0 &&
		bytes.Compare(BroadcastAddress(network), BroadcastAddress(other)) >= 0 {
		return true
	}

	return false
}

func ExcludeSubnet(network *net.IPNet, other *net.IPNet) []*net.IPNet {
	network.IP = validateIPv4(network.IP)
	other.IP = validateIPv4(other.IP)

	if len(network.IP) != len(other.IP) {
		log.Panic().Msg("IP Version mismatch")
	}

	var res []*net.IPNet

	s1, s2 := Subnets(network)
	for !IPNetEqual(s1, other) && !IPNetEqual(s2, other) {
		if ContainsSubnet(s1, other) {
			res = append(res, s2)
			s1, s2 = Subnets(s1)
		} else if ContainsSubnet(s2, other){
			res = append(res, s1)
			s1, s2 = Subnets(s2)
		}
	}

	if IPNetEqual(s1, other) {
		res = append(res, s2)
	} else if IPNetEqual(s2, other) {
		res = append(res, s1)
	}

	return res
}