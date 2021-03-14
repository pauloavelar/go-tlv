package tlv

import (
	"fmt"
	"time"
)

func ExampleParseBytes() {
	message, _, err := ParseSingle(data)
	if err != nil {
		panic(err)
	}

	if message.Tag != tagMessage {
		panic("received message is unknown")
	}

	items, err := message.GetNodes()
	if err != nil {
		panic(err)
	}

	var pushNotifications []pushNotification
	for _, pn := range items.GetByTag(tagPushNotification) {
		nodes, err := pn.GetNodes()
		if err != nil {
			panic(err)
		}

		pushNotifications = append(pushNotifications, newPushNotification(nodes))
	}

	fmt.Println(len(pushNotifications))
	if len(pushNotifications) == 2 {
		fmt.Printf("%+v\n", pushNotifications[0])
		fmt.Printf("%+v\n", pushNotifications[1])
	}

	// Output: 2
	// {Title:Hello there! Silent:false ActionId:12345678 Timestamp:2021-06-30 15:34:56 +0000 UTC}
	// {Title:You there? Silent:true ActionId:240 Timestamp:2021-03-01 04:02:03 +0000 UTC}
}

const (
	tagMessage          Tag = 0x0001
	tagPushNotification Tag = 0x0101
	tagTitle            Tag = 0x0102
	tagActionid         Tag = 0x0103
	tagTimestamp        Tag = 0x0104
	tagSilent           Tag = 0x0105
)

type pushNotification struct {
	Title     string    `json:"title"`
	Silent    bool      `json:"silent"`
	ActionId  uint64    `json:"action_id"`
	Timestamp time.Time `json:"timestamp"`
}

func newPushNotification(nodes Nodes) pushNotification {
	var pn pushNotification

	if title, ok := nodes.GetFirstByTag(tagTitle); ok {
		pn.Title = title.GetString()
	}
	if silent, ok := nodes.GetFirstByTag(tagSilent); ok {
		pn.Silent = silent.GetPaddedBool()
	}
	if actionId, ok := nodes.GetFirstByTag(tagActionid); ok {
		pn.ActionId = actionId.GetPaddedUint64()
	}
	if timestamp, ok := nodes.GetFirstByTag(tagTimestamp); ok {
		ts, _ := timestamp.GetDate()
		pn.Timestamp = ts.UTC()
	}

	return pn
}

/*
 * message:
 *   - push_notification:
 *       title: Hello world!
 *       action_id: 12345678
 *       timestamp: 2021-06-30T12:34:56Z
 *   - push_notification:
 *       title: You there?
 *       action_id: 240
 *       silent: true
 *       timestamp: 2021-03-01T01:02:03Z
 *
 */
var data = []byte{
	0x00, 0x01, // Tag: message
	0x00, 0x47, // Length: 71 bytes

	0x01, 0x01, // Tag: push_notification
	0x00, 0x1f, // Length: 31 bytes
	0x01, 0x02, // Tag: title
	0x00, 0x0c, // Length: 12 bytes
	0x48, 0x65, // Value: Hello there!
	0x6c, 0x6c,
	0x6f, 0x20,
	0x74, 0x68,
	0x65, 0x72,
	0x65, 0x21,
	0x01, 0x03, // Tag: action_id
	0x00, 0x03, // Length: 3 bytes
	0xbc, 0x61, // Value: 12345678
	0x4e,
	0x01, 0x04, // Tag: timestamp
	0x00, 0x04, // Length: 4 bytes
	0x60, 0xdc, // Value: 1625067296 = 2021-06-30T12:34:56Z
	0x8f, 0x20,

	0x01, 0x01, // Tag: push_notification
	0x00, 0x20, // Length: 32 bytes
	0x01, 0x02, // Tag: title
	0x00, 0x0a, // Length: 10
	0x59, 0x6f, // Value: You there?
	0x75, 0x20,
	0x74, 0x68,
	0x65, 0x72,
	0x65, 0x3f,
	0x01, 0x03, // Tag: action_id
	0x00, 0x01, // Length: 1 byte
	0xf0,       // Value: 240
	0x01, 0x05, // Tag: silent
	0x00, 0x01, // Length: 1 byte
	0x01,       // Value: true
	0x01, 0x04, // Tag: timestamp
	0x00, 0x04, // Length: 4 bytes
	0x60, 0x3c, // Value: 1614571323 = 2021-03-01T01:02:03Z
	0x67, 0x3b,
}
