module github.com/octelium/octelium/client/octelium

go 1.23.8

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2
	github.com/creack/pty v1.1.24
	github.com/fatih/color v1.18.0
	github.com/google/uuid v1.6.0
	github.com/manifoldco/promptui v0.9.0
	github.com/miekg/dns v1.1.62
	github.com/moby/term v0.5.0
	github.com/octelium/octelium/apis v0.0.0-00010101000000-000000000000
	github.com/octelium/octelium/client/common v0.0.0-00010101000000-000000000000
	github.com/octelium/octelium/pkg v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/quic-go/quic-go v0.48.2
	github.com/spf13/cobra v1.8.1
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.10.0
	github.com/tj/assert v0.0.3
	github.com/vishvananda/netlink v1.3.0
	go.uber.org/zap v1.27.0
	golang.org/x/crypto v0.38.0
	golang.org/x/net v0.40.0
	golang.org/x/sys v0.33.0
	golang.zx2c4.com/wireguard v0.0.0-20231211153847-12269c276173
	golang.zx2c4.com/wireguard/wgctrl v0.0.0-20230429144221-925a1e7659e6
	golang.zx2c4.com/wireguard/windows v0.5.3
	google.golang.org/protobuf v1.36.1
	gvisor.dev/gvisor v0.0.0-20240722211153-64c016c92987
)

require (
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/chzyer/readline v0.0.0-20180603132655-2972be24d48e // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/denisbrodbeck/machineid v1.0.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-ole/go-ole v1.2.6 // indirect
	github.com/go-resty/resty/v2 v2.16.5 // indirect
	github.com/go-task/slim-sprig v0.0.0-20230315185526-52ccab3ef572 // indirect
	github.com/gofrs/flock v0.8.1 // indirect
	github.com/google/btree v1.1.2 // indirect
	github.com/google/certificate-transparency-go v1.1.2 // indirect
	github.com/google/go-attestation v0.5.1 // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/google/go-tpm v0.9.0 // indirect
	github.com/google/go-tspi v0.3.0 // indirect
	github.com/google/pprof v0.0.0-20210720184732-4bb14d4b1be1 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.4.0 // indirect
	github.com/hashicorp/go-version v1.2.0 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/josharian/native v1.1.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.13 // indirect
	github.com/mdlayher/genetlink v1.3.2 // indirect
	github.com/mdlayher/netlink v1.7.2 // indirect
	github.com/mdlayher/socket v0.4.1 // indirect
	github.com/octelium/octelium/octelium-go v0.0.0-00010101000000-000000000000 // indirect
	github.com/olekukonko/tablewriter v0.0.5 // indirect
	github.com/onsi/ginkgo/v2 v2.9.5 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/vishvananda/netns v0.0.4 // indirect
	github.com/yusufpapurcu/wmi v1.2.4 // indirect
	github.com/zcalusic/sysinfo v1.0.2 // indirect
	go.uber.org/mock v0.4.0 // indirect
	go.uber.org/multierr v1.11.0 // indirect
	golang.org/x/exp v0.0.0-20240506185415-9bf2ced13842 // indirect
	golang.org/x/mod v0.18.0 // indirect
	golang.org/x/sync v0.14.0 // indirect
	golang.org/x/text v0.25.0 // indirect
	golang.org/x/time v0.6.0 // indirect
	golang.org/x/tools v0.22.0 // indirect
	golang.zx2c4.com/wintun v0.0.0-20230126152724-0fa3db229ce2 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/grpc v1.69.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace github.com/octelium/octelium/pkg => ../../pkg

replace github.com/octelium/octelium/apis => ../../apis

replace github.com/octelium/octelium/client/common => ../common

replace github.com/octelium/octelium/client/octeliumctl => ../octeliumctl

replace github.com/octelium/octelium/cluster/apis => ../../cluster/apis

replace github.com/octelium/octelium/cluster/apiserver => ../../cluster/apiserver

replace github.com/octelium/octelium/cluster/common => ../../cluster/common

replace github.com/octelium/octelium/cluster/rscserver => ../../cluster/rscserver

replace github.com/octelium/octelium/octelium-go => ../../octelium-go
