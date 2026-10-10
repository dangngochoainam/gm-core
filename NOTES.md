# Những quyết định đang hoãn

Những thứ cố tình để ngoài phiên bản đầu của `gm-core`. Mỗi mục có lý do riêng; chỉ nhặt lên khi
xuất hiện nhu cầu cụ thể.

## Package `common/v1/dto` dùng chung

`KeysetPagination` và response envelope `trace_id / status_code / reason_code / reason_message /
time_ms` hiện đang bị lặp bên trong `product/v1/dto/product.dto.proto`. Đây là ứng viên rõ ràng cho
`proto/common/v1/dto/common.dto.proto` để mọi service import.

Chưa làm vì: mới đúng một service migrate sang thì chẳng có gì để chia sẻ, và một package dùng
chung buộc phải chốt xem envelope là một message để embed hay là các field lặp lại. Xem lại khi
service thứ hai xuất hiện và chỗ lặp trở thành thật.

## `buf lint` và `buf breaking`

`buf.yaml` mới khai báo workspace, ngoài ra không gì khác. Hai bước kiểm tra làm nên giá trị của
một repo contract thì vẫn còn thiếu:

- `buf lint` sẽ ép bảng đặt tên trong `CLAUDE.md` bằng công cụ thay vì bằng review. Bố cục hiện tại
  đã theo đúng category `STANDARD` của buf (thư mục trùng package, package có version), nên bật lên
  là pass ngay lần đầu.
- `buf breaking --against '.git#branch=main'` sẽ bắt được một field bị đánh số lại hay một rpc bị
  xóa trước khi nó tới tay consumer. Đây mới là thứ thực sự bảo vệ quy tắc `v1`/`v2`: hiện chưa có
  gì ngăn một sửa đổi breaking ngay bên trong `v1`.

Chưa làm: cả hai chỉ hữu ích khi gắn vào CI, mà ở đây chưa có CI. Thêm vào là một target `lint:` và
một target `breaking:` cộng vài khối config, khoảng mười dòng.

## Code sinh ra không commit (đã chốt, không phải hoãn)

`*.pb.go` và `*.pb.gw.go` trong `application/` bị gitignore; chỉ những file đó, không phải cả thư
mục. Lý do: bản `.pb.go` từng commit đã lệch khỏi `.proto` (thiếu `trace_id` và các field envelope),
làm gm-product không build được trên `master`. Không commit thì không có gì để lệch, nên cũng không
cần job CI chạy `gen-proto` rồi `git diff --exit-code application/`.

Giá phải trả: checkout nào cũng phải `make gen-proto` trước khi service `make mod`
(`scripts/task-start.sh` lo cho worktree), CI sau này cũng phải có bước đó, và xem mục "Publish
`gm-core` theo version".

## Codegen được phép chạy ở đâu (đã chốt, không phải hoãn)

`buf` cố tình **không** được cài vào dev image của service, và `gm-product` cùng `gm-gateway` cố ý
**không** có target `gen-proto`. Trước đây cả hai đều ủy quyền về đây bằng
`$(MAKE) -C ../gm-core gen-proto`, và cả hai image đều mang protoc cùng bốn plugin `protoc-gen-*`
còn sót từ template `codebase-go-grpc` cũ. Tất cả đã bị gỡ.

Quy tắc giờ được ép bằng công cụ chứ không bằng thỏa thuận miệng: generate diễn ra trong repo này,
từ shell của host, và không ở đâu khác. Bước duy nhất của một service là `make mod` sau đó.

Nếu sau này ai đó muốn lấy lại lối tắt, hãy thêm `buf` vào image chứ đừng thêm các plugin. Đưa
protoc trở lại sẽ tái tạo đúng tình trạng lệch version mà việc chuyển sang `buf` đã dẹp được.

## Local plugin thay cho remote plugin (đã chốt, không phải hoãn)

`buf.gen.yaml` từng dùng remote plugin trên Buf Schema Registry, và ghi chú cũ ở đây nói rằng chỉ
lần chạy đầu tiên của mỗi version plugin cần mạng. **Điều đó sai.** buf cache module và descriptor,
không cache kết quả gọi plugin: mỗi `buf generate` là một lần gửi `CodeGeneratorRequest` lên
registry. Chặn egress rồi chạy thì hỏng ngay lần nào cũng vậy, với
`Failure: the server hosted at that remote is unavailable.`

Nay ba generator chạy dưới dạng local plugin trỏ vào `bin/`, build từ tool directive trong
`tools/go.mod`. Đây không phải protoc quay lại: buf vẫn tự compile proto, chỉ khác là plugin chạy
tại máy. Cùng version thì output không khác một byte, verify bằng `diff -r` giữa `application/` đang
commit và một lần generate bằng plugin local.

Đánh đổi: version giờ nằm trong `tools/go.mod` chứ không nằm trong `buf.gen.yaml`, và
`buf generate` gõ tay từ thư mục con sẽ không tìm thấy plugin vì đường dẫn là tương đối so với gốc
repo. Đi qua `make gen-proto`.

## Commit `tools/vendor` cho môi trường không có egress

Một bản checkout hiện chỉ cần Go, nhưng lần `make gen-proto` đầu tiên trên mỗi máy vẫn cần mạng để
Go tải source của bốn tool về; sau đó module cache lo phần còn lại. Đo được: build lạnh bốn binary
khoảng 54 giây, build lại khi đã có cache khoảng nửa giây.

Muốn một máy chưa từng thấy internet cũng generate được thì chạy `go mod vendor` trong `tools/` và
commit cây vendor đó. Đã thử: với `GOPROXY=off` và egress bị chặn, build từ vendor rồi generate ra
đúng output.

Chưa làm vì giá là 62 MB và 4089 file Go trong repo, cộng việc mỗi lần bump một tool là một diff
kéo theo cả cây vendor. Repo chưa có yêu cầu air-gapped nào. Nhặt lên khi CI chạy trong môi trường
không có egress.

## Publish `gm-core` theo version

Đang dùng theo path (`replace => ../gm-core`) giống `gopkg`, nên một thay đổi contract là vô hình
với service cho tới khi nó chạy `make mod`. Tag version Go module thật sẽ cho mỗi service một phiên
bản contract tường minh, nhưng điều đó chỉ có ý nghĩa khi các service deploy độc lập. Lưu ý đây là
trục khác với hậu tố proto `v1`, vốn để version hóa wire contract.

Vướng: code sinh ra không được commit, mà `go get` không chạy codegen, nên một tag của repo này
sẽ thiếu package Go. Muốn publish thì phải đưa code sinh ra vào artifact được tag, ví dụ một nhánh
hay repo release do CI generate rồi commit.

## Sinh OpenAPI

Chưa có `protoc-gen-openapiv2` ở đâu cả: không có trong `tools/go.mod`, cũng không nối vào
`gen-proto`. Thêm khi thực sự có người cần output swagger. Việc cần làm là thêm một dòng `tool` vào
`tools/go.mod`, một biến `PKG_` cùng binary tương ứng trong `Makefile`, và một khối `local:` trong
`buf.gen.yaml`.

## Chính sách deprecate ở mức field

Chưa có gì quy định một field được khai tử ra sao: `[deprecated = true]` trong một release, rồi
`reserved`. Viết ra lần đầu tiên có một field thật sự phải ra đi.
