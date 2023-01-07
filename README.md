# Go TLV

[![go version](https://img.shields.io/github/go-mod/go-version/pauloavelar/go-tlv)](https://github.com/pauloavelar/go-tlv/blob/main/go.mod)
[![go docs](https://pkg.go.dev/badge/github.com/pauloavelar/go-tlv.svg)](https://pkg.go.dev/github.com/pauloavelar/go-tlv)
[![license](https://img.shields.io/github/license/pauloavelar/go-tlv)](https://github.com/pauloavelar/go-tlv/blob/main/LICENSE)
[![build](https://img.shields.io/github/actions/workflow/status/pauloavelar/go-tlv/ci.yml?branch=main)](https://github.com/pauloavelar/go-tlv/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/pauloavelar/go-tlv/branch/main/graph/badge.svg?token=4V15TQTKRR)](https://codecov.io/gh/pauloavelar/go-tlv)
[![open issues](https://img.shields.io/github/issues-raw/pauloavelar/go-tlv)](https://github.com/pauloavelar/go-tlv/issues)

## What is TLV?

**Tag-Length-Value (TLV)** is a binary encoding scheme used for data transport.

### Main advantages

* Very flexible, easy to extend and change as needed
* Messages can be easily decoded/displayed as a hierarchical tree-like structure
* New tags can be added/moved without breaking decoder compatibility
* Searching for specific tags in long payloads is easy and efficient

### Format

There are many implementations of the scheme. The example below is one of them:

```
00 0f   # 2 bytes reserved for the tag ID
00 04   # 2 bytes reserved for the length
01 02   # As many bytes as informed by the
03 04   # length value
```

The **value** itself can be any binary format, such as numerical representations, strings and even
other TLV messages. See [data_test.go](https://github.com/pauloavelar/go-tlv/blob/main/tlv/data_test.go)
for an example of a complex structure.

> It is up to the parser to know the **value** type and format based on the **tag**.

## Features

### Byte array decoding as multiple TLV nodes

```go
data := []byte{0x00, 0x01, 0x02 /* ... */}

nodes, err := tlv.DecodeBytes(data)
if err != nil {
    panic(err) // invalid payload length vs bytes available
}

nodes.HasTag(0x0123)        // returns a bool with the tag presence
nodes.GetByTag(0x0f2a)      // returns a filtered Nodes structure
nodes.GetFirstByTag(0xabcd) // returns a Node structure with value accessors 
```

### Byte array decoding as a single TLV node

```go
data := []byte{0x00, 0x01, 0x02 /* ... */}

n, err := tlv.DecodeSingle(data)
if err != nil {
    panic(err) // invalid payload length vs bytes available
}

n.String()          // returns a base64 representation of the raw message
n.GetNodes()        // parses the value as TLV and returns a Nodes structure (or error)
n.GetUint8()        // parses the value as uint8 (returns error if value is too small)
n.GetPaddedUint8()  // parses the value as uint8 and pads it if too small

// all available types: bool, uint8, uint16, uint32, uint64, string, time.Time and Nodes
```

### Custom Decoder with different sizes and endianness

The public functions exposed in the `tlv` package use a **standard decoder** with tags and
lengths always taking 2 bytes per node, and the bytes are parsed using `binary.BigEndian`
as the `ByteOrder`.

In order to decode messages with different configuration, there is Decoder constructor:

```go
decoder, err := tlv.CreateDecoder(4, 4, binary.LittleEndian)
```

> The constructor validates the tag and length sizes, as they must be between `1` and `8`.

### Supported types

| Type     | Max Length (bytes) | Notes                                                             |
|----------|-------------------:|-------------------------------------------------------------------|
| `bool`   |                  1 | Any **non-zero** value is treated as `true`                       | 
| `uint8`  |                  1 |                                                                   |
| `uint16` |                  2 |                                                                   |
| `uint32` |                  4 |                                                                   |
| `uint64` |                  8 |                                                                   |
| `Time`   |                  8 | Value is parsed as padded `uint64` and then as **Unix** (seconds) |
| `string` |      **Unlimited** | Value is parsed as **UTF-8**                                      |
| `Nodes`  |      **Unlimited** |                                                                   |

> If the **value** is bigger than the **max length**, only the first _n_ bytes are used.

## Important details

### Tags are non-unique in TLV messages

When parsing a value to multiple nodes, tags can be **repeated** and will be returned by the decoder.
Use `Nodes#GetByTag(tlv.Tag)` and `Nodes#GetFirstByTag(tlv.Tag)` to fetch **all** or **one** node,
respectively.

#### Example:

```yaml
# Visual representation of a repeated tag in an object-like payload
message:
  - object:
      - repeated_tag: a  # this will be a node 
      - repeated_tag: b  # this will be another node
```

### The decoder supports multiple root level messages

After reading a TLV-encoded message from a byte-array, when using `tlv.DecodeBytes([]byte)` the parser
will continue reading the array until it reaches the end. The returned structure will have **all the
nodes** found in the payload.

> ⚠️&nbsp; The decoder works in an all or none strategy when dealing with multiple messages.

### Manually-created nodes use the default decoder configuration

When a `tlv.Node` is created by declaring the struct, all methods that require context, such as `GetNodes`
or `GetUint8` (or any other integer parser), will use the **standard decoder** definitions. See above for
more details on the decoder. To create a node with custom decoder configuration, first create a decoder
and call the `NewNode` method on it.


```go
var node tlv.Node

node = tlv.Node{Tag: Tag(0x1234), Value: []byte{1}}
node.GetNodes() // uses the standard decoder configuration

customDecoder := tlv.MustCreateDecoder(1, 1, binary.LittleEndian)
node = customDecoder.NewNode(Tag(0x1234), []byte{1})
node.GetNodes() // uses the customDecoder configuration
```

## Caveats

### No bit parity or checksum

The encoding scheme itself does *not* provide **bit parity** or **checksum** to ensure the integrity
of received payloads. It is up to the upper layer or to the payload design to add these features.

### Errors with multiple messages are hard to pinpoint

The bigger the payload, more likely errors will *not* be identified by the parser. The **only**
failproof hint of a malformed payload is a mismatch between the read length and the remaining bytes
in the stream. When that happens, a reading error may have happened *anywhere* in the payload, which
means none of it can be trusted.

> ⚠️&nbsp; If by the end of the stream there is a mismatch between the **provided length** and the
> **remaining bytes**, the whole payload is invalidated, and the decoder will return an error,
> **regardless of how many successful messages it has read**.

## Changelog

* **`v1.1.0`** (2023-06-01)
  * [#23](https://github.com/pauloavelar/go-tlv/pull/23): nil pointer errors on manually-created nodes
    * fix panics when calling value getters on a node without a decoder reference
    * provide functions to create a Node with the proper configuration (standard or custom)

* **`v1.0.0`** (2022-07-01)
  * **Breaking** change: parser has been renamed to decoder
  * [#10](https://github.com/pauloavelar/go-tlv/issues/10): add support to custom tag and length sizes
  * [#11](https://github.com/pauloavelar/go-tlv/issues/11): add support to custom endianness (byte order)

* **`v1.0.0-alpha1`** (2021-03-14)
  * First release with basic parsing support
  * ⚠️&nbsp; Methods and structs may change completely 
