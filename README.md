# go-csv2json

## Prerequisites

Ensure the following dependencies are installed:

- [Go](https://go.dev/doc/install)
- [pre-commit](https://pre-commit.com/) (`pip install pre-commit`)

### Run Setup

Execute the following command to set up the project:

```sh
make setup
```

#### Build

```shell
make build
```

#### Run

```shell
set -o allexport
source .env
set +o allexport
./csv2json
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
