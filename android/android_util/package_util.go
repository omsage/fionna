package android_util

import (
	"fionna/android/gadb"
	"fionna/entity"
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
)

func GetPackageNameList(device *gadb.Device) ([]string, error) {
	output, err := device.RunShellCommand("pm list packages")
	if err != nil {
		return nil, fmt.Errorf("exec command erro : " + "adb shell pm list packages")
	}
	// 解析输出并提取包名列表
	packages := extractPackageNames(output)
	return packages, nil
}

func extractPackageNames(output string) []string {
	lines := strings.Split(output, "\n")
	packages := []string{}

	regex := regexp.MustCompile(`^package:(.*)$`)

	for _, line := range lines {
		// 每行的格式为 "package:<package_name>"
		match := regex.FindStringSubmatch(line)
		if len(match) == 2 {
			pkg := match[1]
			pkg = strings.ReplaceAll(pkg, "\n", "")

			pkg = strings.TrimSpace(pkg)

			if len(pkg) <= 0 {
				continue
			}
			if strings.Contains(pkg, "com.android") {
				continue
			}
			packages = append(packages, pkg)
		}
	}
	return packages
}

func GetCurrentPackageName(device *gadb.Device) (string, error) {
	data, err := device.RunShellCommand("dumpsys activity top | grep ACTIVITY")
	if err != nil {
		return "", fmt.Errorf("exec command erro : " + "dumpsys activity top | grep ACTIVITY")
	}

	var dataSplit []string

	dataSplitTemp := strings.Split(data, "\n")

	for _, lineStr := range dataSplitTemp {
		if lineStr != "" {
			dataSplit = append(dataSplit, lineStr)
		}
	}

	currentActivityLineStr := strings.TrimLeft(dataSplit[len(dataSplit)-1], " ")

	dataSplit = strings.Split(currentActivityLineStr, " ")
	if len(dataSplit) > 1 {
		dataSplit = strings.Split(dataSplit[1], "/")
		return dataSplit[0], nil
	} else {
		return "", errors.New("not get current package name")
	}
}

func GetSerialList(client gadb.Client) ([]entity.SerialInfo, error) {
	devices, err := client.DeviceList()
	if err != nil {
		panic(err)
	}
	var serialList []entity.SerialInfo
	for _, dev := range devices {
		if stat, _ := dev.State(); stat == gadb.StateOnline {
			serialList = append(serialList, entity.SerialInfo{
				SerialName:    dev.Serial(),
				Model:         dev.Model(),
				ProductDevice: dev.Product(),
			},
			)
		}
	}
	return serialList, err
}

func GetSerialInfo(device *gadb.Device) entity.SerialInfo {
	productDevice := device.Product()
	serial := device.Serial()
	version, _ := device.RunShellCommand("getprop ro.system.build.version.release")
	name, _ := device.RunShellCommand("getprop ro.product.name")
	model, _ := device.RunShellCommand("getprop ro.product.model")

	manufacturer, _ := device.RunShellCommand("getprop ro.product.manufacturer")

	cpu, _ := device.RunShellCommand("getprop ro.product.cpu.abi")
	var isHm int

	ringtone, _ := device.RunShellCommand("getprop ro.config.ringtone")

	if ringtone != "" && strings.Contains(ringtone, "Harmony") {
		version, _ = device.RunShellCommand("getprop hw_sc.build.platform.version")
		isHm = 1
	} else {
		isHm = 0
	}
	s := entity.SerialInfo{
		SerialName:    serial,
		ProductDevice: productDevice,
		Version:       strings.TrimSpace(version),
		Platform:      1,
		Name:          strings.TrimSpace(name),
		Model:         strings.TrimSpace(model),
		Manufacturer:  strings.TrimSpace(manufacturer),
		IsHm:          isHm,
		CPUArch:       strings.TrimSpace(cpu), //
		Size:          GetDeviceSize(device),
	}
	return GetBattery(device, s)
}

func GetDeviceSize(device *gadb.Device) string {
	var size = ""
	res, _ := device.RunShellCommand("wm size")
	if strings.Contains(res, "Override size") {
		size = res[:strings.Index(res, "Override size")]
	} else {
		sizeSplit := strings.Split(res, ":")
		if len(sizeSplit) > 1 {
			size = sizeSplit[1]
		}
	}
	if size == "" {
		return size
	}
	size = strings.TrimSpace(size)
	size = strings.ReplaceAll(size, ":", "")
	size = strings.ReplaceAll(size, "\r", "")
	size = strings.ReplaceAll(size, "\n", "")
	size = strings.ReplaceAll(size, " ", "")
	if len(size) > 20 {
		size = "unknown"
	}
	return size
}

func GetBattery(device *gadb.Device, serial entity.SerialInfo) entity.SerialInfo {
	battery, _ := device.RunShellCommand("dumpsys battery")
	battery = strings.ReplaceAll(battery, "Max charging voltage", "")
	batterySplit := strings.Split(battery, "\n")
	for _, line := range batterySplit {
		if strings.Contains(line, "level") {
			levelSplit := strings.Split(line, ":")
			if len(levelSplit) > 1 {
				num, _ := strconv.ParseInt(strings.TrimSpace(levelSplit[1]), 10, 64)
				serial.Level = int(num)
			}
		}
		if strings.Contains(line, "voltage") {
			levelSplit := strings.Split(line, ":")
			if len(levelSplit) > 1 {
				num, _ := strconv.ParseInt(strings.TrimSpace(levelSplit[1]), 10, 64)
				serial.Voltage = int(num)
			}
		}
		if strings.Contains(line, "temperature") {
			levelSplit := strings.Split(line, ":")
			if len(levelSplit) > 1 {
				num, _ := strconv.ParseInt(strings.TrimSpace(levelSplit[1]), 10, 64)
				serial.Temperature = int(num)
			}
		}

	}
	return serial
}

func GetCpuCores(device *gadb.Device) int {
	resStr, err := device.RunShellCommand("cat /proc/cpuinfo | grep processor | wc -l")
	if err != nil {
		return 0
	}
	resStr = strings.TrimSpace(resStr)
	resStr = strings.ReplaceAll(resStr, "\n", "")

	num, err := strconv.Atoi(resStr)
	if err != nil {
		return 0
	}
	return num
}

func GetDevice(client gadb.Client, serial string) (*gadb.Device, error) {
	devices, err := client.DeviceList()
	if err != nil {
		return nil, err
	}

	if len(devices) == 0 {
		return nil, errors.New("not connect device")
	}
	if serial == "" {
		return &devices[0], nil
	} else {
		for _, dev := range devices {
			if dev.Serial() == serial {
				return &dev, nil
			}
		}
		return nil, errors.New("not connect device")
	}
}

func GetPidOnPackageName(device *gadb.Device, appName string) (pid string, err error) {
	data, err := device.RunShellCommand("dumpsys activity " + appName)
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(fmt.Sprintf("ACTIVITY\\s+%s.*pid=[0-9]+", appName))
	regResult := reg.FindString(data)
	dataSplit := strings.Split(regResult, " ")
	if len(dataSplit) < 2 {
		return "", errors.New("unable to find the pid corresponding to app")
	}
	return strings.ReplaceAll(dataSplit[3], "pid=", ""), nil
}
