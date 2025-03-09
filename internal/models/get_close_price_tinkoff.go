package models

type InstrumentRequest struct {
	InstrumentID string `json:"instrumentId"`
}

type GetClosePricesRequest struct {
	Instruments []InstrumentRequest `json:"instruments"`
}

type Price struct {
	Units string `json:"units"`
	Nano  int `json:"nano"`
}

type ClosePrice struct {
	Figi                string `json:"figi"`
	InstrumentUid       string `json:"instrumentUid"`
	Price               Price  `json:"price"`
	EveningSessionPrice Price  `json:"eveningSessionPrice"`
	Time                string `json:"time"`
}

type ClosePricesResponse struct {
	ClosePrices []ClosePrice `json:"closePrices"`
}
