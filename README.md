# OpenOEPGo
开源教育平台
采用CGo编写，依赖libav*,libx265
1.  Linux系统编译步骤，以Ubuntu为例:
apt update && apt upgrade -y && apt install libavdevice-dev libavcodec-dev libavformat-dev libavutil-dev libavresample-dev libswscale-dev libx265-dev
下载源码执行make编译
2.  Windows:
下载msys2并安装
替换为清华源https://mirrors.tuna.tsinghua.edu.cn/help/msys2/
打开msys64执行pacman -Syu && pacman -Sy && pacman -S mingw-w64/ffmpeg_x86_64
修改PATH把go可执行文件加入PATH环境变量，如export PATH=/c/go/bin:$PATH
执行make
