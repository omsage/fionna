package server

import (
	"fionna/entity"
	"fmt"
	"github.com/xuri/excelize/v2"
	"reflect"
	"strings"
)

func writeSheet(xlsx *excelize.File, sheetName string, dataValue interface{}) {
	index, _ := xlsx.NewSheet(sheetName)

	switch reflect.TypeOf(dataValue).Kind() {
	case reflect.Slice, reflect.Array:
		data := reflect.ValueOf(dataValue)
		for i := 0; i < data.Len(); i++ {
			dValue := data.Index(i)
			d := dValue.Type()
			for j := 0; j < d.NumField(); j++ {
				if i == 0 {
					columns := d.Field(j).Tag.Get("xlsx")
					if columns == "" {
						continue
					}
					//根据结构体中绑定的tag，根据分隔符，拿到列号
					column := strings.Split(columns, "-")[0]
					//同理拿到列名
					name := strings.Split(columns, "-")[1]
					// 设置表头
					xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", column, i+1), name)

				}
				//	设置内容
				column := strings.Split(d.Field(j).Tag.Get("xlsx"), "-")[0]

				xlsx.SetCellValue(sheetName, fmt.Sprintf("%s%d", column, i+2), dValue.Field(j))
			}

		}
	}

	xlsx.SetActiveSheet(index)
}

func Export2Excel(uuid string) *excelize.File {
	xlsx := excelize.NewFile()
	var perfConfig entity.PerfConfig
	db.First(&perfConfig, "uuid = ?", uuid)
	// todo dao?this project is not so strict

	if perfConfig.SysCpu {
		var sysCpuDatas []entity.SystemCPUInfo
		db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&sysCpuDatas)
		writeSheet(xlsx, "sys-cpu", sysCpuDatas)

		var cpuSummarys []entity.SystemCPUSummary
		db.Where("uuid = ?", uuid).Find(&cpuSummarys)

		writeSheet(xlsx, "sys-cpu"+"-"+"summary", cpuSummarys)

	}

	if perfConfig.SysMem {
		var SYSMemDatas []entity.SystemMemInfo
		db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&SYSMemDatas)
		writeSheet(xlsx, "sys-mem", SYSMemDatas)

		var sysMemSummarys []entity.SystemMemSummary
		db.Where("uuid = ?", uuid).Find(&sysMemSummarys)
		writeSheet(xlsx, "sys-mem"+"-"+"summary", sysMemSummarys)
	}

	if perfConfig.SysTemperature {
		var sysTemperatureDatas []entity.SysTemperature
		db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&sysTemperatureDatas)
		writeSheet(xlsx, "sys temperature", sysTemperatureDatas)

		var sysTemperature []entity.SystemTemperatureSummary
		db.Where("uuid = ?", uuid).Find(&sysTemperature)
		writeSheet(xlsx, "sys-temperature"+"-"+"summary", sysTemperature)
	}

	if perfConfig.SysNetwork {
		var SysNetworkDatas []entity.SystemNetworkInfo
		db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&SysNetworkDatas)
		writeSheet(xlsx, "sys-network", SysNetworkDatas)

		var netSummarys []entity.SystemNetworkSummary
		db.Where("uuid = ?", uuid).Find(&netSummarys)
		writeSheet(xlsx, "sys-network"+"-"+"summary", netSummarys)
	}

	if perfConfig.FPS || perfConfig.Jank {
		var sysFrameDatas []entity.SysFrameInfo
		db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&sysFrameDatas)
		writeSheet(xlsx, "frame", sysFrameDatas)

		var frameSummary []entity.FrameSummary
		db.Where("uuid = ?", uuid).Find(&frameSummary)
		writeSheet(xlsx, "sys-frame"+"-"+"summary", frameSummary)
	}

	if perfConfig.ProcThread {
		var procThreadDatas []entity.ProcThreadsInfo
		db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&procThreadDatas)
		writeSheet(xlsx, "proc-thread", procThreadDatas)
	}

	if perfConfig.ProcCpu {
		var procCpuDatas []entity.ProcCpuInfo
		db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&procCpuDatas)
		writeSheet(xlsx, "proc-cpu", procCpuDatas)

		var procCpuSummary []entity.ProcCpuSummary
		db.Where("uuid = ?", uuid).Find(&procCpuSummary)
		writeSheet(xlsx, "proc-cpu"+"-"+"summary", procCpuSummary)
	}

	if perfConfig.ProcMem {
		var procMemDatas []entity.ProcMemInfo
		db.Order("timestamp asc").Where("uuid = ?", uuid).Find(&procMemDatas)
		writeSheet(xlsx, "proc-mem", procMemDatas)

		var procMemSummary []entity.ProcMemSummary
		db.Where("uuid = ?", uuid).Find(&procMemSummary)
		writeSheet(xlsx, "proc-mem"+"-"+"summary", procMemSummary)
	}

	//err := xlsx.SaveAs("test_write.xlsx")
	//if err != nil {
	//	fmt.Println(err)
	//}
	return xlsx
}
