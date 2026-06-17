# qed

`qed` 是一个基于 Go 的命令行工具，旨在提供高效、可靠的 `加密`, `解密`功能.

## 安装

> **一键安装**（复制即用）

```shell
curl -fsSL https://raw.githubusercontent.com/CDU-Ge/qed/main/install.sh | sh
```

或使用 wget:

```shell
wget -qO- https://raw.githubusercontent.com/CDU-Ge/qed/main/install.sh | sh
```

默认安装到 `~/.local/bin/qed`，确保该目录在 PATH 中:

```shell
export PATH="$HOME/.local/bin:$PATH"
```

### 安装选项

```shell
curl -fsSL https://raw.githubusercontent.com/CDU-Ge/qed/main/install.sh | sh -s -- --help
```

| 选项                  | 说明                      |
|---------------------|-------------------------|
| `--version <tag>`   | 指定版本标签，例如 v0.1.0        |
| `--dir <path>`      | 安装目录，默认 ~/.local/bin    |
| `--repo <repo>`     | GitHub 仓库，默认 CDU-Ge/qed |
| `--progress <mode>` | 进度模式: auto, bar, quiet  |
| `--quiet`           | 静默安装                    |

### 环境变量

| 变量                | 说明             |
|-------------------|----------------|
| `QED_VERSION`     | 同 `--version`  |
| `QED_INSTALL_DIR` | 同 `--dir`      |
| `QED_REPO`        | 同 `--repo`     |
| `QED_BINARY_NAME` | 安装的命令名称，默认 qed |
| `QED_PROGRESS`    | 同 `--progress` |

### 示例

安装指定版本:

```shell
curl -fsSL https://raw.githubusercontent.com/CDU-Ge/qed/main/install.sh | sh -s -- --version v0.1.0
```

安装到自定义目录:

```shell
curl -fsSL https://raw.githubusercontent.com/CDU-Ge/qed/main/install.sh | sh -s -- --dir /usr/local/bin
```

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

## 示例

```shell
cat examples/example.json | qed -e 123456 > examples/example.json.enc
cat examples/example.json.enc | qed -d 123456
```