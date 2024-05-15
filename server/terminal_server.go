package server

import (
	"context"
	"encoding/json"
	"fionna/android/android_util"
	"fionna/entity"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
	"strings"
	"sync"
)

var storeTerminalLock sync.Mutex
var storeTerminalDevice = make(map[string]net.Conn)

func getTerminalIo(uuid string) net.Conn {
	var commandReader net.Conn
	storeTerminalLock.Lock()
	commandReader = storeTerminalDevice[uuid]
	storeTerminalLock.Unlock()
	return commandReader
}

func putTerminalIo(uuid string, commandReader net.Conn) {
	storeTerminalLock.Lock()
	storeTerminalDevice[uuid] = commandReader
	storeTerminalLock.Unlock()
}

func deleteTerminalIo(uuid string) {
	storeTerminalLock.Lock()
	delete(storeTerminalDevice, uuid)
	storeTerminalLock.Unlock()
}

func WebSocketTerminal(r *gin.Engine) {
	r.GET("/android/terminal", func(c *gin.Context) {

		serial := c.Query("udid")
		if serial == "" {
			log.Error("serial is empty")
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: "serial is empty",
				Code: entity.ParameterErr,
			})
			return
		}
		device, err := android_util.GetDevice(client, serial)
		if err != nil {
			c.JSON(http.StatusOK, entity.ResponseData{
				Data: entity.CodeDefaultMessage[entity.GetDeviceErr],
				Code: entity.GetDeviceErr,
			})
			log.Error(err)
		}

		terminalWs, err := upGrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			log.Print("Error during connection upgradation:", err)
			return
		}

		ws := NewSafeWebsocket(terminalWs)

		exitCtx, exitFn := context.WithCancel(context.Background())

		go func() {
			for {
				select {
				case <-exitCtx.Done():
					return
				default:
					defer func() {
						if r := recover(); r != nil {
							log.Error("terminal ws recovered:", r)
						}
					}()
					var message entity.TerminalRecvMessage

					err := ws.ReadJSON(&message)
					if err != nil {
						log.Error("terminal read message steam err:", err)
						ws.WriteJSON(entity.NewTerminalErrorInfo("terminal read message steam err:" + err.Error()))
						break
					} else {
						if message.MessageType == entity.CommandTerminalType {
							if command, ok := message.Data.(string); ok {
								if strings.Contains(command, "reboot") || strings.Contains(command, "rm") || strings.Contains(command, "su ") {
									ws.WriteJSON(&entity.TerminalSendMessage{
										Uuid:        message.Uuid,
										MessageType: entity.SendCommandRespEndType,
									})
									return
								}

								if message.Uuid == "" {
									uuidObj, _ := uuid.NewUUID()
									message.Uuid = uuidObj.String()
									output, err := device.RunShellLoopCommandSock(command)
									if err != nil {
										log.Error("terminal execute command err:", err)
										ws.WriteJSON(entity.NewTerminalErrorInfo("terminal execute command err:" + err.Error()))

										return
									}

									putTerminalIo(message.Uuid, output)

								}

								output := getTerminalIo(message.Uuid)

								if output == nil {
									log.Error("terminal execute command not get ")
									ws.WriteJSON(entity.NewTerminalErrorInfo("terminal execute command not get"))

									return
								}

								go func() {
									for {
										var bytesOutput = make([]byte, 1024)
										n, err1 := output.Read(bytesOutput)
										if err1 != nil {
											ws.WriteJSON(&entity.TerminalSendMessage{
												Uuid:        message.Uuid,
												MessageType: entity.SendCommandRespEndType,
											})
											log.Error(err1)
											return
										}
										err1 = ws.WriteJSON(&entity.TerminalSendMessage{
											Uuid:        message.Uuid,
											MessageType: entity.SendCommandRespType,
											Data:        string(bytesOutput[:n]),
										})
										if err1 != nil {
											panic(err1)
										}
									}

								}()

							}
						}
						if message.MessageType == entity.LogcatTerminalType {

							data, err1 := json.Marshal(message.Data)
							if err1 != nil {
								log.Error("terminal logcat the data sent is not json")
								ws.WriteJSON(entity.NewTerminalErrorInfo("terminal logcat the data sent is not json"))

								break
							}
							// todo uuid
							var logcatCommandConfig = &entity.LogcatRecvMessage{}
							err1 = json.Unmarshal(data, logcatCommandConfig)

							if err1 == nil {
								logcatMessage := "logcat *:" + logcatCommandConfig.Level
								if logcatCommandConfig.Filter != "" {
									logcatMessage = logcatMessage + " | grep " + logcatCommandConfig.Filter
								}

								if message.Uuid == "" {
									uuidObj, _ := uuid.NewUUID()
									message.Uuid = uuidObj.String()
									output, err := device.RunShellLoopCommandSock(logcatMessage)
									if err != nil {
										log.Error("logcat execute command err:", err)
										ws.WriteJSON(entity.NewTerminalErrorInfo("logcat execute command err:" + err.Error()))
										return
									}

									putTerminalIo(message.Uuid, output)

								}

								output := getTerminalIo(message.Uuid)

								if output == nil {
									log.Error("terminal logcat execute command not get ")
									ws.WriteJSON(entity.NewTerminalErrorInfo("terminal logcat execute command not get "))
									return
								}

								go func() {
									for {
										var bytesOutput = make([]byte, 1024)
										n, err1 := output.Read(bytesOutput)
										if err1 != nil {
											ws.WriteJSON(&entity.TerminalSendMessage{
												Uuid:        message.Uuid,
												MessageType: entity.SendLogcatRespEndType,
											})
											log.Error(err1)
											return
										}
										err1 = ws.WriteJSON(&entity.TerminalSendMessage{
											Uuid:        message.Uuid,
											MessageType: entity.SendLogcatRespType,
											Data:        string(bytesOutput[:n]),
										})
										if err1 != nil {
											panic(err1)
										}
									}

								}()

							} else {
								log.Error("logcat conversion message error,", err1)
								ws.WriteJSON(entity.NewTerminalErrorInfo("logcat conversion message error:" + err1.Error()))
								break
							}
						}
						if message.MessageType == entity.StopCommandType {
							log.Info("close command")
							output := getTerminalIo(message.Uuid)
							output.Close()
							deleteTerminalIo(message.Uuid)
							//exitFn()
						}
						if message.MessageType == entity.StopLogcatType {
							log.Info("close logcat")
							output := getTerminalIo(message.Uuid)
							output.Close()
							deleteTerminalIo(message.Uuid)
							//exitFn()
						}
						if message.MessageType == entity.CloseTerminalType {
							// todo close all io
							exitFn()
						}
						if message.MessageType == entity.PongTerminalType {
							continue
						}
					}
				}
			}
		}()
	})
}
