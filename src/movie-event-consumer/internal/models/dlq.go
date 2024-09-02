package models

type Dlq struct {
	PK            string `dynamodbav:"PK"`
	SK            string `dynamodbav:"SK"`
	GSI_PK        string `dynamodbav:"GSI_PK"`
	GSI_SK        string `dynamodbav:"GSI_SK"`
	CorrelationId string `dynamodbav:"CorrelationId"`
	EventName     string `dynamodbav:"EventName"`
	Payload       string `dynamodbav:"Payload"`
}

func NewDlq(correlationId string, eventName string, payload string) *Dlq {
	return &Dlq{
		PK:            "DLQ#" + eventName,
		SK:            correlationId,
		GSI_PK:        "DLQ",
		GSI_SK:        correlationId,
		CorrelationId: correlationId,
		EventName:     eventName,
		Payload:       payload,
	}
}
