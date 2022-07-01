package tlv

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type testScenarios map[*Node]struct {
	res interface{}
	ok  bool
}

func TestNode_String(t *testing.T) {
	node := Node{Raw: []byte{0x01, 0x02, 0x03, 0x04, 0x05}}

	require.Equal(t, "AQIDBAU=", node.String())
}

func TestNode_GetNodes(t *testing.T) {
	node := Node{Value: data, decoder: stdDecoder}

	nodes, err := node.GetNodes()

	require.Nil(t, err)
	require.NotNil(t, nodes)
	require.Equal(t, 1, len(nodes))
	require.Equal(t, tagMessage, nodes[0].Tag)
}

func TestNode_GetNodes_WhenTheValueIsInvalid(t *testing.T) {
	node := Node{Value: data[:3], decoder: stdDecoder}

	nodes, err := node.GetNodes()

	require.NotNil(t, err)
	require.Empty(t, nodes)
}

func TestNode_GetBool(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte{0xff, 0x12}): {true, true},   // bigger
		newNode([]byte{0xff}):       {true, true},   // exact size (>1)
		newNode([]byte{0x01}):       {true, true},   // exact size (=1)
		newNode([]byte{0x00}):       {false, true},  // exact size (=0)
		newNode([]byte{}):           {false, false}, // empty
		newNode(nil):                {false, false}, // nil
	}

	for node, expected := range scenarios {
		res, ok := node.GetBool()

		require.Equal(t, expected.res, res)
		require.Equal(t, expected.ok, ok)
	}
}

func TestNode_GetPaddedBool(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte{0xff, 0x12}): {res: true},  // bigger
		newNode([]byte{0xff}):       {res: true},  // exact size (>1)
		newNode([]byte{0x01}):       {res: true},  // exact size (=1)
		newNode([]byte{0x00}):       {res: false}, // exact size (=0)
		newNode([]byte{}):           {res: false}, // empty
		newNode(nil):                {res: false}, // nil
	}

	for node, expected := range scenarios {
		res := node.GetPaddedBool()

		require.Equal(t, expected.res, res)
	}
}

func TestNode_GetString(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte("abc")): {res: "abc"}, // valid
		newNode([]byte{}):      {res: ""},    // empty
		newNode(nil):           {res: ""},    // nil
	}

	for node, expected := range scenarios {
		res := node.GetString()

		require.Equal(t, expected.res, res)
	}
}

func TestNode_GetDate(t *testing.T) {
	date := time.Date(2021, 3, 14, 19, 26, 45, 0, time.UTC)
	node := newNode([]byte{0x60, 0x4e, 0x63, 0x75})

	res, ok := node.GetDate()

	require.True(t, ok)
	require.Equal(t, date, res)
}

func TestNode_GetDate_WhenValueIsEmpty(t *testing.T) {
	node := newNode([]byte{})

	res, ok := node.GetDate()

	require.False(t, ok)
	require.Empty(t, res)
}

func TestNode_GetUint8(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte{0x02, 0xff}): {uint8(2), true},  // bigger
		newNode([]byte{0x02}):       {uint8(2), true},  // exact size
		newNode([]byte{0x00}):       {uint8(0), true},  // exact size (zero)
		newNode([]byte{}):           {uint8(0), false}, // empty
		newNode(nil):                {uint8(0), false}, // nil
	}

	for node, expected := range scenarios {
		res, ok := node.GetUint8()

		require.Equal(t, expected.res, res)
		require.Equal(t, expected.ok, ok)
	}
}

func TestNode_GetPaddedUint8(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte{0xff, 0x12}): {res: uint8(255)}, // bigger
		newNode([]byte{0xff}):       {res: uint8(255)}, // exact size
		newNode([]byte{0x00}):       {res: uint8(0)},   // exact size (zero)
		newNode([]byte{}):           {res: uint8(0)},   // empty
		newNode(nil):                {res: uint8(0)},   // nil
	}

	for node, expected := range scenarios {
		res := node.GetPaddedUint8()

		require.Equal(t, expected.res, res)
	}
}

func TestNode_GetUint16(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte{0xab, 0xcd, 0xff}): {uint16(43981), true}, // bigger
		newNode([]byte{0xab, 0xcd}):       {uint16(43981), true}, // exact size
		newNode([]byte{0x00, 0x00}):       {uint16(0), true},     // exact size (0)
		newNode([]byte{0x00}):             {uint16(0), false},    // smaller
		newNode([]byte{}):                 {uint16(0), false},    // empty
		newNode(nil):                      {uint16(0), false},    // nil
	}

	for node, expected := range scenarios {
		res, ok := node.GetUint16()

		require.Equal(t, expected.res, res)
		require.Equal(t, expected.ok, ok)
	}
}

func TestNode_GetPaddedUint16(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte{0x00, 0xf0, 0xff}): {res: uint16(240)}, // bigger
		newNode([]byte{0x00, 0xf0}):       {res: uint16(240)}, // exact size
		newNode([]byte{0x00, 0x00}):       {res: uint16(0)},   // exact size (zero)
		newNode([]byte{0x00}):             {res: uint16(0)},   // smaller
		newNode([]byte{}):                 {res: uint16(0)},   // empty
		newNode(nil):                      {res: uint16(0)},   // nil
	}

	for node, expected := range scenarios {
		res := node.GetPaddedUint16()

		require.Equal(t, expected.res, res)
	}
}

func TestNode_GetUint32(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte{0x12, 0x34, 0x56, 0x78, 0x9a}): {uint32(305419896), true}, // bigger
		newNode([]byte{0x12, 0x34, 0x56, 0x78}):       {uint32(305419896), true}, // exact size
		newNode([]byte{0x00, 0x00, 0x12}):             {uint32(0), false},        // smaller (3 bytes)
		newNode([]byte{0x00, 0x34}):                   {uint32(0), false},        // smaller (2 bytes)
		newNode([]byte{0x56}):                         {uint32(0), false},        // smaller (1 byte)
		newNode([]byte{}):                             {uint32(0), false},        // empty
		newNode(nil):                                  {uint32(0), false},        // nil
	}

	for node, expected := range scenarios {
		res, ok := node.GetUint32()

		require.Equal(t, expected.res, res)
		require.Equal(t, expected.ok, ok)
	}
}

func TestNode_GetPaddedUint32(t *testing.T) {
	scenarios := testScenarios{
		newNode([]byte{0x12, 0x34, 0x56, 0x78, 0x9a}): {res: uint32(305419896)}, // bigger
		newNode([]byte{0x12, 0x34, 0x56, 0x78}):       {res: uint32(305419896)}, // exact size
		newNode([]byte{0x00, 0x00, 0x12}):             {res: uint32(18)},        // smaller (3 bytes)
		newNode([]byte{0x00, 0x34}):                   {res: uint32(52)},        // smaller (2 bytes)
		newNode([]byte{0x56}):                         {res: uint32(86)},        // smaller (1 byte)
		newNode([]byte{}):                             {res: uint32(0)},         // empty
		newNode(nil):                                  {res: uint32(0)},         // nil
	}

	for node, expected := range scenarios {
		res := node.GetPaddedUint32()

		require.Equal(t, expected.res, res)
	}
}

func TestNode_GetUint64(t *testing.T) {
	fullValue := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12} // bigger (9 bytes)

	scenarios := testScenarios{
		newNode(fullValue[:9]): {uint64(1311768467463790320), true}, // bigger
		newNode(fullValue[:8]): {uint64(1311768467463790320), true}, // exact size
		newNode(fullValue[:7]): {uint64(0), false},                  // smaller (7 bytes)
		newNode(fullValue[:6]): {uint64(0), false},                  // smaller (6 bytes)
		newNode(fullValue[:5]): {uint64(0), false},                  // smaller (5 bytes)
		newNode(fullValue[:4]): {uint64(0), false},                  // smaller (4 bytes)
		newNode(fullValue[:3]): {uint64(0), false},                  // smaller (3 bytes)
		newNode(fullValue[:2]): {uint64(0), false},                  // smaller (2 bytes)
		newNode(fullValue[:1]): {uint64(0), false},                  // smaller (1 byte)
		newNode([]byte{}):      {uint64(0), false},                  // empty
		newNode(nil):           {uint64(0), false},                  // nil
	}

	for node, expected := range scenarios {
		res, ok := node.GetUint64()

		require.Equal(t, expected.res, res)
		require.Equal(t, expected.ok, ok)
	}
}

func TestNode_GetPaddedUint64(t *testing.T) {
	fullValue := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0, 0x12} // bigger (9 bytes)

	scenarios := testScenarios{
		newNode(fullValue[:9]): {res: uint64(1311768467463790320)}, // bigger
		newNode(fullValue[:8]): {res: uint64(1311768467463790320)}, // exact size
		newNode(fullValue[:7]): {res: uint64(5124095576030430)},    // smaller (7 bytes)
		newNode(fullValue[:6]): {res: uint64(20015998343868)},      // smaller (6 bytes)
		newNode(fullValue[:5]): {res: uint64(78187493530)},         // smaller (5 bytes)
		newNode(fullValue[:4]): {res: uint64(305419896)},           // smaller (4 bytes)
		newNode(fullValue[:3]): {res: uint64(1193046)},             // smaller (3 bytes)
		newNode(fullValue[:2]): {res: uint64(4660)},                // smaller (2 bytes)
		newNode(fullValue[:1]): {res: uint64(18)},                  // smaller (1 byte)
		newNode([]byte{}):      {res: uint64(0)},                   // empty
		newNode(nil):           {res: uint64(0)},                   // nil
	}

	for node, expected := range scenarios {
		res := node.GetPaddedUint64()

		require.Equal(t, expected.res, res)
	}
}

func newNode(value []byte) *Node {
	return &Node{binParser: binary.BigEndian, Value: value}
}
