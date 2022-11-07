package main

import (
	"fmt"
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/networking/v2/ports"
	"net"
	"os"
	"strings"
)

func main() {
	fmt.Printf("osportpatch .....\n")

	if len(os.Args) < 4 {
		usage()
	}

	var remove bool = false
	if strings.ToLower(os.Args[1]) == "remove" {
		remove = true
	} else if strings.ToLower(os.Args[1]) != "add" {
		usage()
	}
	ipaddr := os.Args[2]
	if net.ParseIP(ipaddr) == nil {
		fmt.Printf("Invalid IP address (%s)\n", ipaddr)
		usage()
	}

	opts := gophercloud.AuthOptions{
		IdentityEndpoint: getEnvChecked("OS_AUTH_URL"),
		Username:         getEnvChecked("OS_USERNAME"),
		Password:         getEnvChecked("OS_PASSWORD"),
		TenantName:       getEnvChecked("OS_PROJECT_NAME"),
		DomainName:       "Default",
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		panic(err)
	}

	computeClient, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")})
	if err != nil {
		panic(err)
	}

	networkClient, err := openstack.NewNetworkV2(provider, gophercloud.EndpointOpts{Region: os.Getenv("OS_REGION_NAME")})

	serverNameSet := make(map[string]bool)
	for _, n := range os.Args[3:] {
		serverNameSet[n] = true
	}

	serverList, err := listServers(computeClient)
	if err != nil {
		panic(err)
	}

	portList, err := listPorts(networkClient)
	if err != nil {
		panic(err)
	}

	fmt.Printf("osportpatch found %d servers and %d ports in this tenant\n", len(serverList), len(portList))

	serverById := make(map[string]servers.Server)

	for _, srv := range serverList {
		if _, ok := serverNameSet[srv.Name]; ok {
			//fmt.Printf("Server %s to be handled\n", srv.Name)
			serverById[srv.ID] = srv
		} else {
			//fmt.Printf("Server %s NOT to be handled\n", srv.Name)
		}
	}

	for _, port := range portList {
		if server, ok := serverById[port.DeviceID]; ok {
			newAddressPairs, err := handlePort(networkClient, port, ipaddr, remove)
			if err != nil {
				panic(err)
			}
			fmt.Printf("Server '%s': others allowed IP:%v)\n", server.Name, newAddressPairs)
		}
	}
}

func getEnvChecked(vname string) string {
	s := os.Getenv(vname)
	if s == "" {
		fmt.Printf("\nERROR: Missing \"%s\" environment variable\n", vname)
		os.Exit(2)
	}
	return s
}

func handlePort(client *gophercloud.ServiceClient, p ports.Port, ipAddress string, remove bool) ([]ports.AddressPair, error) {
	var newAddressPairs []ports.AddressPair
	if remove {
		newAddressPairs = removeIp(p.AllowedAddressPairs, ipAddress)
	} else {
		newAddressPairs = addIp(p.AllowedAddressPairs, ipAddress)
	}
	//fmt.Printf("newAddressPairs:%v\n", newAddressPairs)
	r := ports.Update(client, p.ID, ports.UpdateOpts{AllowedAddressPairs: &newAddressPairs})
	_, err := r.Extract()
	return newAddressPairs, err
}

func addIp(src []ports.AddressPair, ipAddress string) []ports.AddressPair {
	newAddressPairs := make([]ports.AddressPair, 0, len(src)+1)
	found := false
	for _, oldap := range src {
		newAddressPairs = append(newAddressPairs, oldap)
		if oldap.IPAddress == ipAddress {
			found = true
		}
	}
	if !found {
		newAddressPairs = append(newAddressPairs, ports.AddressPair{
			IPAddress: ipAddress,
		})
	}
	return newAddressPairs
}

func removeIp(src []ports.AddressPair, ipAddress string) []ports.AddressPair {
	newAddressPairs := make([]ports.AddressPair, 0, len(src)+1)
	for _, oldap := range src {
		if oldap.IPAddress != ipAddress {
			newAddressPairs = append(newAddressPairs, oldap)
		}
	}
	return newAddressPairs
}

func listServers(client *gophercloud.ServiceClient) ([]servers.Server, error) {

	page, err := servers.List(client, servers.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}
	if b, _ := page.IsEmpty(); b {
		return []servers.Server{}, nil
	}

	serverList, err := servers.ExtractServers(page)
	if err != nil {
		return nil, err
	}
	return serverList, nil
}

func listPorts(client *gophercloud.ServiceClient) ([]ports.Port, error) {
	page, err := ports.List(client, ports.ListOpts{}).AllPages()
	if err != nil {
		return nil, err
	}
	if b, _ := page.IsEmpty(); b {
		return []ports.Port{}, nil
	}
	portList, err := ports.ExtractPorts(page)
	if err != nil {
		return nil, err
	}
	return portList, err
}

func usage() {
	fmt.Printf("USAGE: osportpatch [add|remove] <ipaddr> <server> [<server>...]\n")
	os.Exit(1)
}
