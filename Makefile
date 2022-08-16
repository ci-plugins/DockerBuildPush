CURRENT_DIR=/data/landun/workspace/src/DockerBuildPushGo
BUILD_TAGS="exclude_graphdriver_devicemapper exclude_graphdriver_btrfs containers_image_openpgp"
BUILD_LD= -ldflags "-extldflags '-static'"
all:clean golang_update proot53_build skopeo_build kaniko_build rice_build executor_build open_source

golang_update:
	echo "update bkci build golang 1.14 to 1.16 or higher"
	if [ -e  "/usr/bin/yum" ]; then  rm -rf /var/lib/rpm/.dbenv.lock &&  mkdir -p /temprpmdb && rpm --rebuilddb --dbpath=/temprpmdb/rpm && mv -f /temprpmdb/rpm/* /var/lib/rpm  && yum install golang -y; else echo "can not find /usr/bin/yum"; fi
proot_build:
	echo "start build proot"
	cd ${CURRENT_DIR}/ && ./build_proot.sh
	cp -f ${CURRENT_DIR}/proot-5.1.0/src/proot   ${CURRENT_DIR}/bin_file/proot
proot53_build:
	curl -L -o ${CURRENT_DIR}/bin_file/proot https://github.com/proot-me/proot/releases/download/v5.3.0/proot-v5.3.0-x86_64-static
skopeo_build:
	echo "start build skopeo"
	go get -d github.com/containers/skopeo/cmd/skopeo
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v -ldflags '-extldflags "-static" ' -tags ${BUILD_TAGS} -o ${CURRENT_DIR}/bin_file/skopeo github.com/containers/skopeo/cmd/skopeo
	echo "skopeo build finished"
	ls -la ${CURRENT_DIR}/bin_file/
kaniko_build:
	echo "start build kaniko"
	#go get -d github.com/GoogleContainerTools/kaniko/cmd/executor
	cd ${CURRENT_DIR} && git clone http://github.com/GoogleContainerTools/kaniko
	cd ${CURRENT_DIR}/kaniko && GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -v -ldflags '-extldflags "-static" ' -o ${CURRENT_DIR}/bin_file/executor  github.com/GoogleContainerTools/kaniko/cmd/executor
	echo "kaniko build finished"
	ls -la ${CURRENT_DIR}/bin_file/
rice_build:
	echo "start build rice"
	go get -d github.com/GeertJohan/go.rice/rice
	cd ${CURRENT_DIR} && go build -v  -o ${CURRENT_DIR}/rice github.com/GeertJohan/go.rice/rice
	cd ${CURRENT_DIR} &&  ${CURRENT_DIR}/rice -i github.com/ci-plugins/DockerBuildPush  embed-go
	echo "finished rice embed"
executor_build:
	echo "start build kaniko"
	mkdir -p ${CURRENT_DIR}/../../bin
	go mod tidy -v
	cd ${CURRENT_DIR} && GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags '-extldflags "-static" ' -tags ${BUILD_TAGS}  -o out/executor github.com/ci-plugins/DockerBuildPush
	cp -f -v -r  ${CURRENT_DIR}/out/executor  ${CURRENT_DIR}/../../bin/app
	cp -f -v -r  ${CURRENT_DIR}/task.json  ${CURRENT_DIR}/../../bin/task.json
	echo "finished build kaniko"

clean:
	rm -rf ${CURRENT_DIR}/../../bin
	rm -rf ${CURRENT_DIR}/out/executor
	rm -rf ${CURRENT_DIR}/rice
	rm -rf ${CURRENT_DIR}/bin_file/executor
	rm -rf ${CURRENT_DIR}/bin_file/skopeo
	rm -rf ${CURRENT_DIR}/DockerBuildPush
	rm -rf ${CURRENT_DIR}/kaniko
	mkdir -p ${CURRENT_DIR}/bin_file


open_source:
	cd ${CURRENT_DIR}/../../bin && zip app.zip app task.json