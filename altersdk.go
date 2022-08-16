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

// 蓝盾SDK获取凭证提取出来合并到开源SDK,所有相关函数添加my开头
import (
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ci-plugins/DockerBuildPush/api"

	"github.com/ci-plugins/DockerBuildPush/log"
)

// http head keys
const (
	AuthHeaderBuildId              = "X-SODA-BID"
	AuthHeaderProjectId            = "X-SODA-PID"
	AuthHeaderDevopsBuildType      = "X-DEVOPS-BUILD-TYPE"
	AuthHeaderDevopsProjectId      = "X-DEVOPS-PROJECT-ID"
	AuthHeaderDevopsBuildId        = "X-DEVOPS-BUILD-ID"
	AuthHeaderDevopsVmSeqId        = "X-DEVOPS-VM-SID"
	AuthHeaderDevopsVmSeqName      = "X-DEVOPS-VM-NAME"
	AuthHeaderDevopsAgentId        = "X-DEVOPS-AGENT-ID"
	AuthHeaderDevopsAgentSecretKey = "X-DEVOPS-AGENT-SECRET-KEY"
)

// MyCertificate   请求凭证结果
type MyCertificate struct {
	status int
	Data   map[string]string `json:"data"`
}

// MyGetCertificate   获取指定ID的凭证
func MyGetCertificate(certificateId string) map[string]string {

	//log.Info("Begin to get certificate")
	url := mybuildUrl("/ticket/api/build/credentials/" + certificateId + "/detail")
	var build = MyBuildRequest{path: url, requestBody: nil, headers: mygetAllHeaders()}
	req, err := mybuildGet(build)
	if err != nil {
		log.Error("build request failed: " + err.Error())
		log.Warn("获取蓝盾内置凭证不成功,请降级使用手工输入的方式.")

		return nil
	}

	respByte, err := myrequest(*req, "failed to get certificate")
	if err != nil {
		log.Error("get certificate failed: " + err.Error())
		log.Warn("获取蓝盾内置凭证不成功,请降级使用手工输入的方式.")
		return nil
	}
	var certificate = new(MyCertificate)
	err = json.Unmarshal(respByte, &certificate)

	if err != nil {
		log.Warn("获取蓝盾内置凭证不成功,请降级使用手工输入的方式.")
		log.Error("resolve response json failed: " + err.Error())
	}

	return certificate.Data
}

// MyBuildRequest   构建请求
type MyBuildRequest struct {
	path        string
	headers     map[string]string
	requestBody io.Reader
}

var myclient = http.Client{
	Timeout: 30 * time.Second,
}

func myrequest(r http.Request, errMessage string) ([]byte, error) {
	response, err := myclient.Do(&r)
	if err != nil {
		log.Error("do http request failed: " + err.Error())
		return nil, errors.New(errMessage)
	}

	if !(response.StatusCode >= 200 && response.StatusCode < 300) {
		log.Error("http request failed, status: " + strconv.Itoa(response.StatusCode))
		return nil, errors.New(errMessage)
	}

	respStr, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error("get response content failed: " + err.Error())
		return nil, errors.New(errMessage)
	}

	return respStr, nil
}

func mybuildGet(build MyBuildRequest) (*http.Request, error) {
	if build.path == "" {
		return nil, errors.New("can not generate request without path")
	}

	url := mybuildUrl(build.path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	if build.headers != nil {
		for k, v := range build.headers {
			req.Header.Add(k, v)
		}
	}

	return req, nil
}

func mygetAllHeaders() map[string]string {
	var headers = make(map[string]string)
	headers[AuthHeaderDevopsBuildType] = api.GlobalSdkEvn.BuildType
	headers[AuthHeaderProjectId] = api.GlobalSdkEvn.ProjectId
	headers[AuthHeaderDevopsProjectId] = api.GlobalSdkEvn.ProjectId
	headers[AuthHeaderDevopsBuildId] = api.GlobalSdkEvn.BuildId
	headers[AuthHeaderDevopsAgentSecretKey] = api.GlobalSdkEvn.SecretKey
	headers[AuthHeaderDevopsAgentId] = api.GlobalSdkEvn.AgentId
	headers[AuthHeaderDevopsVmSeqId] = api.GlobalSdkEvn.VmSeqId
	headers[AuthHeaderBuildId] = api.GlobalSdkEvn.BuildId
	return headers
}

func mybuildUrl(path string) string {
	var gateway = strings.TrimSuffix(api.GlobalSdkEvn.Gateway, "/")
	if strings.HasPrefix(gateway, "http") {
		return gateway + "/" + strings.TrimPrefix(strings.TrimSpace(path), "/")
	} else {
		return "http://" + gateway + "/" + strings.TrimPrefix(strings.TrimSpace(path), "/")
	}
}
