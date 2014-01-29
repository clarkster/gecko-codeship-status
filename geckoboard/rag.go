package geckoboard

import "encoding/json"

type RAG struct {
	RedValue int
	RedText  string

	AmberValue int
	AmberText  string

	GreenValue int
	GreenText  string
}

type ragPayload struct {
	Item []ragPayloadItem `json:"item"`
}

type ragPayloadItem struct {
	Value int    `json:"value"`
	Text  string `json:"text"`
}

func (rag *RAG) MarshalJSON() ([]byte, error) {
	red := ragPayloadItem{
		Value: rag.RedValue,
		Text:  rag.RedText,
	}
	amber := ragPayloadItem{
		Value: rag.AmberValue,
		Text:  rag.AmberText,
	}
	green := ragPayloadItem{
		Value: rag.GreenValue,
		Text:  rag.GreenText,
	}
	payload := ragPayload{
		Item: []ragPayloadItem{red, amber, green},
	}
	return json.Marshal(payload)
}
