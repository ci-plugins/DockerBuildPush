/*
 * Tencent is pleased to support the open source community by making BK-CI 蓝鲸持续集成平台 available.
 *
 * Copyright (C) 2019 THL A29 Limited, a Tencent company.  All rights reserved.
 *
 * BK-CI 蓝鲸持续集成平台 is licensed under the MIT license.
 *
 * A copy of the MIT License is included in this file.
 *
 *
 * Terms of the MIT License:
 * ---------------------------------------------------
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this
 * software and associated documentation
 * files (the "Software"), to deal in the Software without restriction, including
 * without limitation the rights to use,copy,
 * modify, merge, publish, distribute, sublicense, and/or sell copies of the Software,
 * and to permit persons to whom the
 * Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies
 * or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT
 * LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN
 * NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY,
 * WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
 * SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 */

package main

import (
	"bytes"
	"encoding/base64"
	"errors"
	"github.com/ci-plugins/bkci-DockerBuildPush/api"
	"github.com/ci-plugins/bkci-DockerBuildPush/log"
	"github.com/syyongx/php2go"
	"github.com/tidwall/sjson"
	"github.com/txn2/txeh"
	"io"
	"net/url"
	"os"
	"os/exec"
	"strings"
)

func saveStringToFile(savePath string, jsonStr string) {

	err := php2go.FilePutContents(savePath, jsonStr, 0644)
	if err != nil {

	}

}

func addOrUpdateDockerConfigJson(savePath string, domain string, username string, password string) {
	if username != "" && domain != "" {

		encodeUsernameAndPassword := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
		//转义域名中的特殊字符避开jq语法关键字
		slashDomain := strings.Replace(domain, ".", "\\.", 100)

		//文件如果存在,先读取旧内容
		oldJsonStr, _ := php2go.FileGetContents(savePath)

		//有旧内容改成,无旧内容新增
		newJsonStr, _ := sjson.Set(oldJsonStr, "auths."+slashDomain+".auth", encodeUsernameAndPassword)
		saveStringToFile(savePath, newJsonStr)
	} else {
		log.Error("docker login域名和用户名不能为空")
	}
}

func createDockerConfigFile(savePath string) {
	targetImage := api.GetInputParam("targetImage")
	targetTicketId := api.GetInputParam("targetTicketId")
	username := ""
	password := ""
	username, password = getUserAndPass(targetTicketId, "dst")
	var repo *url.URL
	if !strings.Contains(targetImage, "http://") {
		repo, _ = url.Parse("http://" + targetImage)

	} else {
		repo, _ = url.Parse(targetImage)

	}

	addOrUpdateDockerConfigJson(savePath, repo.Host, username, password)
}

func myCopy(src, dst string) {

	_, err := php2go.Copy(src, dst)
	if err != nil {

	}
	log.Info("copy " + src + " to  " + dst)

}

func removeCRLF(input string) []string {
	inputNormalized := strings.Replace(input, "\r\n", "\n", -1)

	lines := strings.Split(inputNormalized, "\n")
	//dataLines := lines[:len(lines)-1]
	return lines
}

func fixedPath(input string) string {
	inputNormalized := ""
	inputNormalized = strings.Replace(input, "\r\n", "\n", -1)
	inputNormalized = strings.Replace(inputNormalized, "./", "/", -1)
	inputNormalized = strings.Replace(inputNormalized, "//", "/", -1)
	return inputNormalized
}

func createBashFile51(landunWorkSpacePath string) {

	targetImage := api.GetInputParam("targetImage")
	targetImageTag := api.GetInputParam("targetImageTag")
	ignorePath := api.GetInputParam("prootIgnorePath")
	dockerBuildDir := api.GetInputParam("dockerBuildDir")
	dockerFilePath := api.GetInputParam("dockerFilePath")
	dockerBuildArgs := api.GetInputParam("dockerBuildArgs")

	var LandunWorkSpacePath = landunWorkSpacePath
	var b strings.Builder

	if targetImageTag == "" {
		targetImageTag = "latest"
	}

	tagLines := removeCRLF(targetImageTag)
	var dstImagesTag string
	for _, v := range tagLines {
		if v != "" {
			dstImagesTag = dstImagesTag + " -d " + targetImage + ":" + v + " "
		}
	}
	if dstImagesTag == "" {
		dstImagesTag = " -d " + targetImage + ":latest "
	}

	argLines := removeCRLF(dockerBuildArgs)
	var dstArgs string
	for _, arg := range argLines {
		if strings.TrimSpace(arg) != "" {
			dstArgs = dstArgs + " --build-arg " + arg + " "
		}
	}

	//proot 5.1
	b.WriteString("export DOCKER_CONFIG=/kaniko/.docker/\n")

	b.WriteString("export DOCKER_CREDENTIAL_GCR_CONFIG=/kaniko/.config/gcloud/docker_credential_gcr_config.json\n")
	//proot执行参数
	//proot 5.1
	b.WriteString(LandunWorkSpacePath + "/proot -S " + LandunWorkSpacePath + "/kaniko_rootfs ")

	//kaniko执行参数
	//proot 5.1
	b.WriteString(" /kaniko/executor -v info ")

	//添加编译参数,内置环境变量
	b.WriteString(dstArgs)

	//在镜像打包做快照时要忽略的目录,这些目录的文件不会加进镜像里. etc下因为有/etc/hosts,/etc/resolv.conf整个目录被加进来.

	ignorePathLines := removeCRLF(ignorePath)
	var dstIgnores string
	for _, arg := range ignorePathLines {
		if strings.TrimSpace(arg) != "" {
			dstIgnores = dstIgnores + " --ignore-path=" + arg + " "
		}
	}

	b.WriteString(dstIgnores)
	b.WriteString(" --skip-tls-verify  ")
	//docker build的工作目录和Dockerfile路径
	fixedDockerfilePath := fixedPath("/workspace/" + dockerFilePath)
	fixedBuildWorkspacePath := fixedPath("/workspace/" + dockerBuildDir)

	b.WriteString(" -f  " + fixedDockerfilePath)
	b.WriteString(" -c " + fixedBuildWorkspacePath)
	//如果有多tag需要推送进行拼接
	b.WriteString(dstImagesTag)
	saveStringToFile(LandunWorkSpacePath+"/myrun.sh", b.String())
	err := os.Chmod(LandunWorkSpacePath+"/myrun.sh", 755)
	if err != nil {

	}
}

func addHostsToFile() {

	savePath := "/etc/hosts"

	hostContent := api.GetInputParam("dockerBuildHosts")
	if hostContent == "" {
		return
	}
	oldHostsFile, _ := txeh.ParseHosts(savePath)
	hosts, err := txeh.NewHostsDefault()
	if err != nil {
		panic(err)
	}
	//添加回文件已有的hosts内容.
	for _, v := range oldHostsFile {
		hosts.AddHosts(v.Address, v.Hostnames)
	}
	//添加插件上有的内容
	if hostContent != "" {
		//解析有bug.在尾部添加一个换行符,最后一个才能被读到
		addHost, _ := ParseStringHosts(hostContent + "\n")
		for _, l := range addHost {
			hosts.AddHosts(l.Address, l.Hostnames)
			log.Debug("add to stage : "+l.Address+" ,", l.Hostnames[0])
		}
	}

	log.Debug("==New hosts file content is==:\n" + hosts.RenderHostsFile() + "\n=================")
	errSave := hosts.SaveAs(savePath)
	if errSave != nil {

	}
}

func initDockerFileEnv(landunWorkSpacePath string, dockerBuildWorkpacePath string) {
	dockerFilePath := api.GetInputParam("dockerFilePath")
	myCopy(landunWorkSpacePath+"/"+dockerFilePath, dockerBuildWorkpacePath+"/Dockerfile")

	dockerBuildDir := api.GetInputParam("dockerBuildDir")

	//复制docker workspace工作目录进kaniko_rootfs workspace目录
	if dockerBuildDir == "" || dockerBuildDir == "/" || dockerBuildDir == "." || dockerBuildDir == "./" ||
		dockerBuildDir == landunWorkSpacePath {
		log.Error("docker build工作空间范围太大,有可能拷贝时间很长,或者镜像变太大.")
		log.Error("请缩小docker build目录范围,以及减少docker build目录文件数")
		os.Exit(1)
	} else {
		log.Info("正在复制docker workspace工作目录进kaniko rootfs workspace目录,请稍候,文件较多的话可能需等待较长时间.")
		workspaceSrcPath := fixedPath(landunWorkSpacePath + "/" + dockerBuildDir + "/.")
		workspaceDstPath := fixedPath(dockerBuildWorkpacePath + "/" + dockerBuildDir + "/")
		log.Info("copy " + workspaceSrcPath + " to  " + workspaceDstPath)

		err := exeCommandStdout("mkdir -p " + workspaceDstPath +
			" && cp -f  -r " + workspaceSrcPath + " " + workspaceDstPath)
		if err != nil {

		}
	}

}

func initKanikoEnv() (string, string, string, string) {

	log.Debug("开始构造kaniko和容器的私有工作目录")

	defaultWorkSpace := "/data/landun/workspace"

	if api.GetWorkspace() != "/" && api.GetWorkspace() != "" {
		defaultWorkSpace = api.GetWorkspace()
	}
	LandunWorkSpacePath := defaultWorkSpace
	KanikoRootFSPath := LandunWorkSpacePath + "/kaniko_rootfs"
	KanikoExecutePath := KanikoRootFSPath + "/kaniko"
	dockerBuildWorkpacePath := KanikoRootFSPath + "/workspace/"
	DockerConfigPath := KanikoExecutePath + "/.docker"
	err := errors.New("")
	err = os.MkdirAll(LandunWorkSpacePath, 755)
	err = os.MkdirAll(KanikoRootFSPath, 755)
	err = os.MkdirAll(KanikoExecutePath, 755)
	err = os.MkdirAll(DockerConfigPath, 755)
	err = os.MkdirAll(dockerBuildWorkpacePath, 755)
	if err != nil {

	}

	return LandunWorkSpacePath, KanikoExecutePath, dockerBuildWorkpacePath, DockerConfigPath
}

func exeCommandStdout(command string) error {
	var err error
	cmd := exec.Command("/bin/bash", "-c", command)

	var stdBuffer bytes.Buffer
	mw := io.MultiWriter(os.Stdout, &stdBuffer)

	cmd.Stdout = mw
	cmd.Stderr = mw

	err = cmd.Run()

	return err
}

func loginAndBuildPush() {

	//预创建各种所需要的工作目录
	LandunWorkSpacePath, KanikoExecutePath, dockerBuildWorkpacePath, DockerConfigPath := initKanikoEnv()
	//从可执行文件释放proot和kaniko到指定的工作目录
	unzipBinExeFileToDisk(LandunWorkSpacePath+"/proot", "proot")
	unzipBinExeFileToDisk(KanikoExecutePath+"/executor", "executor")
	//复制dockerfile和工作空间到工作目录
	initDockerFileEnv(LandunWorkSpacePath, dockerBuildWorkpacePath)
	//创建登录文件
	createDockerConfigFile(DockerConfigPath + "/config.json")
	//修改/etc/hosts
	addHostsToFile()
	//创建最终的bash运行脚本
	createBashFile51(LandunWorkSpacePath)
	//调用bash
	err := exeCommandStdout(LandunWorkSpacePath + "/myrun.sh")

	if err != nil {
		os.Exit(1)
	}
	log.Info("============")
	log.Info("编译镜像和push成功....")
	log.Info("============")

}

// ParseStringHosts Hosts文件读取
func ParseStringHosts(input string) ([]txeh.HostFileLine, error) {

	inputNormalized := strings.Replace(input, "\r\n", "\n", -1)

	lines := strings.Split(inputNormalized, "\n")
	dataLines := lines[:len(lines)-1]

	hostFileLines := make([]txeh.HostFileLine, len(dataLines))

	for i, l := range dataLines {
		curLine := &hostFileLines[i]
		curLine.OriginalLineNum = i
		curLine.Raw = l

		curLine.Trimed = strings.TrimSpace(l)

		if strings.HasPrefix(curLine.Trimed, "#") {
			curLine.LineType = txeh.COMMENT
			continue
		}

		if curLine.Trimed == "" {
			curLine.LineType = txeh.EMPTY
			continue
		}

		curLineSplit := strings.SplitN(curLine.Trimed, "#", 2)
		if len(curLineSplit) > 1 {
			curLine.Comment = curLineSplit[1]
		}
		curLine.Trimed = curLineSplit[0]

		curLine.Parts = strings.Fields(curLine.Trimed)

		if len(curLine.Parts) > 1 {
			curLine.LineType = txeh.ADDRESS
			curLine.Address = strings.ToLower(curLine.Parts[0])
			for _, p := range curLine.Parts[1:] {
				curLine.Hostnames = append(curLine.Hostnames, strings.ToLower(p))
			}

			continue
		}

		curLine.LineType = txeh.UNKNOWN

	}

	return hostFileLines, nil
}

// LoginBuildPush 登录仓库编译镜像并推送
func LoginBuildPush() {
	log.Info("插件开始执行......")

	selectOp := api.GetInputParam("selectOp")

	if selectOp == "login_build_push" {
		loginAndBuildPush()
	}
}
