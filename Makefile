CURRENT_DIR=/data/landun/workspace/src/DockerBuildPush
BUILD_TAGS="exclude_graphdriver_devicemapper exclude_graphdriver_btrfs containers_image_openpgp"
BUILD_LD= -ldflags "-extldflags '-static'"
all:clean skopeo_build kaniko_build rice_build executor_build open_source

proot_build:
	echo "start build proot"
	cd ${CURRENT_DIR}/ && ./build_proot.sh
	cp -f ${CURRENT_DIR}/proot-5.1.0/src/proot   ${CURRENT_DIR}/bin_file/proot

proot52_build:
	curl -L -o ${CURRENT_DIR}/bin_file/proot https://github.com/proot-me/proot/releases/download/v5.2.0/proot-v5.2.0-x86_64-static
skopeo_build:
	echo "start build skopeo"
	go get -d github.com/containers/skopeo/cmd/skopeo
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v -ldflags '-extldflags "-static" ' -tags ${BUILD_TAGS} -o ${CURRENT_DIR}/bin_file/skopeo github.com/containers/skopeo/cmd/skopeo
kaniko_build:
	echo "start build kaniko"
	cd ${CURRENT_DIR} && git clone https://github.com/GoogleContainerTools/kaniko
	cd ${CURRENT_DIR}/kaniko && GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v -ldflags '-extldflags "-static" ' -o ${CURRENT_DIR}/bin_file/executor  github.com/GoogleContainerTools/kaniko/cmd/executor
rice_build:
	echo "start build rice"
	go get -d github.com/GeertJohan/go.rice/rice
	cd ${CURRENT_DIR} && go build -v  -o ${CURRENT_DIR}/rice github.com/GeertJohan/go.rice/rice
	cd ${CURRENT_DIR} &&  ${CURRENT_DIR}/rice -i github.com/ci-plugins/bkci-DockerBuildPush  embed-go
	echo "finished rice embed"
executor_build:
	echo "start build kaniko"
	mkdir -p ${CURRENT_DIR}/../../bin
	go mod tidy -v
	cd ${CURRENT_DIR} && GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static" ' -tags ${BUILD_TAGS}  -o out/executor github.com/ci-plugins/bkci-DockerBuildPush
	cp -f -v -r  ${CURRENT_DIR}/out/executor  ${CURRENT_DIR}/../../bin/app
	cp -f -v -r  ${CURRENT_DIR}/task.json  ${CURRENT_DIR}/../../bin/task.json
	echo "finished build kaniko"

clean:
	rm -rf ${CURRENT_DIR}/../../bin
	rm -rf ${CURRENT_DIR}/out/executor
	rm -rf ${CURRENT_DIR}/rice
	rm -rf ${CURRENT_DIR}/bin_file/executor
	rm -rf ${CURRENT_DIR}/bin_file/skopeo
	rm -rf ${CURRENT_DIR}/DockerBuildPushGo
	rm -rf ${CURRENT_DIR}/kaniko
	mkdir -p ${CURRENT_DIR}/bin_file


open_source:
	cd ${CURRENT_DIR}/../../bin && zip app.zip app task.json