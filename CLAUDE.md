# CLAUDE.md

File này hướng dẫn Claude Code (claude.ai/code) khi làm việc với code trong repo này.

## Repo này là gì

`gm-core` là **nguồn sự thật duy nhất cho contract protobuf** của workspace. Đây là module
library, không phải service: không `cmd/`, không server, không database. Mọi service `gm-*` và
`gm-gateway` import package sinh ra từ đây thay vì tự generate.

```
gm-core/proto/<service>/<version>/{dto,service}/  ──make gen-proto──>  gm-core/application/<service>/<version>/{dto,service}/
```

## Bố cục

```
proto/                                    contract của chúng ta
  common/v1/
    dto/common.dto.proto                  package common.v1.dto       — enum SystemCode (reason + systemCode)
  product/v1/
    dto/product.dto.proto                 package product.v1.dto      — mọi message
    service/product.service.proto         package product.v1.service  — chỉ khối service
third_party/                              import vendor về, không bao giờ generate code
  google/api/{annotations,http}.proto
  buf/validate/validate.proto
tools/                                    module riêng, pin buf + 3 plugin bằng tool directive
bin/                                      binary tool do make build ra, gitignore
apperror/                                 Go viết tay: sentinel lỗi nghiệp vụ dùng chung, import gopkg
application/                              code Go sinh ra, phản chiếu proto/ — file *.pb.go, *.pb.gw.go bị gitignore
  product/v1/
    dto/product.dto.pb.go                 package productdtov1
    service/product.service.pb.go         package productsvcv1
    service/product.service_grpc.pb.go
    service/product.service.pb.gw.go
```

**`apperror/` là package Go viết tay duy nhất, có chủ đích.** `enum SystemCode` là danh mục
`reason` (tên giá trị) và `systemCode` (số) của toàn hệ thống: service gửi tên trong
`ErrorInfo.reason`, gm-gateway tra số để trả `systemCode` cho client. `apperror` dựng các sentinel
dùng chung (`ErrBadRequest`, `ErrNotFound`, …) từ `responsehelper.New` của gopkg và tên trong enum,
nên gõ sai reason thì không compile. Vì vậy gm-core import gopkg (`replace ... => ../gopkg`);
gopkg trung lập và không bao giờ import gm-core. Thêm lỗi nghiệp vụ: thêm giá trị vào enum,
`make gen-proto`, thêm sentinel vào `apperror/`. Reason mà gopkg tự gắn
(`grpcmiddlewarehelper.Reason*`) phải có trong enum; `apperror_test.go` kiểm tra điều đó.

**Đường dẫn thư mục bắt buộc phải trùng proto package.** `proto/product/v1/dto/` chứa
`package product.v1.dto;`, không chứa gì khác. Đây là quy tắc style của buf, giữ cho import không
nhập nhằng khi repo này chứa hàng chục service.

**dto và service là hai proto package tách rời, không chỉ là hai file tách rời.** Hai file trong
cùng một proto package buộc phải chung `go_package`, nên dù sắp xếp thư mục kiểu gì chúng cũng rơi
vào một Go package. Chính việc tách proto package mới làm `dto` và `service` thành hai Go package
thật. Cái được: caller chỉ dựng message thì import `dto` và không bao giờ kéo theo runtime của gRPC
hay grpc-gateway.

**Mọi package đều mang hậu tố version.** Một thay đổi breaking sẽ thêm `proto/product/v2/` bên cạnh
`v1` và cả hai cùng generate, nên service phía dưới migrate theo lịch của chính nó thay vì migrate
trong đợt deploy của người khác. Thay đổi không breaking (thêm field, thêm rpc) đi vào version
đang có.

**`third_party/` là import root, không phải input.** `buf.yaml` liệt kê nó như một workspace module
để import resolve được, nhưng `gen-proto` chỉ truyền `proto` làm input, nên proto vendor về không
bao giờ sinh ra code Go.

## Đặt tên

| Thứ | Quy tắc | Ví dụ |
|---|---|---|
| dto package | `<service>.<version>.dto` | `product.v1.dto` |
| service package | `<service>.<version>.service` | `product.v1.service` |
| thư mục | trùng với package | `proto/product/v1/dto/` |
| file dto | `<service>.dto.proto` | `product.dto.proto` |
| file service | `<service>.service.proto` | `product.service.proto` |
| `go_package` của dto | `<module>/application/<service>/<version>/dto;<service>dto<version>` | `...;productdtov1` |
| `go_package` của service | `<module>/application/<service>/<version>/service;<service>svc<version>` | `...;productsvcv1` |
| khối service | `<Service>Service`, không bao giờ `API` | `service ProductService` |

File service import file dto bằng **đường dẫn đầy đủ tính từ proto root**
(`import "product/v1/dto/product.dto.proto";`) và tham chiếu message ở dạng **fully qualified**
(`product.v1.dto.HelloReply`). Resolve tương đối cũng chạy được, nhưng khi vượt qua ranh giới
package thì trong một repo contract, viết rõ ra là đáng.

## Thêm một service

1. `mkdir -p proto/<service>/v1/{dto,service}`, viết `<service>.dto.proto` và
   `<service>.service.proto` với package và `go_package` lấy từ bảng trên.
2. `make gen-proto` → `application/<service>/v1/{dto,service}/` xuất hiện.
3. `make mod && make build`.
4. Trong service tiêu thụ, thêm cặp `replace`/`require` (bên dưới) và import nửa nào nó cần:
   ```go
   <service>dtov1 "github.com/dangngochoainam/gm-core/application/<service>/v1/dto"
   <service>svcv1 "github.com/dangngochoainam/gm-core/application/<service>/v1/service"
   ```
5. Trong `gm-gateway/cmd/server/main.go`, dial service đó và gọi
   `<service>svcv1.Register<Service>ServiceHandler`.

## Codegen

```bash
make tools       # build bin/buf và 3 plugin từ tools/ — gen-proto tự gọi nếu thiếu
make gen-proto   # bin/buf generate proto -> application/
make mod         # go mod tidy
make build       # go build ./...
make fmt         # sửa format: gofmt + goimports
make lint        # golangci-lint v2.14.0 (config ở .golangci.yml) — lint toàn repo
make verify      # lint + build, không chạy test
make test        # go test ./...
```

golangci-lint không nằm trong `tools/`: nó cần go >= 1.26, còn mọi module giữ `go 1.24.0`.
Makefile chạy nó qua `go run ...@v2.14.0`, và `GOTOOLCHAIN=auto` tự tải go1.26 chỉ để build linter.

**Một bản checkout chỉ cần Go, không cần cài thêm gì.** `make gen-proto` thấy `bin/` trống thì tự
build `buf` cùng ba plugin rồi generate. Lần đầu trên một máy mới tốn khoảng một phút và cần mạng để
Go tải module về; sau đó cache lại, và những lần sau chỉ còn bước generate.

Version nằm ở **một chỗ duy nhất là `tools/go.mod`**, pin bằng tool directive của Go 1.24 và
checksum trong `tools/go.sum`. `tools/` là module riêng vì bốn tool đó kéo theo khoảng một trăm
dependency; khai báo trong `go.mod` của `gm-core` thì số đó vào module graph của mọi service, vì
service nào cũng `require gm-core`. Bump một version là sửa dòng `require` tương ứng, chạy
`make tools` rồi `make gen-proto`.

**Ba generator chạy dưới dạng local plugin, trỏ vào `bin/`.** Đây không phải protoc quay lại: buf
vẫn tự compile proto, chỉ khác là plugin chạy tại máy thay vì trên Buf Schema Registry. Lý do đổi:
remote plugin gọi registry ở **mọi** lần chạy chứ không chỉ lần đầu, nên codegen phụ thuộc vào mạng.
Cùng version thì output byte-identical, đã verify bằng `diff -r`.

Đường dẫn plugin trong `buf.gen.yaml` là tương đối so với thư mục buf chạy, tức gốc repo. Gõ tay
`buf generate` từ thư mục con sẽ không tìm thấy plugin; đi qua `make gen-proto`.

Mỗi plugin dùng `module=github.com/dangngochoainam/gm-core` thay vì `paths=source_relative`, nên
mọi file sinh ra rơi đúng đường dẫn ghi trong `go_package` của nó. Input là thư mục `proto`, nên
thư mục service mới được nhận diện mà không phải sửa config.

Dev container của các service **không** mang công cụ codegen nào — xem `NOTES.md`.

Code sinh ra **không được commit**: `.gitignore` bỏ qua `*.pb.go` và `*.pb.gw.go` trong
`application/`, chỉ những file đó, nên code viết tay đặt trong `application/` vẫn được track.
Một bản checkout mới phải chạy `make gen-proto` trước khi service nào `make mod` được;
`scripts/task-start.sh` ở root workspace làm việc đó cho worktree của task.

## Các service tiêu thụ repo này thế nào

Dùng theo path, y hệt `gopkg`:

```
// go.mod của mọi consumer
replace github.com/dangngochoainam/gm-core => ../gm-core
require github.com/dangngochoainam/gm-core v0.0.0
```

Hệ quả:

- Sau khi đổi proto ở đây, chạy `make gen-proto` ở đây, rồi `make mod` trong từng consumer.
- Consumer có vendor, nên `go mod vendor` phải chạy lại thì build của nó mới thấy thay đổi.
- Dev container bind-mount `../gm-core:/gm-core` để replace `../gm-core` resolve được từ `/app`.
- Consumer **không** có target `gen-proto`. Bước duy nhất của nó sau khi contract đổi là `make mod`.

## Nguyên tắc làm việc

Giống các repo khác: ưu tiên đơn giản, sửa đúng chỗ cần sửa, không trừu tượng hóa trước cho tương
lai. Ý tưởng hoãn lại thuộc về `NOTES.md`, không thuộc về code viết sẵn.
