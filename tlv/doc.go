/*
Package tlv holds all logic to decode TLV messages into [Nodes].

There are two ways to use this package: with the standard decoder or creating
a custom [Decoder] with custom tag/length sizes and byte order. The standard
decoder uses 2-byte tags, 2-byte lengths and big endian [binary.ByteOrder].

[Nodes] are a representation of a collection of decoded TLV messages.
Specific messages can be filtered by [Tag] and indexes can be accessed
directly in an array-like syntax:

	firstNode := nodes[0]

[Node] is a representation of a single TLV message, comprised of a [Tag],
a [Length] and a Value. The struct has many helper methods to parse the
value as different types (e.g. integers, strings, dates and booleans),
as well as nested TLV [Nodes].

	node.GetNodes()
	node.GetPaddedUint8()

Note: [Nodes] decoded with a custom configuration retain the configuration
when parsing their values as other nodes, so messages always have consistent
tag/length sizes and byte order.
*/
package tlv
