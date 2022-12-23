package ips

import (
	"net"
	"strconv"
	"testing"
)

func Test_IPPortEqual(t *testing.T) {
	t.Parallel()

	tcs := []struct {
		ip1    IPPort
		ip2    IPPort
		expect bool
	}{
		{
			ip1:    IPPort{IP: net.ParseIP("127.0.0.1"), Port: 1},
			ip2:    IPPort{IP: net.ParseIP("127.0.0.1"), Port: 1},
			expect: true,
		},
		{
			ip1:    IPPort{IP: net.ParseIP("::1"), Port: 1},
			ip2:    IPPort{IP: net.ParseIP("::1"), Port: 1},
			expect: true,
		},
		{
			ip1:    IPPort{IP: net.ParseIP("127.0.0.1"), Port: 1},
			ip2:    IPPort{IP: net.ParseIP("127.0.0.1"), Port: 2},
			expect: false,
		},
	}

	for i, tc := range tcs {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			result := tc.ip1.Equal(tc.ip2)
			if result != tc.expect {
				t.Errorf("extect %v but return %v\n", tc.expect, result)
			}
		})
	}
}
