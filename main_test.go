package main

import (
	"net"
	"reflect"
	"testing"
)

func TestEncode(t *testing.T) {
	ips, err := encode("vtrust-flash", "flashmeifyoucan")
	if err != nil {
		t.Fatal(err)
	}

	var want []net.IP
	for _, ip := range []string{
		"226.0.1.2",
		"226.1.2.3",
		"226.2.3.4",
		"226.32.15.15",
		"226.64.61.54",
		"226.65.32.51",
		"226.66.56.60",
		"226.67.62.51",
		"226.68.32.62",
		"226.69.46.53",
		"226.70.60.63",
		"226.71.0.48",
		"226.96.97.65",
		"226.97.245.78",
		"226.16.12.12",
		"226.48.116.118",
		"226.49.117.114",
		"226.50.116.115",
		"226.51.102.45",
		"226.52.97.108",
		"226.53.104.115",
		"226.80.204.136",
		"226.81.94.247",
		"226.112.35.35",
	} {
		want = append(want, net.ParseIP(ip))
	}

	if !reflect.DeepEqual(want, ips) {
		t.Fatalf("want:\n%#v\ngot:\n%#v", want, ips)
	}
}
