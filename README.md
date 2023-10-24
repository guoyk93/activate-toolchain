# activate-toolchain

A toolchain download and activation tool with various mirror sites embedded

## Installation

You can either build binary from source, or just download pre-built binary.

* Build from source

    ```shell
   git clone https://github.com/guoyk93/activate-toolchain.git
   cd activate-toolchain
   go build -o activate-toolchain ./cmd/activate-toolchain
    ```

* Download pre-built binaries

  View <https://github.com/guoyk93/activate-toolchain/releases>

## Usage

Only support shell in `POSIX` environment.

```shell
eval "$(activate-toolchain node@16.2 jdk@17 maven@3.8)"
```

Or create a `toolchains.txt` file with each line a toolchain, and run

```shell
eval "$(activate-toolchain)"
```

## Supported Toolchains and Version Examples

| Toolchain | Version Examples                      |
|-----------|---------------------------------------|
| node      | `node@16`, `node@16.2`, `node@16.2.0` |
| jdk       | `jdk@8`, `jdk@8.0`, `jdk@8.0.372`     |
| maven     | `maven@3`, `maven@3.8`, `maven@3.8.1` |
| ossutil   | `ossutil@1`                           |
| pnpm      | `pnpm@8`, `pnpm@8.9`, `pnpm@8.9.2`    |

## Credits

GUO YANKE, MIT License
