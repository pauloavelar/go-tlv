package tlv

import (
	"fmt"
	"time"
)

func ExampleDecodeSingle() {
	message, _, err := DecodeSingle(data)
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
	// {Title:Hello there! Silent:false ActionID:12345678 Timestamp:2021-06-30 15:34:56 +0000 UTC}
	// {Title:You there? Silent:true ActionID:240 Timestamp:2021-03-01 04:02:03 +0000 UTC}
}

const (
	tagMessage          Tag = 0x0001
	tagPushNotification Tag = 0x0101
	tagTitle            Tag = 0x0102
	tagActionID         Tag = 0x0103
	tagTimestamp        Tag = 0x0104
	tagSilent           Tag = 0x0105
)

type pushNotification struct {
	Title     string    `json:"title"`
	Silent    bool      `json:"silent"`
	ActionID  uint64    `json:"action_id"`
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
	if actionId, ok := nodes.GetFirstByTag(tagActionID); ok {
		pn.ActionID = actionId.GetPaddedUint64()
	}
	if timestamp, ok := nodes.GetFirstByTag(tagTimestamp); ok {
		ts, _ := timestamp.GetDate()
		pn.Timestamp = ts
	}

	return pn
}
