package server_test

import (
	"encoding/json"
	"fionna/entity"
	"fionna/server"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"
)

func getGinEngine() *gin.Engine {
	r := gin.Default()
	r.Use(server.Cors())
	return r
}

// 将结构体转换为 URL 参数
func mapToURLValues(s map[string]string) url.Values {
	values := url.Values{}

	for k, v := range s {
		values.Set(k, v)
	}
	return values
}

func simpleTestHttpResults(t *testing.T, resStruct interface{}, rec *httptest.ResponseRecorder) {

	result := rec.Result()
	if result.StatusCode != 200 {
		t.Fatalf("请求状态码不符合预期")
	}

	body, err := ioutil.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("读取返回内容失败, err:%v", err)
	}
	defer result.Body.Close()

	err = json.Unmarshal(body, resStruct)
	if err != nil {
		t.Fatalf("转换结果失败,err: %v", err)
	}
	if response, ok := resStruct.(entity.ResponseData); ok && response.Code != entity.RequestSucceed {
		t.Fatalf("请求失败了, 相应结果:%v", string(body))
	}
	t.Log("res:", string(body))
	t.Log("用例测试通过")
}

func TestAndroidSerialListUrl(t *testing.T) {
	r := getGinEngine()
	server.GroupAndroidSerialUrl(r)
	req, err := http.NewRequest(http.MethodGet, "/android/serial/list", nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var listRes = &entity.ResponseData{
		Data: []entity.SerialInfo{},
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, listRes, rec)
}

func TestAndroidDefaultSerialUrl(t *testing.T) {
	r := getGinEngine()
	server.GroupAndroidSerialUrl(r)
	req, err := http.NewRequest(http.MethodGet, "/android/serial/default", nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: "", // serial value
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestAndroidSerialInfoUrl(t *testing.T) {
	r := getGinEngine()
	server.GroupAndroidSerialUrl(r)
	req, err := http.NewRequest(http.MethodGet, "/android/serial/info?udid="+"91cf5f1c", nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var listRes = &entity.ResponseData{
		Data: []string{},
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, listRes, rec)
}

func TestAndroidPackageNameList(t *testing.T) {
	r := getGinEngine()
	server.GroupAndroidPackageUrl(r)

	params := map[string]string{
		"udid": "91cf5f1c",
	}

	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/android/app/list"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: []string{}, // package name list value
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestAndroidCurrentPackageName(t *testing.T) {
	r := getGinEngine()
	server.GroupAndroidPackageUrl(r)

	params := map[string]string{
		"serial": "192.168.2.198:5555",
	}

	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/android/app/current"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: "", // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportList(t *testing.T) {
	r := getGinEngine()
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"name": "OnePlus8T_CH_KB2000_com.baidu.tieba_mini_2024-04-24 02:36:30",
		"page": "1",
		"size": "10",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/list"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: map[string]interface{}{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetConfig(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "b022bc54-dc62-41d1-80c5-57e71cf91be9",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/config"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: entity.PerfConfig{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetSummary(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "6dd166a2-f6b2-426f-960d-45393a0aaed6",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/summary"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: entity.OverallSummary{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetProcCpuData(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "6dd166a2-f6b2-426f-960d-45393a0aaed6",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/proc/cpu"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: []entity.ProcCpuInfo{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetProcMemData(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "6dd166a2-f6b2-426f-960d-45393a0aaed6",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/proc/mem"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: []entity.ProcMemInfo{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetThreadData(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "6dd166a2-f6b2-426f-960d-45393a0aaed6",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/proc/thread"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: []entity.ProcThreadsInfo{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetSysNetworkData(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "6dd166a2-f6b2-426f-960d-45393a0aaed6",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/sys/network"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: []entity.SystemNetworkInfo{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetSysMemData(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "6dd166a2-f6b2-426f-960d-45393a0aaed6",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/sys/mem"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: []entity.SystemMemInfo{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetSysCpuData(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "6dd166a2-f6b2-426f-960d-45393a0aaed6",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/sys/cpu"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: []entity.SystemCPUInfo{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

func TestReportGetSysFrameData(t *testing.T) {
	r := getGinEngine()
	server.SetSQLDebug(true)
	server.InitDB()
	server.GroupReportUrl(r)

	params := map[string]string{
		"uuid": "6dd166a2-f6b2-426f-960d-45393a0aaed6",
	}
	encode := mapToURLValues(params).Encode()

	req, err := http.NewRequest(http.MethodGet, "/report/sys/frame"+"?"+encode, nil)
	if err != nil {
		t.Fatalf("构建请求失败, err: %v", err)
	}

	rec := httptest.NewRecorder()

	r.ServeHTTP(rec, req)
	var defaultSerialNameRes = &entity.ResponseData{
		Data: []entity.SysFrameInfo{}, // current package name
		Code: entity.RequestSucceed,
	}
	simpleTestHttpResults(t, defaultSerialNameRes, rec)
}

//func TestAndroidFramePerfWebSocket(t *testing.T) {
//	// 这个测试要先运行main
//	header := make(http.Header)
//	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/android/perf/sys/frame"}
//	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
//	if err != nil {
//		t.Fatalf("连接 WebSocket 服务器失败, err: %v", err)
//		return
//	}
//	defer conn.Close()
//
//	// 启动一个 goroutine 用于接收 WebSocket 服务器的响应
//	go func() {
//		perfConfig := &entity.PerfConfig{
//			UUID: "cae028b3-52ad-43d6-bcd0-129ae9f23669", // 这个UUID需要存在于数据库中，如果没有，则需要先调用/android/perf/start接口
//		}
//
//		err = conn.WriteJSON(perfConfig)
//		if err != nil {
//			t.Fatalf("发送数据失败, err: %v", err)
//		}
//		for {
//			_, message, err := conn.ReadMessage()
//			if err != nil {
//				log.Printf("server> ERROR! %v\n", err)
//				return
//			}
//			fmt.Println(string(message))
//		}
//	}()
//	time.Sleep(50 * time.Second)
//}

func TestAndroidScrcpyWebSocket(t *testing.T) {
	// 这个测试要先运行main
	header := make(http.Header)
	u := url.URL{Scheme: "ws", Host: "127.0.0.1:8080", Path: "/android/scrcpy"}
	conn, _, err := websocket.DefaultDialer.Dial(u.String(), header)
	if err != nil {
		t.Fatalf("连接 WebSocket 服务器失败, err: %v", err)
		return
	}
	defer conn.Close()

	ticker := time.NewTicker(30 * time.Second)
	ticker1 := time.NewTicker(1 * time.Second)

	// 启动一个 goroutine 用于接收 WebSocket 服务器的响应
	go func() {
		devConfig := &entity.ScrcpyDevice{
			UDID: "91cf5f1c",
		}

		err = conn.WriteJSON(devConfig)
		if err != nil {
			t.Fatalf("发送数据失败, err: %v", err)
		}

		go func() {
			for {
				<-ticker1.C
				conn.WriteJSON(&entity.ScrcpyRecvMessage{
					MessageType: entity.ScrcpyPongType,
					Data:        "",
				})
			}
		}()

		for {
			select {
			case <-ticker.C:
				break
			case <-ticker1.C:
				conn.WriteJSON(entity.ScrcpyRecvMessage{
					MessageType: entity.ScrcpyPongType,
				})
			default:
				dataType, message, err := conn.ReadMessage()
				if err != nil {
					log.Printf("server> ERROR! %v\n", err)
					return
				}
				if dataType == websocket.BinaryMessage {
					//log.Println("the is binary message!")
				} else {
					fmt.Println(string(message))
				}
			}
		}
	}()
	time.Sleep(50 * time.Second)
}
