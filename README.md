# GoFetch

[![Go Reference](https://pkg.go.dev/badge/image)](https://pkg.go.dev/github.com/roarc0/gofetch)
[![Go Report](https://goreportcard.com/badge/github.com/roarc0/gofetch)](https://goreportcard.com/report/github.com/roarc0/gofetch)
[![Go Coverage](https://github.com/roarc0/gofetch/wiki/coverage.svg)](https://raw.githack.com/wiki/roarc0/gofetch/coverage.html)
![go workflow](https://github.com/roarc0/gofetch/actions/workflows/go.yml/badge.svg)

A simple magnet link scraper that filters the items based on regular expressions.

## Work in progress 🚧

This is a work in progress, and it's not yet ready for use.

## Installation  💾

```bash
go install github.com/roarc0/gofetch/cmd/gofetch-cli@latest
```

## Usage 🏄

Create the configuration file in `~/.config/gofetch/config.yaml` replace variables with actual values.

```yaml
memory:
    filepath: gofetch.db
sources:
    nyaa:
        name: nyaa
        uris:
            - https://$URL/?c=1_2&s=seeders&o=desc
entries:
    animeName:
        sourcename: nyaa
        filter:
            matchers:
                - type: regex
                  matcher:
                    regex: .*AnimeName.*
                    matchtype: required
                - type: regex
                  matcher:
                    regex: ^\[Releaser\].*
                    matchtype: required
                - type: regex
                  matcher:
                    regex: .*1080p.*
                    matchtype: required
                - type: regex
                  matcher:
                    regex: .*(480|720)p.*
                    matchtype: exclude
    animeName2:
        sourcename: nyaa
        filter:
            matchers:
                - type: regex
                  matcher:
                    regex: .*AnimeName2.*
                    matchtype: required
                - type: regex
                  matcher:
                    regex: ^\[Releaser2\].*
                    matchtype: required
                - type: regex
                  matcher:
                    regex: .*1080p.*
                    matchtype: required
                - type: regex
                  matcher:
                    regex: .*(480|720)p.*
                    matchtype: exclude
```

## Credits :star:

- [Alessandro Rosetti](https://github.com/roarc0)

## License :scroll:

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
