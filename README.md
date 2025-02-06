![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/joshua-temple/gotag)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/joshua-temple/gotag/test.yml?label=tests)
![GitHub issues](https://img.shields.io/github/issues/joshua-temple/gotag)
![GitHub stars](https://img.shields.io/github/stars/joshua-temple/gotag?style=social)


# gotag

gotag is a CLI tool written in Go that allows you to manipulate struct tags in your Go source files. It supports adding, deleting, or overwriting struct tags on specific structs, entire files, or directories. This tool ensures that all tags are uniformly ordered across fields and supports interactive mode with confirmation prompts (which can be bypassed using the force flag).

## Features

- **Targeting:**
  - **Single Struct:** Specify a struct in a file using `pkg/here/struct.go@StructName`.
  - **Entire File:** Process all structs in a file using `pkg/here/struct.go`.
  - **Directory:** Process all structs in a directory (with an optional recursive flag) using `pkg/here/`.
- **Directives:**
  - **Add (`-a`):** Add a tag if it does not exist.
  - **Delete (`-d`):** Delete a tag.
  - **Overwrite (`-o`):** Overwrite an existing tag's value (use an empty value to trigger interactive prompts/defaults).
- **Case Styles:** Specify a case style for tag values (default is `camelCase`; alternatives include `snake_case` and `kebab-case`).
- **Uniform Ordering:** Tags are stored in a uniform, sorted order across all fields.
- **Interactive Mode (`-i`):** Optionally prompt for tag values or confirmations when adding or overwriting tags.
- **Force Mode (`-f`):** Bypass confirmation prompts and force changes.

## Installation

Ensure you have Go (1.18+) installed.

Clone the repository:

```bash
git clone https://github.com/joshua-temple/gotag.git
cd gotag
```

Build the CLI tool:

```bash
go build -o gotag ./cmd
```

Or install it:

```bash
go install github.com/joshua-temple/gotag/cmd/gotag@latest
```

## Usage

### Non-interactive Mode

- **Single struct update:**

  ```bash
  gotag --target pkg/here/struct.go@MyStruct -a json -o json=new_value --case snake_case
  ```

- **Entire file update:**

  ```bash
  gotag --target pkg/here/struct.go -a xml -d json
  ```

- **Recursive directory update:**

  ```bash
  gotag --target pkg/here/ --recursive -a db
  ```

- **Multi-tag update**
  ```bash
    gotag --target pkg/here/struct.go@MyStruct -a json,xml -c kebab
  ```

### Interactive Mode

Run in interactive mode to be prompted for tag values and confirmations:

```bash
gotag --target pkg/here/struct.go@MyStruct -o json= -i
```

You will see prompts like:

```text
Enter new value for tag 'json' in field 'FieldName' (default: camelCase conversion):
```

### Force Mode

To bypass confirmation prompts (for overwrites and uniform ordering updates), use the `-f` flag:

```bash
gotag --target pkg/here/struct.go -a json -o json=new_value -f
```

## Flags

- `-t, --target` : Target file, directory, or struct (required).
- `-r, --recursive` : Recursively scan directories.
- `-i, --interactive` : Interactive mode for tag prompts/confirmations.
- `-f, --force` : Force changes and bypass confirmation prompts.
- `-c, --case` : Case style for tag values (`camelCase`, `snake_case`, `kebab-case`). Default is `camelCase`.
  - aliases are also available: `camel`, `snake`, `kebab`.
- `-a, --add` : Tag keys to add (comma-separated for multiple, e.g. `-a json,xml,db`).
- `-d, --delete` : Tag keys to delete (comma-separated for multiple, e.g. `-d json,xml`).
- `-o, --overwrite` : Tag keys to overwrite in the format `key=newValue` (use an empty value for interactive/defaults).

## Running Tests

Run all tests with:

```bash
go test ./...
```

## Contributing

[Contributions](./CONTRIBUTING.md) are welcome! Please open an issue or submit a pull request.

## License

This project is licensed under the MIT License.

## Donations

![Static Badge](https://img.shields.io/badge/XRP_(2174028412)-rw2ciyaNshpHe7bCHo4bRWq6pqqynnWKQg-blue?style=for-the-badge&logo=XRP&logoColor=white)

![Static Badge](https://img.shields.io/badge/Hedera_(3927122871)-0.0.1133968-blue?style=for-the-badge&logo=Hedera&logoColor=white)

![Static Badge](https://img.shields.io/badge/Bitcoin-bc1q7f49uuzrq2hwyclct78whaaekwpcd6r4p69mj9-blue?style=for-the-badge&logo=Bitcoin&logoColor=yellow)

![Static Badge](https://img.shields.io/badge/Polygon-0x242F47E66d725fDd5116c87a586066c5D0Cd1d3C-blue?style=for-the-badge&logo=Polygon&logoColor=pink)

![Static Badge](https://img.shields.io/badge/SUI-0x4e637f0dca15d71f3ed967ffd188e20ce9fe49f372396f251c85a1e0541d7de4-blue?style=for-the-badge&logo=Sui&logoColor=blue)

![Static Badge](https://img.shields.io/badge/Base-0x242F47E66d725fDd5116c87a586066c5D0Cd1d3C-blue?style=for-the-badge&logo=Ethereum&logoColor=white)

![Static Badge](https://img.shields.io/badge/Ethereum-0x242F47E66d725fDd5116c87a586066c5D0Cd1d3C-blue?style=for-the-badge&logo=Ethereum&logoColor=white)

![Static Badge](https://img.shields.io/badge/Solana-FuTked8fnRAiFsYNjQPsnuPQBcdQd2QZ2K822LvjYxvu-blue?style=for-the-badge&logo=Solana&logoColor=teal)



