## `slogpfx` - Prefix Logging for `slog`

Easily prefix your log messages with attributes from the log record using `slogpfx`.

### Features

- Customizable prefix formatter.
- Colored output by default.
- Works best with [lmittmann/tint](https://github.com/lmittmann/tint)
  and [mattn/go-colorable](https://github.com/mattn/go-colorable) packages (some examples below).

### Installation

```bash
go get -u github.com/dpotapov/slogpfx
```

### Usage

Here's a quick example to get you started:

```go
h := slog.NewTextHandler(os.Stdout, nil)

// Use the prefix based on the attribute "service"
prefixed := slogpfx.NewHandler(h, &slogpfx.HandlerOptions{
    PrefixKeys:      []string{"service"}
})

logger := slog.New(prefixed)

logger.Info("Hello World!")
logger.Info("Hello World!", "service", "billing")
logger.With("service", "database").Error("Connection error")
```

Using in conjunction with [lmittmann/tint](https://github.com/lmittmann/tint)
and [mattn/go-colorable](https://github.com/mattn/go-colorable) packages:

```go
h := tint.NewHandler(colorable.NewColorable(os.Stdout), nil)

// Use the prefix based on the attribute "service"
prefixed := slogpfx.NewHandler(h, &slogpfx.HandlerOptions{
    PrefixKeys:      []string{"service"},
})

logger := slog.New(prefixed)
```

### Customization

You can customize the way prefixes are displayed using the `PrefixFormatter` option.

---

Happy Logging! 📜🖋️
