# qed

`qed` 是一个基于 Go 的命令行工具，旨在提供高效、可靠的 `加密`, `解密`功能.

## Commands

```shell
cat example.json | qed -e password > example.json.enc
cat example.json.enc | qed -d password > example.json
```

其中 `enc` 为二进制文件，格式如下:

```
MN(4) + VERSION(4) + METHOD(32) + BODY(N) + SIGNATURE(sha256) + RE-SIGNATURE(sha256)
```

其中 `SIGNATURE` 为 `sha256(MN+VERSION+METHOD+BODY)`
`RE-SIGNATURE` 为 `sha256(MN+VERSION+METHOD+BODY+SIGNATURE+INTERNAL-CODE)`

`INTERNAL-CODE` 为 `qed` 内置的数据，二次签名用于校验工具链。

## 架构

- 使用 cobra.
