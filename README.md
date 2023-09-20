# gotool

> tool collection written in Go

## commands

### base64

> base64 encode/decode tool

**Flags:**

| global flags            | description                                            |
|:------------------------|:-------------------------------------------------------|
| -b, --base int          | base(32/64) (default 64)                               |
| -d, --data string       | input data                                             |
| --encoder string        | custom encode/decode char sequence (default: STDCodec) |
| -h, --help              | help for base64                                        |
| --hexOutput             | output hex string (upper case)                         |
| -f, --inputFile string  | output file (leave empty using --data)                 |
| -o, --outputFile string | output file (leave empty using stdout)                 |

#### decode

> base64 decoder

| flags | description                |
|:------|:---------------------------|
| --hex | hexEncoding                |
| --raw | no padding, (base 64 only) |
| --url | urlEncoding                |

#### encode

> base64 encoder

| flags      | description                 |
|:-----------|:----------------------------|
| --hex      | hexEncoding                 |
| --hexInput | process input as hex string |
| --raw      | no padding, (base 64 only)  |
| --url      | urlEncoding                 |

---

### digit

> digit converter

| flags      | description                    |
|:-----------|:-------------------------------|
| --from int | input radix (default 10)       |
| --to int   | destination radix (default 10) |

---

### hex

> hex encode/decode tool

| flags           | description                    |
|:----------------|:-------------------------------|
| --input string  | input string or file path      |
| --output string | outputFile, empty using stdout |

---

### image

> image utils

| global flags        | description |
|:--------------------|:------------|
| -i, --input string  | input image |
| -o, --output string | input image |

#### compress

> image compressor

| flags             | description              |
|:------------------|:-------------------------|
| -q, --quality int | input image (default 70) |

#### exif

> exif parser

| flags     | description   |
|:----------|:--------------|
| --json    | Print JSON    |
| --verbose | Print logging |

#### parse

> jpeg header blocks parser

| flags                | description          |
|:---------------------|:---------------------|
| -e, --extract string | extract specific tag |

#### resize

> resizer

| flags                  | description                                  |
|:-----------------------|:---------------------------------------------|
| -b, --flag             | only modify flag bytes                       |
| -f, --format string    | output format (default: jpg) (default "jpg") |
| -H, --height int       | output height                                |
| -h, --help             | help for resize                              |
| -W, --width int        | output width                                 |
| -g, --algorithm string |                                              |

algorithms:

+ Lanczos3(Default)
+ NearestNeighbor
+ Bilinear
+ Bicubic
+ MitchellNetravali
+ Lanczos2

---

### reverse

> A tool reverse bytes blocks of files
> input : <block1><block2><block3>...
> output: <1kcolb><2kcolb><3kcolb>...

| flags      | description                        |
|:-----------|:-----------------------------------|
| --base int | base32(32)/base64(64) (default 32) |
| --bit int  | block size (default 1280)          |
| --delete   | delete origin file                 |
| --enc_name | do not encode/decode file name     |
| --force    | overwrite file                     |

---

### stego

> Unicode Text Steganography Encoders/Decoders

#### advanced

> uses better looking Homoglyphs

| flags            | description                           |
|:-----------------|:--------------------------------------|
| --input string   | cover text                            |
| --message string | content to hidden(only support ASCII) |

#### mhomoglyph

> uses more Homoglyphs encode bits

| flags               | description                  |
|:--------------------|:-----------------------------|
| --input string      | input text(to decode/encode) |
| --mapping string    | custom char mapping file     |
| --message string    | message to hidden            |
| -o, --output string | output path                  |

#### tags

> ASCII map to Unicode Tags (U+E0000 to U+E007F)

| flags               | description                  |
|:--------------------|:-----------------------------|
| --all               | show all unicode code point  |
| --detail            | show mapping info            |
| --input string      | input text(to decode/encode) |
| --message string    | message to hidden            |
| -o, --output string | output path                  |
| --print_clean       | show clean text              |

---

### sort

> sort files into category

| flags                   | description |
|:------------------------|:------------|
| -D, --dest string       | sort dest   |
| -F, --files stringArray | input files |
