package models

import "time"

type HistoricCandle struct {
	InstrumentId string    `json:"instrumentID" gorm:"primaryKey;size:255"`
	Time         time.Time `json:"time" gorm:"primaryKey"`
	Volume       string    `json:"volume"`

	High  High  `json:"high" gorm:"foreignKey:InstrumentId,Time;references:InstrumentId,Time"`
	Low   Low   `json:"low" gorm:"foreignKey:InstrumentId,Time;references:InstrumentId,Time"`
	Close Close `json:"close" gorm:"foreignKey:InstrumentId,Time;references:InstrumentId,Time"`
	Open  Open  `json:"open" gorm:"foreignKey:InstrumentId,Time;references:InstrumentId,Time"`
}

type High struct {
	InstrumentId string    `json:"-" gorm:"primaryKey;size:255"`
	Time         time.Time `json:"-" gorm:"primaryKey"`
	Units        string    `json:"units"`
	Nano         int       `json:"nano"`
}

type Low struct {
	InstrumentId string    `json:"-" gorm:"primaryKey;size:255"`
	Time         time.Time `json:"-" gorm:"primaryKey"`
	Units        string    `json:"units"`
	Nano         int       `json:"nano"`
}

type Close struct {
	InstrumentId string    `json:"-" gorm:"primaryKey;size:255"`
	Time         time.Time `json:"-" gorm:"primaryKey"`
	Units        string    `json:"units"`
	Nano         int       `json:"nano"`
}

type Open struct {
	InstrumentId string    `json:"-" gorm:"primaryKey;size:255"`
	Time         time.Time `json:"-" gorm:"primaryKey"`
	Units        string    `json:"units"`
	Nano         int       `json:"nano"`
}

type HistoricCandles struct {
	Candles []HistoricCandle `json:"candles"`
}
