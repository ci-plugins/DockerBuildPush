# 安装方法
## 1.静态链接生成proot执行文件
必须在非docker环境中进行本步操作.
```
cd ./bin_file/
chmod +x ./build_proot.sh
./build_proot.sh

```
成功的话将在`./bin_file/`目录生成`proot`执行文件

## 2.构建bkci-DockerBuildPush主程序
执行2之前先确保`./bin_file/proot`执行文件已经存在了
```
#./是项目源代码根目录
cd ./
make -f Makefile
```
成功的话将在当前目录的上两级目录生成`app`(无扩展名)和`task.json`文件.
`../../bin/app`和`../../bin/task.json`

## 3.打包插件
请把app和task.json打包到zip格式的根目录中.不要放在任何二级目录下.
```
cd ../../bin
zip -r ./app.zip app task.json
```