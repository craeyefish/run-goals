package models

type StravaWebhookPayload struct {
	AspectType     string                 `json:"aspect_type"`
	EventTime      int64                  `json:"event_time"`
	ObjectType     string                 `json:"object_type"`
	ObjectID       int64                  `json:"object_id"`
	OwnerID        int64                  `json:"owner_id"`
	SubscriptionID int64                  `json:"subscription_id"`
	Updates        map[string]interface{} `json:"updates"`
}
