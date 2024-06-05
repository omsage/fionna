# Fionna

[English](https://github.com/omsage/fionna/blob/master/README_EN.md )

## 概述

fionna是一个针对安卓端的性能采集工具，其功能非常简洁。命名灵感来自于《**冒险时光**》中的Fionna Campbell。



当前已实现的主要功能有：

- 解决安卓SurfaceFlinger、GFXInfo这两个方案在游戏等方向上无法采集Frame数据的痛点，具体请参考frame benchmark
- 其他性能数据，如CPU、Memory、温度、网络......等数据的采集
- 设备投屏和触控
- 测试报告
- 数据导出
- 数据对比

**Home**

![](./doc/Home.png)

**perf**

![](doc/Perf.png)

**Terminal**

![](./doc/Terminal.png)

**report**

![](./doc/Report.png)

## Frame Benchmark

具体的代码在test/frame_benchmark_test.go中，我在里面定义了三个用于对比的方法：TestFPSBySurfaceFlinger、TestFPSByGFXInfo、TestFPSByFrameTool，分别对应SurfaceFlinger、GFXInfo、新方案工具的获取方式。

测试使用机型为iqoo11s，软件为《白尘禁区》（com.dragonli.projectsnow.lhm），三个方法的获取结果为：

**TestFPSBySurfaceFlinger**

![](./doc/TestFPSBySurfaceFlinger.png)

**TestFPSByGFXInfo**

![](./doc/TestFPSByGFXInfo.png)

**OMSage Frame Tool**

![](./doc/TestFPSByOMSageFrameTool.png)

可以看到，使用OMSage Frame Tool工具获取Frame的性能数据有更好的兼容性，基本上能适用安卓大部分场景下的Frame性能采集场景。

**OMSage Frame Tool的使用请参考具体代码，并严格按照其代码逻辑使用，二开过程中出现的任何问题，不予解答。**

## 使用

### 快速使用

通过Release下载对应的构建产物，解压后直接执行对应的目标程序即可

命令行中直接执行：

```
fionna
```

之后浏览器访问对应的地址： http://127.0.0.1:3417 

### 命令行模式

fionna内置简单的命令行模式，可以输入以下命令查看说明：

```
fionna --help
```

主要有两个模式：

#### web模式

在命令行中直接执行以下命令，程序会直接开启一个web服务，之后可以访问对应的地址：http://127.0.0.1:3417 

```
fionna web
```

当然，在web模式下，也可以指定使用的SQLite DB文件：

```
fionna web -d <db path>
# 参考示例:fionna fionna web -d ./exapmle/my.db
```

#### cli-perf模式

由于web的依赖性太强，所以提供一个简单的性能输出模式，输入以下命令，即可在命令行中得到对应的性能数据：

```
fionna cli-perf [flags]
#参考示例:fionna cli-perf --fps
```

**可选参数**

| 快捷使用 | 选项名         | 参数类型 | 描述信息                     |
| -------- | -------------- | -------- | ---------------------------- |
| -h       | --help         |          | 获取帮助指南                 |
| -p       | --package      | string   | 应用包名                     |
| -d       | --pid          | int      | 应用PID (默认 -1)            |
|          | --proc-cpu     |          | 获取进程cpu数据              |
|          | --fps          |          | 获取系统FPS                  |
|          | --jank         |          | 获取系统jank信息             |
|          | --proc-mem     |          | 获取进程内存数据             |
|          | --proc-threads |          | 获取进程线程数量             |
| -s       | --serial       | string   | 设备序列号（默认第一个设备） |
|          | --sys-cpu      |          | 获取系统cpu数据              |
|          | --sys-mem      |          | 获取系统内存数据             |
|          | --sys-network  |          | 获取系统网络数据             |

## 开发

**二开本项目请遵守AGPL协议！！！！！！**

**二开本项目请遵守AGPL协议！！！！！！**

**二开本项目请遵守AGPL协议！！！！！！**



**本项目嵌入了一个简单的vue前端，所以开发前必须先进行前端的设置。如果需要自定义开发，请先按照以下步骤操作**

- 拉取项目后，首先进入到fionna-web目录

```
cd fionna-web
```

- 拉取前端依赖

```
npm install
```

- 生成构建产物

```
npm run build
```

接下来即可正常根据需求开发项目

## 注意事项

- PC上必须有adb环境
- 应用的CPU计算和常规的计算有所不同，项目使用的指标为：https://blog.csdn.net/weixin_39451323/article/details/118083713
- 避免占用3417端口 

## 感谢

- https://github.com/SonicCloudOrg/sonic-android-supply
- https://github.com/SonicCloudOrg/sonic-client-web
- https://github.com/electricbubble/gadb
- https://github.com/Genymobile/scrcpy