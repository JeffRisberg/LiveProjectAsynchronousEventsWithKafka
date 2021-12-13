package config

const (
	// OrderReceivedTopicName name of the topic that handles OrderReceived events
	OrderReceivedTopicName = "OrderReceived"

	// OrderConfirmedTopicName name of the topic that handles OrderCofirmed events
	OrderConfirmedTopicName = "OrderConfirmed"

	// ErrorsTopicName name of the topic that handles Error events
	ErrorsTopicName = "DeadLetterQueue"
)
