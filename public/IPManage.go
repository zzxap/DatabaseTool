package public

import (
	"net"
	//"sort"
	"strconv"
	"strings"
)

func IsLocalIp(ip net.IP) bool {
	for _, cidr := range []string{"10.0.0.0/8", "172.16.0.0/12", "192.168.0.0/16"} {
		_, subnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic("failed to parse hardcoded rfc1918 cidr: " + err.Error())
		}
		if subnet.Contains(ip) {
			return true
		}
	}
	return false
}

func rfc4193private(ip net.IP) bool {
	_, subnet, err := net.ParseCIDR("fd00::/8")
	if err != nil {
		panic("failed to parse hardcoded rfc4193 cidr: " + err.Error())
	}
	return subnet.Contains(ip)
}

func IsLoopback(ip net.IP) bool {
	for _, cidr := range []string{"127.0.0.0/8", "::1/128"} {
		_, subnet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic("failed to parse hardcoded loopback cidr: " + err.Error())
		}
		if subnet.Contains(ip) {
			return true
		}
	}
	return false
}

func IsPublicIp(ip net.IP) bool {
	return !IsLocalIp(ip) && !rfc4193private(ip) && !IsLoopback(ip)
}

type name_and_ip struct {
	string
	net.IP
}

func heuristic(ni name_and_ip) (ret int) {
	a := strings.ToLower(ni.string)
	ip := ni.IP
	if IsLoopback(ip) {
		ret += 1000
	}
	if IsLocalIp(ip) || rfc4193private(ip) {
		ret += 500
	}
	if strings.Contains(a, "dyn") {
		ret += 100
	}
	if strings.Contains(a, "dhcp") {
		ret += 99
	}
	for i := 0; i < len(ip); i++ {
		if strings.Contains(a, strconv.Itoa(int(ip[i]))) {
			ret += 5
		}
	}
	return ret
}

type nameAndIPByStabilityHeuristic []name_and_ip

func (nis nameAndIPByStabilityHeuristic) Len() int      { return len(nis) }
func (nis nameAndIPByStabilityHeuristic) Swap(i, j int) { nis[i], nis[j] = nis[j], nis[i] }
func (nis nameAndIPByStabilityHeuristic) Less(i, j int) bool {
	return heuristic(nis[i]) < heuristic(nis[j])
}

func PublicAddresses() ([]name_and_ip, error) {
	var ret []name_and_ip

	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		ip, _, err := net.ParseCIDR(addr.String())
		if err != nil {
			return nil, err
		}
		// ignore unresolvable addresses
		names, err := net.LookupAddr(ip.String())
		if err != nil {
			continue
		}
		for _, name := range names {
			//Log("name=" + name + "ip=" + ip)
			ret = append(ret, name_and_ip{name, ip})
		}
	}

	//sort.Sort(nameAndIPByStabilityHeuristic(ret))
	return ret, nil
}

func GetIpList() {
	if names, err := PublicAddresses(); err == nil {
		for _, ni := range names {
			println(ni.string, ni.IP.String(), heuristic(ni))
		}
	} else {
		//panic(err)
	}

}
