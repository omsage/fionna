package gadb

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

type DeviceFileInfo struct {
	Name         string
	Mode         os.FileMode
	Size         uint32
	LastModified time.Time
}

func (info DeviceFileInfo) IsDir() bool {
	return (info.Mode & (1 << 14)) == (1 << 14)
}

const DefaultFileMode = os.FileMode(0664)

type DeviceState string

const (
	StateUnknown      DeviceState = "UNKNOWN"
	StateOnline       DeviceState = "online"
	StateOffline      DeviceState = "offline"
	StateDisconnected DeviceState = "disconnected"
)

var deviceStateStrings = map[string]DeviceState{
	"":        StateDisconnected,
	"offline": StateOffline,
	"device":  StateOnline,
}

func deviceStateConv(k string) (deviceState DeviceState) {
	var ok bool
	if deviceState, ok = deviceStateStrings[k]; !ok {
		return StateUnknown
	}
	return
}

type DeviceForward struct {
	Serial string
	Local  string
	Remote string
	// LocalProtocol string
	// RemoteProtocol string
}

type Device struct {
	adbClient Client
	serial    string
	attrs     map[string]string
}

func (d Device) Product() string {
	return d.attrs["product"]
}

func (d Device) Model() string {
	return d.attrs["model"]
}

func (d Device) Usb() string {
	return d.attrs["usb"]
}

func (d Device) transportId() string {
	return d.attrs["transport_id"]
}

func (d Device) DeviceInfo() map[string]string {
	return d.attrs
}

func (d Device) Serial() string {
	// 	resp, err := d.adbClient.executeCommand(fmt.Sprintf("host-serial:%s:get-serialno", d.serial))
	return d.serial
}

func (d Device) IsUsb() bool {
	return d.Usb() != ""
}

func (d Device) State() (DeviceState, error) {
	resp, err := d.adbClient.executeCommand(fmt.Sprintf("host-serial:%s:get-state", d.serial))
	return deviceStateConv(resp), err
}

func (d Device) DevicePath() (string, error) {
	resp, err := d.adbClient.executeCommand(fmt.Sprintf("host-serial:%s:get-devpath", d.serial))
	return resp, err
}

func (d Device) ForwardLocalAbstract(localPort int, remotePort string, noRebind ...bool) (err error) {
	local := fmt.Sprintf("tcp:%d", localPort)
	remote := fmt.Sprintf("localabstract:%s", remotePort)
	return d.forward(local, remote, noRebind...)
}

func (d Device) FrowardTcp(localPort int, remotePort string, noRebind ...bool) (err error) {
	local := fmt.Sprintf("tcp:%d", localPort)
	remote := fmt.Sprintf("tcp:%s", remotePort)
	return d.forward(local, remote, noRebind...)
}

func (d Device) forward(local, remote string, noRebind ...bool) (err error) {
	command := ""
	if len(noRebind) != 0 && noRebind[0] {
		command = fmt.Sprintf("host-serial:%s:forward:norebind:%s;%s", d.serial, local, remote)
	} else {
		command = fmt.Sprintf("host-serial:%s:forward:%s;%s", d.serial, local, remote)
	}

	_, err = d.adbClient.executeCommand(command, true)
	return
}

func (d Device) ForwardList() (deviceForwardList []DeviceForward, err error) {
	var forwardList []DeviceForward
	if forwardList, err = d.adbClient.ForwardList(); err != nil {
		return nil, err
	}

	deviceForwardList = make([]DeviceForward, 0, len(deviceForwardList))
	for i := range forwardList {
		if forwardList[i].Serial == d.serial {
			deviceForwardList = append(deviceForwardList, forwardList[i])
		}
	}
	// resp, err := d.adbClient.executeCommand(fmt.Sprintf("host-serial:%s:list-forward", d.serial))
	return
}

func (d Device) ForwardKill(localPort int) (err error) {
	local := fmt.Sprintf("tcp:%d", localPort)
	_, err = d.adbClient.executeCommand(fmt.Sprintf("host-serial:%s:killforward:%s", d.serial, local), true)
	return
}

func (d Device) ReverseLocalAbstract(remotePort string, localPort int, noRebind ...bool) (err error) {
	local := fmt.Sprintf("tcp:%d", localPort)
	remote := fmt.Sprintf("localabstract:%s", remotePort)
	return d.reverse(remote, local, noRebind...)
}

func (d Device) ReverseTcp(remotePort, localPort int, noRebind ...bool) (err error) {
	local := fmt.Sprintf("tcp:%d", localPort)
	remote := fmt.Sprintf("tcp:%d", remotePort)
	return d.reverse(remote, local, noRebind...)
}

func (d Device) reverse(remote, local string, noRebind ...bool) (err error) {
	//_, err = d.adbClient.executeCommand("host:transport:"+d.serial, true)
	command := ""
	if len(noRebind) != 0 && noRebind[0] {
		command = fmt.Sprintf("reverse:forward:norebind:%s;%s", remote, local)
	} else {
		command = fmt.Sprintf("reverse:forward:%s;%s", remote, local)
	}

	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return err
	}
	defer func() { _ = tp.Close() }()

	if err = tp.Send(command); err != nil {
		return err
	}

	if err = tp.VerifyResponse(); err != nil {
		return err
	}

	return
}

func (d Device) ReverseList() (deviceForward []DeviceForward, err error) {
	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return nil, err
	}
	defer func() { _ = tp.Close() }()

	err = tp.Send(fmt.Sprintf("reverse:list-forward"))
	if err != nil {
		return nil, err
	}
	data, err := tp.ReadStringAll()
	if err != nil {
		return nil, err
	}
	lines := strings.Split(data, "\n")
	deviceForward = make([]DeviceForward, 0, len(lines))

	for i := range lines {
		line := strings.TrimSpace(lines[i])
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) > 2 {
			deviceForward = append(deviceForward, DeviceForward{Serial: fields[0], Remote: fields[1], Local: fields[2]})
		}

	}

	return
}

func (d Device) ReverseKillLocalAbstract(remotePort string) (err error) {
	local := fmt.Sprintf("localabstract:%s", remotePort)
	return d.reverseKill(local)
}

func (d Device) ReverseKillTcp(localPort int) (err error) {
	local := fmt.Sprintf("tcp:%d", localPort)
	return d.reverseKill(local)
}

func (d Device) reverseKill(remote string) (err error) {
	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return err
	}
	defer func() { _ = tp.Close() }()

	err = tp.Send(fmt.Sprintf("reverse:killforward:%s", remote))
	if err != nil {
		return err
	}
	return
}

func (d Device) ReverseKillAll() (err error) {
	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return err
	}
	defer func() { _ = tp.Close() }()

	err = tp.Send(fmt.Sprintf("reverse:killforward-all"))
	if err != nil {
		return err
	}
	return
}

func (d Device) RunShellCommand(cmd string, args ...string) (string, error) {
	raw, err := d.RunShellCommandWithBytes(cmd, args...)
	return string(raw), err
}

func (d Device) RunShellCommandWithBytes(cmd string, args ...string) ([]byte, error) {
	if len(args) > 0 {
		cmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	}
	if strings.TrimSpace(cmd) == "" {
		return nil, errors.New("adb shell: command cannot be empty")
	}
	raw, err := d.executeCommand(fmt.Sprintf("shell:%s", cmd))
	return raw, err
}

func (d Device) RunShellLoopCommand(cmd string, args ...string) (io.Reader, error) {
	//if len(args) > 0 {
	//	cmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	//}
	//if strings.TrimSpace(cmd) == "" {
	//	return nil, errors.New("adb shell: command cannot be empty")
	//}
	//return d.executeLoopCommand(fmt.Sprintf("shell:%s", cmd))
	return d.RunShellLoopCommandSock(cmd, args...)
}
func (d Device) RunShellLoopCommandSock(cmd string, args ...string) (net.Conn, error) {
	if len(args) > 0 {
		cmd = fmt.Sprintf("%s %s", cmd, strings.Join(args, " "))
	}
	if strings.TrimSpace(cmd) == "" {
		return nil, errors.New("adb shell: command cannot be empty")
	}
	return d.executeLoopCommand(fmt.Sprintf("shell:%s", cmd))
}

func (d Device) EnableAdbOverTCP(port ...int) (err error) {
	if len(port) == 0 {
		port = []int{AdbDaemonPort}
	}

	_, err = d.executeCommand(fmt.Sprintf("tcpip:%d", port[0]), true)
	return
}

func (d Device) createDeviceTransport() (tp transport, err error) {
	if tp, err = newTransport(fmt.Sprintf("%s:%d", d.adbClient.host, d.adbClient.port)); err != nil {
		return transport{}, err
	}

	if err = tp.Send(fmt.Sprintf("host:transport:%s", d.serial)); err != nil {
		return transport{}, err
	}
	err = tp.VerifyResponse()
	return
}

//
//func (d Device) createDeviceTransportTimeout(readTimeout time.Duration) (tp transport, err error) {
//	if tp, err = newTransport(fmt.Sprintf("%s:%d", d.adbClient.host, d.adbClient.port), readTimeout); err != nil {
//		return transport{}, err
//	}
//
//	if err = tp.Send(fmt.Sprintf("host:transport:%s", d.serial)); err != nil {
//		return transport{}, err
//	}
//	err = tp.VerifyResponse()
//	return
//}

func (d Device) executeCommand(command string, onlyVerifyResponse ...bool) (raw []byte, err error) {
	if len(onlyVerifyResponse) == 0 {
		onlyVerifyResponse = []bool{false}
	}

	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return nil, err
	}
	defer func() { _ = tp.Close() }()

	if err = tp.Send(command); err != nil {
		return nil, err
	}

	if err = tp.VerifyResponse(); err != nil {
		return nil, err
	}

	if onlyVerifyResponse[0] {
		return
	}

	raw, err = tp.ReadBytesAll()
	return
}

func (d Device) executeLoopCommand(command string, onlyVerifyResponse ...bool) (raw net.Conn, err error) {
	if len(onlyVerifyResponse) == 0 {
		onlyVerifyResponse = []bool{false}
	}

	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return nil, err
	}
	//defer func() { _ = tp.Close() }()

	if err = tp.Send(command); err != nil {
		return nil, err
	}

	if err = tp.VerifyResponse(); err != nil {
		return nil, err
	}

	if onlyVerifyResponse[0] {
		return
	}

	return tp.sock, err
}

//func (d Device) executeLoopCommandSock(command string, readTimeout time.Duration, onlyVerifyResponse ...bool) (raw net.Conn, err error) {
//	if len(onlyVerifyResponse) == 0 {
//		onlyVerifyResponse = []bool{false}
//	}
//
//	var tp transport
//	if tp, err = d.createDeviceTransportTimeout(readTimeout); err != nil {
//		return nil, err
//	}
//	//defer func() { _ = tp.Close() }()
//
//	if err = tp.Send(command); err != nil {
//		return nil, err
//	}
//
//	if err = tp.VerifyResponse(); err != nil {
//		return nil, err
//	}
//
//	if onlyVerifyResponse[0] {
//		return
//	}
//
//	return tp.sock, err
//}

func (d Device) List(remotePath string) (devFileInfos []DeviceFileInfo, err error) {
	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return nil, err
	}
	defer func() { _ = tp.Close() }()

	var sync syncTransport
	if sync, err = tp.CreateSyncTransport(); err != nil {
		return nil, err
	}
	defer func() { _ = sync.Close() }()

	if err = sync.Send("LIST", remotePath); err != nil {
		return nil, err
	}

	devFileInfos = make([]DeviceFileInfo, 0)

	var entry DeviceFileInfo
	for entry, err = sync.ReadDirectoryEntry(); err == nil; entry, err = sync.ReadDirectoryEntry() {
		if entry == (DeviceFileInfo{}) {
			break
		}
		devFileInfos = append(devFileInfos, entry)
	}

	return
}

func (d Device) PushFile(local *os.File, remotePath string, modification ...time.Time) (err error) {
	if len(modification) == 0 {
		var stat os.FileInfo
		if stat, err = local.Stat(); err != nil {
			return err
		}
		modification = []time.Time{stat.ModTime()}
	}

	return d.Push(local, remotePath, modification[0], DefaultFileMode)
}

func (d Device) Push(source io.Reader, remotePath string, modification time.Time, mode ...os.FileMode) (err error) {
	if len(mode) == 0 {
		mode = []os.FileMode{DefaultFileMode}
	}

	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return err
	}
	defer func() { _ = tp.Close() }()

	var sync syncTransport
	if sync, err = tp.CreateSyncTransport(); err != nil {
		return err
	}
	defer func() { _ = sync.Close() }()

	data := fmt.Sprintf("%s,%d", remotePath, mode[0])
	if err = sync.Send("SEND", data); err != nil {
		return err
	}

	if err = sync.SendStream(source); err != nil {
		return
	}

	if err = sync.SendStatus("DONE", uint32(modification.Unix())); err != nil {
		return
	}

	if err = sync.VerifyStatus(); err != nil {
		return
	}
	return
}

func (d Device) Pull(remotePath string, dest io.Writer) (err error) {
	var tp transport
	if tp, err = d.createDeviceTransport(); err != nil {
		return err
	}
	defer func() { _ = tp.Close() }()

	var sync syncTransport
	if sync, err = tp.CreateSyncTransport(); err != nil {
		return err
	}
	defer func() { _ = sync.Close() }()

	if err = sync.Send("RECV", remotePath); err != nil {
		return err
	}

	err = sync.WriteStream(dest)
	return
}
