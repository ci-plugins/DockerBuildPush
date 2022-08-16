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
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/ci-plugins/DockerBuildPush/api"
	"github.com/ci-plugins/DockerBuildPush/log"
)

//SelectOP 动作执行分支选择和输入参数粗略检测
func SelectOP() {

	selectOp := api.GetInputParam("selectOp")
	targetImage := api.GetInputParam("targetImage")
	targetImageName := api.GetInputParam("targetImageName")
	targetTicketId := api.GetInputParam("targetTicketId")
	targetImageTag := api.GetInputParam("targetImageTag")
	dockerBuildDir := api.GetInputParam("dockerBuildDir")
	dockerFilePath := api.GetInputParam("dockerFilePath")
	dockerBuildArgs := api.GetInputParam("dockerBuildArgs")
	dockerBuildHosts := api.GetInputParam("dockerBuildHosts")
	targetRepoItemStr := api.GetInputParam("targetRepoItemStr")
	sourceRepoItemsStr := api.GetInputParam("sourceRepoItemsStr")
	dockerCommand := api.GetInputParam("dockerCommand")
	//sourceMirrorTicketPair = api.GetInputParam("sourceMirrorTicketPair")
	selectOp = strings.TrimSpace(selectOp)
	targetImage = strings.TrimSpace(targetImage)
	targetImageName = strings.TrimSpace(targetImageName)
	targetTicketId = strings.TrimSpace(targetTicketId)
	targetImageTag = strings.TrimSpace(targetImageTag)
	dockerBuildDir = strings.TrimSpace(dockerBuildDir)
	dockerFilePath = strings.TrimSpace(dockerFilePath)
	dockerBuildArgs = strings.TrimSpace(dockerBuildArgs)
	dockerBuildHosts = strings.TrimSpace(dockerBuildHosts)
	targetRepoItemStr = strings.TrimSpace(targetRepoItemStr)
	sourceRepoItemsStr = strings.TrimSpace(sourceRepoItemsStr)
	dockerCommand = strings.TrimSpace(dockerCommand)

	if selectOp == "" {
		selectOp = "login_build_push"
	}

	if selectOp == "login_build_push" {
		LoginBuildPush()
	}

	if selectOp == "copy_image_to" {
		CopyImageTo()
	}

	log.Info("有问题看日志仍解决不了,请联系插件作者reynaldliu进行协助.")

}

func checkContains(in []string, check string) bool {
	for _, v := range in {
		if v == check {
			return true
		}
	}
	return false
}

// FormatStack 格式化stack trace
func FormatStack(in map[string]string) string {
	var sb strings.Builder
	sb.WriteString("")
	for k, v := range in {
		sb.WriteString(k + "\n")
		sb.WriteString("\t" + v + "\n")
	}
	return sb.String()
}

// ParseStack stacktrace
func ParseStack(stackStr string) map[string]string {
	retObj := make(map[string]string)
	strArray := strings.Split(stackStr, "\n")
	maxLen := len(strArray)
	ignoreStr := []string{"runtime/debug.Stack()",
		"runtime/debug.PrintStack()",
		"SelectOP()",
		"main.main()", "main.main.func1()"}
	for k, _ := range strArray {
		//出现第一个()号是函数行,不是最后一行越界
		if k <= maxLen-1 {
			if strings.Contains(strArray[k], "(") && strings.Contains(strArray[k], ")") {
				//他的下一行也没有越界，并且包含.go
				if k <= maxLen-1-1 {
					if strings.Contains(strArray[k+1], ".go") {
						if _, ok := retObj[strArray[k]]; !ok {
							if !checkContains(ignoreStr, strArray[k]) {
								retObj[strArray[k]] = strArray[k+1]
							}
						}
					}
				}
			}
		}

	}
	return retObj
}
func main() {
	runtime.GOMAXPROCS(4)
	defer func() {
		if err := recover(); err != nil {
			errStack := FormatStack(ParseStack(string(debug.Stack())))

			log.Error("panic: ", err)
			log.Info("有问题看日志仍解决不了,请联系插件作者reynaldliu进行协助.")
			//api.FinishBuild(api.StatusError, "panic occurs")

			api.FinishBuildWithError(api.StatusError,
				"插件级系统错误"+errStack,
				2, api.PluginError)

		}
	}()
	SelectOP()
}
