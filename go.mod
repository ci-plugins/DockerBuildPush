module github.com/ci-plugins/DockerBuildPush

go 1.14

require (
	github.com/GeertJohan/go.rice v1.0.2
	github.com/syyongx/php2go v0.9.4
	github.com/tidwall/gjson v1.11.0 // indirect
	github.com/tidwall/sjson v1.2.3
	github.com/txn2/txeh v1.3.0
)

replace (
	github.com/Azure/go-autorest => github.com/Azure/go-autorest v14.2.0+incompatible
	github.com/containerd/containerd v1.4.0-0.20191014053712-acdcf13d5eaf => github.com/containerd/containerd v0.0.0-20191014053712-acdcf13d5eaf
	github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c => github.com/docker/docker v0.0.0-20190319215453-e7b5f7dbe98c
	github.com/tonistiigi/fsutil v0.0.0-20190819224149-3d2716dd0a4d => github.com/tonistiigi/fsutil v0.0.0-20191018213012-0f039a052ca1
)