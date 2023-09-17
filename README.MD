# qrc-go

QRC decoder for QQMusic's QRC Lyrics file.

## Build

```sh
go generate ./...
go build ./cmd/qrc-decode
```

## Usage

```sh
# xml: QRC with word-sync
qrc-decode.exe input_qm.qrc output.xml
qrc-decode.exe input_qmts.qrc qrc output.lrc
```

## Golang API

Install as dependency:

```sh
go get github.com/jixunmoe-go/qrc
```

```go
import github.com/jixunmoe-go/qrc

func decoder() {
    decrypted_bytes := qrc_go.DecodeQRC([]byte("qrc file content here"))
}
```

## Note

This code is not optimized for streaming.
