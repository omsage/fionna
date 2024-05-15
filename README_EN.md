# Fionna

## Overview

Fionna is a performance monitoring tool designed for Android devices, with a focus on simplicity. It draws inspiration from Fionna Campbell in "**Adventure Time**".

Key features currently implemented include:

- Addressing the limitations of Android's SurfaceFlinger and GFXInfo solutions in capturing frame data for games and similar applications; refer to the frame benchmark for details
- Monitoring other performance metrics such as CPU, memory, temperature, network, etc.
- Device screen mirroring and touch input
- Test reporting
- Data export
- Data comparison

**Home**

![](./doc/Home.png)

**perf**

![](doc/Perf.png)

**Terminal**

![](./doc/Terminal.png)

**report**

![](./doc/Report.png)

## Frame Benchmark

The specific code is in test/frame_benchmark_test.go, where I defined three methods for comparison: TestFPSBySurfaceFlinger, TestFPSByGFXInfo, and TestFPSByFrameTool, corresponding to the ways of obtaining information from SurfaceFlinger, GFXInfo, and the new tool solution, respectively.

The testing was conducted on the iqoo11s device, running the software "**Snowbreak: Containment Zone**" (com.dragonli.projectsnow.lhm). The results obtained from the three methods are as follows:

**TestFPSBySurfaceFlinger**

![](./doc/TestFPSBySurfaceFlinger.png)

**TestFPSByGFXInfo**

![](./doc/TestFPSByGFXInfo.png)

**OMSage Frame Tool**

![](./doc/TestFPSByOMSageFrameTool.png)

You can see that using the OMSage Frame Tool to obtain performance data for Frames has better compatibility, making it suitable for most Android scenarios involving Frame performance collection.

**Please refer to the specific code for the usage of the OMSage Frame Tool and strictly follow its code logic. Any issues encountered during customization will not be addressed.**

## Usage

Download the corresponding build artifact through Release, extract it, and then directly execute the target program.

Execute directly in the command line:

```
fionna
```

Then access the corresponding address in a browser: [http://127.0.0.1:3417](http://127.0.0.1:3417/)

## Development

**When customizing this project, please adhere to the AGPL license!!!!**

**This project embeds a simple Vue frontend, so frontend setup must be done before development. If you need to customize development, follow the steps below first:**

- After cloning the project, navigate to the fionna-web directory

```
cd fionna-web
```

- Install frontend dependencies

```
npm install
```

- Generate build artifacts

```
npm run build
```

You can now proceed with project development according to your requirements.

## Notes

- ADB environment is required on the PC
- Application's CPU computation differs from regular calculations; the project uses the following metrics:https://blog.csdn.net/weixin_39451323/article/details/118083713

- Avoid occupying port 3417.

## Thank

- https://github.com/SonicCloudOrg/sonic-android-supply
- https://github.com/SonicCloudOrg/sonic-client-web
- https://github.com/electricbubble/gadb
- https://github.com/Genymobile/scrcpy