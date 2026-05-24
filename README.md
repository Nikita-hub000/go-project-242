### Hexlet tests and linter status:
[![Actions Status](https://github.com/Nikita-hub000/go-project-242/actions/workflows/hexlet-check.yml/badge.svg)](https://github.com/Nikita-hub000/go-project-242/actions)

# hexlet-path-size

CLI-утилита для расчёта размера файла или директории.

## Установка

```bash
go build -o bin/hexlet-path-size ./cmd/hexlet-path-size
```

## Использование

```bash
./bin/hexlet-path-size [flags] <path>
```

Пример:

```bash
./bin/hexlet-path-size --recursive --all --human ./testdata
```

## Флаги

- `--human`, `-H` - вывести размер в человекочитаемом формате (`KB`, `MB`, `GB`)
- `--all`, `-a` - учитывать скрытые файлы и директории
- `--recursive`, `-r` - рекурсивно обходить вложенные директории

## Разработка

```bash
make build
make test
make lint
```

https://asciinema.org/a/GH6MHzEPRV0pstlT
