# 一.说明
最新版本请从 (https://github.com/ci-plugins/DockerBuildPush) 获取

本软件是免费开源软件 [bk-ci](https://github.com/Tencent/bk-ci) 
的第三方插件,主要用于编译,推送,复制OCI容器镜像三项功能.
与本软件类似或同类的还有 [drone-kaniko](https://github.com/drone/drone-kaniko) 和  
[action-kaniko](https://github.com/aevea/action-kaniko) 等

本软件主程序是一个命令行包装器,主要功能就是将蓝盾的外部输入转化为第三方外部依赖软件的输入.
其核心编译,推送,复制功能并不由本软件直接提供.

其中编译,推送功能由 [kaniko](https://github.com/GoogleContainerTools/kaniko) 提供,
复制功能由 [skopeo](https://github.com/containers/skopeo) 提供.

本软件通过`os.exec`或者生成`myrun.sh`脚本来进行第三方外部依赖软件的调用实现最终编译,推送功能.

# 二.安装方法
## 1.获取源代码
源代码请固定存放到`/data/landun/workspace/src/DockerBuildPushGo`
否则请修改Makefile第一行的`CURRENT_DIR`变量定义为实际路径
```
mkdir -p /data/landun/workspace/src/DockerBuildPushGo
cd /data/landun/workspace/src/DockerBuildPushGo
git clone https://github.com/ci-plugins/DockerBuildPush
```

## 2.下载proot执行文件
```
curl -L -o ./bin_file/proot https://github.com/proot-me/proot/releases/download/v5.3.0/proot-v5.3.0-x86_64-static
chmod +x ./bin_file/proot
```
成功的话将在`./bin_file/`目录生成`proot`执行文件

## 3.构建DockerBuildPush主程序
执行2之前先确保`./bin_file/proot`执行文件已经存在了
```
#./是项目源代码根目录
cd ./
make -f Makefile
```
成功的话将在当前目录的上两级目录生成`app`(无扩展名)和`task.json`文件.
`../../bin/app`和`../../bin/task.json`

## 4.打包插件
请把app和task.json打包到zip格式的根目录中.不要放在任何二级目录下.
```
cd ../../bin
zip -r ./app.zip app task.json
```

# 三.依赖情况
# 1.编译时第三方依赖
[golang](https://golang.google.cn/) v1.14+

[go.rice](https://github.com/GeertJohan/go.rice) 固定v1.0.2

[glibc-static](https://www.gnu.org/software/libc/) v2.12+

[gcc](https://www.gnu.org/software/gcc/) v4.8+

[curl](https://github.com/curl/curl) v7.2+

# 2.运行时第三方依赖
[kaniko](https://github.com/GoogleContainerTools/kaniko) 固定v1.7.0

[skopeo](https://github.com/containers/skopeo) v1.5.2+

[proot](https://github.com/proot-me/proot) 固定v5.3.0
