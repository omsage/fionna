package android_util

import (
	"fionna/android/gadb"
	"fmt"
	"strings"
	"testing"
)

var (
	client gadb.Client
)

func SetClient() {
	client, _ = gadb.NewClient()
}

func TestGetDevice(t *testing.T) {
	SetClient()
	device, err := GetDevice(client, "231341")
	if err != nil {
		panic(err)
	}
	fmt.Println(device.Serial())
	data, _ := device.RunShellCommand("ps")
	fmt.Println(data)
}

func TestGetPackageNameList(t *testing.T) {
	device, err := GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	packageList, err := GetPackageNameList(device)
	if err != nil {
		panic(err)
	}
	fmt.Println(strings.Join(packageList, "\n"))
}

func TestGetCurrentPackageName(t *testing.T) {
	device, err := GetDevice(client, "")
	if err != nil {
		panic(err)
	}
	packageName, err := GetCurrentPackageName(device)
	if err != nil {
		panic(err)
	}
	fmt.Println(packageName)
}

func TestGetSerialList(t *testing.T) {
	serialList, err := GetSerialList(client)
	if err != nil {
		panic(err)
	}
	fmt.Println(serialList)
}
