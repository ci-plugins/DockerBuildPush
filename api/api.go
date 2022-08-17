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

package api

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"strings"

	"github.com/ci-plugins/DockerBuildPush/log"
)

// GlobalSdkEvn
var GlobalSdkEvn *SdkEnv
var gAtomBaseParam *AtomBaseParam
var gAllAtomParam map[string]interface{}
var gAtomOutput *AtomOutput

var gDataDir string
var gInputFile string
var gOutputFile string

// init
func init() {
	gAtomOutput = NewAtomOutput()
	gDataDir = getDataDir()
	gInputFile = getInputFile()
	gOutputFile = getOutputFile()
	initSdkEnv()
	initAtomParam()
}

// initAtomParam
func initAtomParam() {
	err := LoadInputParam(&gAllAtomParam)
	if err != nil {
		log.Error("init atom base param failed: ", err.Error())
		FinishBuildWithErrorCode(StatusError, "init atom base param failed", 16015100)
	}

	gAtomBaseParam = new(AtomBaseParam)
	err = LoadInputParam(gAtomBaseParam)
	if err != nil {
		log.Error("init atom base param failed: ", err.Error())
		FinishBuildWithErrorCode(StatusError, "init atom base param failed", 16015100)
	}
}

// GetInputParam 获取外部输入参数
func GetInputParam(name string) string {
	value := gAllAtomParam[name]
	if value == nil {
		return ""
	}
	strValue, ok := value.(string)
	if !ok {
		return ""
	}
	return strValue
}

// LoadInputParam 从json文件读取输入参数
func LoadInputParam(v interface{}) error {
	file := gDataDir + "/" + gInputFile
	//log.Debug("load input.json file:" + file)

	data, err := ioutil.ReadFile(file)

	if err != nil {
		log.Error("load input param failed:", err.Error())
		return errors.New("load input param failed")
	}
	//log.Debug("input data:" + string(data))
	err = json.Unmarshal(data, v)
	if err != nil {
		log.Error("parse input param failed:", err.Error())
		return errors.New("parse input param failed")
	}
	return nil
}

// initSdkEnv
func initSdkEnv() {
	filePath := gDataDir + "/.sdk.json"
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Warn("警告信息,可以无视,read .sdk.json failed: ", err.Error())
		//FinishBuildWithErrorCode(StatusError, "read .sdk.json failed", 16015102)
		log.Warn("警告信息,可以无视,.sdk.json not found,guess running not in landun/bkci/streamci \n" +
			"will use a fake .sdk.json debug string\n" + "if you run ./app  command line ignore .sdk.json file.")
		jsonData := []byte("{\"buildType\": \"DOCKER\" ,\"projectId\": \"232\",\"agentId\": \"1x\"," +
			"\"secretKey\": \"2323232\",\"gateway\": \"testurl\",\"buildId\": \"vvadsaaf\",\"vmSeqId\": \"33223\"}")
		data = jsonData
	}

	GlobalSdkEvn = new(SdkEnv)
	err = json.Unmarshal(data, GlobalSdkEvn)
	if err != nil {
		log.Warn("警告信息,可以无视,parse .sdk.json failed: ", err.Error())
		FinishBuildWithErrorCode(StatusError, "parse .sdk.json failed", 16015102)
	}

	os.Remove(filePath)
}

// getDataDir
func getDataDir() string {
	dir := strings.TrimSpace(os.Getenv(DataDirEnv))
	if len(dir) == 0 {
		dir, _ = os.Getwd()
	}
	return dir
}

// getInputFile
func getInputFile() string {
	file := strings.TrimSpace(os.Getenv(InputFileEnv))
	if len(file) == 0 {
		file = "input.json"
	}
	return file
}

// getOutputFile
func getOutputFile() string {
	file := strings.TrimSpace(os.Getenv(OutputFileEnv))
	if len(file) == 0 {
		file = "output.json"
	}
	return file
}

// GetOutputData 获取输入参数键值
func GetOutputData(key string) interface{} {
	return gAtomOutput.Data[key]
}

// AddOutputData 组装输出数据
func AddOutputData(key string, data interface{}) {
	gAtomOutput.Data[key] = data
}

// RemoveOutputData 删除输出临时文件
func RemoveOutputData(key string) {
	delete(gAtomOutput.Data, key)
}

// WriteOutput 写入输出文件
func WriteOutput() error {
	data, _ := json.Marshal(gAtomOutput)

	file := gDataDir + "/" + gOutputFile
	err := ioutil.WriteFile(file, data, 0644)
	if err != nil {
		log.Error("write output failed: ", err.Error())
		return errors.New("write output failed")
	}
	return nil
}

// FinishBuild 后置hook完成编译
func FinishBuild(status Status, msg string) {
	gAtomOutput.Message = msg
	gAtomOutput.Status = status
	WriteOutput()
	switch status {
	case StatusSuccess:
		os.Exit(0)
	case StatusFailure:
		os.Exit(1)
	case StatusError:
		os.Exit(2)
	default:
		os.Exit(0)
	}
}

// FinishBuildWithErrorCode 完成编译并指定退出码
func FinishBuildWithErrorCode(status Status, msg string, errorCode int) {
	gAtomOutput.Message = msg
	gAtomOutput.Status = status
	gAtomOutput.ErrorCode = errorCode
	WriteOutput()
	switch status {
	case StatusSuccess:
		os.Exit(0)
	case StatusFailure:
		os.Exit(1)
	case StatusError:
		os.Exit(2)
	default:
		os.Exit(0)
	}
}

// FinishBuildWithError 结束构建
// @status		任务状态
// @msg			消息
// @errorCode	错误码
// @errorType	错误类型
func FinishBuildWithError(status Status, msg string, errorCode int, errorType ErrorType) {
	gAtomOutput.Message = msg
	gAtomOutput.Status = status
	gAtomOutput.ErrorCode = errorCode
	gAtomOutput.ErrorType = errorType
	WriteOutput()
	switch status {
	case StatusSuccess:
		os.Exit(0)
	case StatusFailure:
		os.Exit(1)
	case StatusError:
		os.Exit(2)
	default:
		os.Exit(0)
	}
}

// SetPlatformCode 设置插件对接平台代码
func SetPlatformCode(platformCode string) {
	gAtomOutput.PlatformCode = platformCode
}

// SetPlatformErrorCode 设置插件对接平台错误码
func SetPlatformErrorCode(platformErrorCode int) {
	gAtomOutput.PlatformErrorCode = platformErrorCode
}

// SetAtomOutputType 设置输出参数类型
func SetAtomOutputType(atomOutputType string) {
	gAtomOutput.Type = atomOutputType
}

// GetProjectName 获取当前构件项目内部名
func GetProjectName() string {
	return gAtomBaseParam.ProjectName
}

// GetProjectDisplayName 获取当前构件项目中文名
func GetProjectDisplayName() string {
	return gAtomBaseParam.ProjectNameCn
}

// GetPipelineId 获取流水线ID
func GetPipelineId() string {
	return gAtomBaseParam.PipelineId
}

// GetPipelineName 获取流水线名
func GetPipelineName() string {
	return gAtomBaseParam.PipelineName
}

// GetPipelineBuildId 获取流水线构建ID
func GetPipelineBuildId() string {
	return gAtomBaseParam.PipelineBuildId
}

// GetPipelineBuildNumber 获取流水线构件号
func GetPipelineBuildNumber() string {
	return gAtomBaseParam.PipelineBuildNum
}

// GetPipelineStartType 获取流水线启动类型
func GetPipelineStartType() string {
	return gAtomBaseParam.PipelineStartType
}

// GetPipelineStartUserId 获取流水线启动人ID
func GetPipelineStartUserId() string {
	return gAtomBaseParam.PipelineStartUserId
}

// GetPipelineStartUserName 获取流水线启动人名字
func GetPipelineStartUserName() string {
	return gAtomBaseParam.PipelineStartUserName
}

// GetPipelineStartTimeMills 获取流水线开始运行的时间点
func GetPipelineStartTimeMills() string {
	return gAtomBaseParam.PipelineStartTimeMills
}

// GetPipelineVersion 获取流水线版本号
func GetPipelineVersion() string {
	return gAtomBaseParam.PipelineVersion
}

// GetWorkspace 获取构件机工作目录
func GetWorkspace() string {
	return gAtomBaseParam.BkWorkspace
}
