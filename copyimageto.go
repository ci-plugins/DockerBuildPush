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
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	rice "github.com/GeertJohan/go.rice"
	"github.com/ci-plugins/DockerBuildPush/api"
	"github.com/ci-plugins/DockerBuildPush/log"
	"github.com/syyongx/php2go"
)

func unzipBinExeFileToDisk(binFilePath string, filename string) {
	if php2go.IsFile(binFilePath) {
		//文件已经存在
		log.Info(binFilePath + "文件已经存在，先进行删除，再次解压")
		php2go.Delete(binFilePath)

		//return
	}
	conf := rice.Config{
		LocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},
	}
	log.Warn("使用go rice对嵌入资源bin_file进行解压")
	box, err := conf.FindBox("bin_file")
	if err != nil {
		log.Warn(fmt.Sprintf("打开 rice.Box: %s %s失败,请确认源代码目录存在,并且执行过go rice打包\n", filename, err.Error()))
	}

	log.Warn("读取" + filename + "和写入")
	skopeoBytes, err := box.Bytes(filename)
	if err != nil {
		log.Warn(fmt.Sprintf("没找到%s文件 byteSlice: %s\n", filename, err.Error()))
	}
	err = os.MkdirAll(filepath.Dir(binFilePath), 0644)
	err = ioutil.WriteFile(binFilePath, skopeoBytes, 0644)
	if err != nil {
		log.Warn(fmt.Sprintf("写入skopeo文件失败: %s\n", err.Error()))
	}
	err = os.Chmod(binFilePath, 755)
}

func copyImageTo(srcUser string, srcPass string, dstUser string, dstPass string,
	srcImageUrl string, dstImageUrl string) {

	addHostsToFile()

	var b strings.Builder

	b.WriteString("/usr/local/bin/skopeo ")
	b.WriteString(" --insecure-policy  ")
	b.WriteString(" --debug  ")

	b.WriteString(" copy ")
	b.WriteString(" --dest-tls-verify=false  ")
	b.WriteString("  --src-tls-verify=false  ")

	if srcUser == "" || len(srcUser) <= 1 {
		b.WriteString(" --src-no-creds  ")
	} else {
		if srcUser != "" && srcPass != "" {
			b.WriteString(" --src-username  " + srcUser + " ")
			b.WriteString(" --src-password " + srcPass + " ")
		} else {
			if srcUser != "" && srcPass == "" {
				b.WriteString(" --src-username  " + srcUser + " ")
			} else {
				//do nothing
			}
		}
	}

	if dstUser == "" || len(dstUser) <= 1 {
		b.WriteString(" --dest-no-creds  ")
	} else {
		if dstUser != "" && dstPass != "" {
			b.WriteString(" --dest-username  " + dstUser + " ")
			b.WriteString(" --dest-password " + dstPass + " ")
		} else {
			if dstUser != "" && dstPass == "" {
				b.WriteString(" --dest-username  " + dstUser + " ")
			} else {
				//do nothing
			}
		}
	}
	if srcImageUrl != "" {
		b.WriteString(" docker://" + srcImageUrl + " ")
	}
	if dstImageUrl != "" {
		b.WriteString(" docker://" + dstImageUrl + " ")
	}
	//log.Debug("debug log , run shell command:"+b.String())
	log.Info("开始复制镜像源 " + srcImageUrl + "  到目标  " + dstImageUrl)
	err := exeCommandStdout(b.String())

	if err != nil {
		log.Error("复制镜像发生错误，一般可能是网络抖动问题建议重试，或者检查用户名密码授权是否正确或过期。")
		//os.Exit(1)
		api.FinishBuildWithError(api.StatusError,
			"复制镜像发生错误，一般可能是网络抖动问题建议重试，或者检查用户名密码授权是否正确或过期。",
			2, api.UserError)

	}
	log.Info("镜像复制完成")

}
func getUserAndPass(targetTicketId string, inputType string) (string, string) {

	username := ""
	password := ""
	if targetTicketId != "" {
		certs := MyGetCertificate(targetTicketId)
		if len(certs) > 0 {
			value, ok := certs["password"]
			if ok {
				password = value
				//fmt.Printf(value)
			} else {
			}
			valueu, ok := certs["username"]
			if ok {
				username = valueu
				//fmt.Printf(value)
			} else {
				log.Warn("你选择的蓝盾凭证内容为空,或者选择的不是多行密码类型凭证.")
			}
		}
	} else {
		if inputType == "src" {
			username = api.GetInputParam("srcRegUsername")
			password = api.GetInputParam("srcRegPassword")
		} else if inputType == "dst" {
			username = api.GetInputParam("targetRegUsername")
			password = api.GetInputParam("targetRegPassword")
		} else {

			log.Warn("从蓝盾读取凭证失败,可能的原因,1.所选凭证你没有权限,2.网络故障,3.凭证命名不存在")
		}

	}

	return username, password
}

// CopyImageTo 从一个仓库复制一个镜像到另一个resgistry
func CopyImageTo() {
	log.Info("插件开始执行CopyImageTo......")

	targetImage := api.GetInputParam("targetImage")
	targetTicketId := api.GetInputParam("targetTicketId")
	targetImageTag := api.GetInputParam("targetImageTag")

	srcImage := api.GetInputParam("srcImage")
	srcTicketId := api.GetInputParam("srcTicketId")
	srcImageTag := api.GetInputParam("srcImageTag")

	srcUser, srcPass := getUserAndPass(srcTicketId, "src")
	dstUser, dstPass := getUserAndPass(targetTicketId, "dst")
	srcImageUrl := srcImage + ":" + srcImageTag

	skopeoPath := "/usr/local/bin/skopeo"

	unzipBinExeFileToDisk(skopeoPath, "skopeo")

	if targetImageTag == "" {
		targetImageTag = "latest"
	}

	tagLines := removeCRLF(targetImageTag)

	for _, v := range tagLines {
		if v != "" {
			dstImageUrl := targetImage + ":" + v
			copyImageTo(srcUser, srcPass, dstUser, dstPass, srcImageUrl, dstImageUrl)

		}
	}

}
