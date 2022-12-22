# go-dalle2 #

![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/Andrew-peng/go-dalle2)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/Andrew-peng/go-dalle2/dalle2)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Andrew-peng/go-dalle2/lint.yaml?label=lint)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/Andrew-peng/go-dalle2/test.yaml?label=test)
[![codecov](https://codecov.io/gh/Andrew-peng/go-dalle2/branch/master/graph/badge.svg?token=RKA5UFHKS0)](https://codecov.io/gh/Andrew-peng/go-dalle2)
![GitHub](https://img.shields.io/github/license/Andrew-peng/go-dalle2?color=blue)

Unofficial Dalle-2 API golang client library

## Setup ##

go-dalle2 is a simple client library, so with Go installed:

```bash
go get github.com/Andrew-peng/go-dalle2/dalle2
```

or

```go
import "github.com/Andrew-peng/go-dalle2/dalle2"
```

## Usage ##

For details about DALLE-2 (and api) visit the OpenAI documentation.

### Instantiate client ###

Instantiate the client with a valid OpenAI API key

```go
client, err := dalle2.MakeNewClientV1(apiKey)
if err != nil {
    log.Fatalf("Error initializing client: %s", err)
}
```

### Creating images ###

```go
resp, err := client.Create(
    context.Background(),
    "A skyline view of New York during the sunset, watercolor",
    dalle2.WithNumImages(1),
    dalle2.WithSize(dalle2.SMALL),
    dalle2.WithFormat(dalle2.URL),
)
if err != nil {
    log.Fatal(err)
}
for _, img := range resp.Data {
    fmt.Println("%s", img.Url)
}
```

| Prompt | Output |
| --- | --- |
| A skyline view of New York during the sunset, watercolor | ![A skyline view of New York during the sunset, watercolor](examples/image/output.png) |

### Editing images ###

```go
imgBytes := ...
maskBytes := ...
resp, err := client.Edit(
    context.Background(),
    imgBytes,
    maskBytes,
    "A cute baby sea otter wearing a large sombrero",
    dalle2.WithNumImages(1),
    dalle2.WithSize(dalle2.SMALL),
    dalle2.WithFormat(dalle2.URL),
)
if err != nil {
    log.Fatal(err)
}
for _, img := range resp.Data {
    fmt.Println("%s", img.Url)
}
```

| Prompt | Image | Mask | Output |
| --- | --- | --- | --- |
| A cute baby sea otter wearing a large sombrero | ![otter](examples/edit/otter.png) | ![mask](examples/edit/mask.png) | ![A cute baby sea otter wearing a large sombrero](examples/edit/output.png) |

### Creating variations ###

```go
imgBytes := ...
resp, err := client.Edit(
    context.Background(),
    imgBytes,
    dalle2.WithNumImages(1),
    dalle2.WithSize(dalle2.SMALL),
    dalle2.WithFormat(dalle2.URL),
)
if err != nil {
    log.Fatal(err)
}
for _, img := range resp.Data {
    fmt.Println("%s", img.Url)
}
```

| Image | Output |
| --- | --- |
| ![otter](examples/variations/otter.png) | ![output variation](examples/variations/output.png) |

## License ##

go-dalle2 is distributed under the MIT-style license found in the [LICENSE](./LICENSE)
file.
