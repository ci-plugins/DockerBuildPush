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

const (
	AuthHeaderBuildId              = "X-SODA-BID"
	AuthHeaderProjectId            = "X-SODA-PID"
	AuthHeaderDevopsBuildType      = "X-DEVOPS-BUILD-TYPE"
	AuthHeaderDevopsProjectId      = "X-DEVOPS-PROJECT-ID"
	AuthHeaderDevopsBuildId        = "X-DEVOPS-BUILD-ID"
	AuthHeaderDevopsVmSeqId        = "X-DEVOPS-VM-SID"
	AuthHeaderDevopsAgentId        = "X-DEVOPS-AGENT-ID"
	AuthHeaderDevopsAgentSecretKey = "X-DEVOPS-AGENT-SECRET-KEY"
)

const (
	DataDirEnv    = "bk_data_dir"
	InputFileEnv  = "bk_data_input"
	OutputFileEnv = "bk_data_output"
)

// SdkEnv sdk环境变量输入参数
type SdkEnv struct {
	BuildType string `json:"buildType"`
	ProjectId string `json:"projectId"`
	AgentId   string `json:"agentId"`
	SecretKey string `json:"secretKey"`
	Gateway   string `json:"gateway"`
	BuildId   string `json:"buildId"`
	VmSeqId   string `json:"vmSeqId"`
}

// AtomBaseParam 流水线内部变量
type AtomBaseParam struct {
	PipelineVersion        string `json:"pipeline.version"`
	ProjectName            string `json:"project.name"`
	ProjectNameCn          string `json:"project.name.chinese"`
	PipelineId             string `json:"pipeline.id"`
	PipelineBuildNum       string `json:"pipeline.build.num"`
	PipelineBuildId        string `json:"pipeline.build.id"`
	PipelineName           string `json:"pipeline.name"`
	PipelineStartTimeMills string `json:"pipeline.time.start"`
	PipelineStartType      string `json:"pipeline.start.type"`
	PipelineStartUserId    string `json:"pipeline.start.user.id"`
	PipelineStartUserName  string `json:"pipeline.start.user.name"`
	BkWorkspace            string `json:"bkWorkspace"`
}

// BuildType 构建类型
type BuildType string

const (
	BuildTypeWorker      = "WORKER"
	BuildTypeAgent       = "AGENT"
	BuildTypePluginAgent = "PLUGIN_AGENT"
	BuildTypeDocker      = "DOCKER"
	BuildTypeDockerHost  = "DOCKER_HOST"
	BuildTypeTstackAgent = "TSTACK_AGENT"
)

// DataType 输出类型
type DataType string

const (
	DataTypeString   DataType = "string"
	DataTypeArtifact DataType = "artifact"
	DataTypeReport   DataType = "report"
)

// Status 插件执行状态
type Status string

// 插件执行状态
const (
	StatusSuccess Status = "success"
	StatusFailure Status = "failure"
	StatusError   Status = "error"
)

// ErrorType 执行错误类型
type ErrorType int

// 插件执行状态
const (
	UserError       ErrorType = 1
	ThirdPartyError ErrorType = 2
	PluginError     ErrorType = 3
)

// ArtifactData 构件类型
type ArtifactData struct {
	Type  DataType `json:"type"`
	Value []string `json:"value"`
}

// AddArtifact 添加构件
func (a *ArtifactData) AddArtifact(artifact string) {
	a.Value = append(a.Value, artifact)
}

// AddArtifactAll 添加多个构件
func (a *ArtifactData) AddArtifactAll(artifacts []string) {
	a.Value = append(a.Value, artifacts...)
}

// StringData 字符型
type StringData struct {
	Type  DataType `json:"type"`
	Value string   `json:"value"`
}

// ReportData 报告型
type ReportData struct {
	Type   DataType `json:"type"`
	Label  string   `json:"label"`
	Path   string   `json:"path"`
	Target string   `json:"target"`
}

// NewReportData 报告
func NewReportData(label string, path string, target string) *ReportData {
	return &ReportData{
		Type:   DataTypeReport,
		Label:  label,
		Path:   path,
		Target: target,
	}
}

// NewStringData 字符
func NewStringData(value string) *StringData {
	return &StringData{
		Type:  DataTypeString,
		Value: value,
	}
}

// NewArtifactData 新构件
func NewArtifactData() *ArtifactData {
	return &ArtifactData{
		Type:  DataTypeArtifact,
		Value: []string{},
	}
}

// Qualitydata 质量红线数据
type Qualitydata struct {
	Value string `json:"value"`
}

// String 质量红线数据
func (a *Qualitydata) String() string {
	return a.Value
}

// NewQualityData 创建质量红线数据
func NewQualityData(value string) *Qualitydata {
	return &Qualitydata{
		Value: value,
	}
}

// AtomOutput 插件输出
type AtomOutput struct {
	Status            Status                  `json:"status"`
	Message           string                  `json:"message"`
	ErrorCode         int                     `json:"errorCode"`
	ErrorType         ErrorType               `json:"errorType"`
	Type              string                  `json:"type"`
	Data              map[string]interface{}  `json:"data"`
	QualityData       map[string]*Qualitydata `json:"qualityData"`
	PlatformCode      string                  `json:"platformCode"`
	PlatformErrorCode int                     `json:"platformErrorCode"`
}

// NewAtomOutput 创建插件输出
func NewAtomOutput() *AtomOutput {
	output := new(AtomOutput)
	output.Status = StatusSuccess
	output.Message = "success"
	output.Type = "default"
	output.Data = make(map[string]interface{})
	output.QualityData = make(map[string]*Qualitydata)
	return output
}
