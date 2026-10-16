package eticket

import (
	"encoding/json"
	"ticket-backend/app/dto"
)

/////////// Request

type SearchDto struct {
	Directions DirectionsDto `json:"directions"`
	RouteType  string        `json:"routeType"`
}

type DirectionsDto struct {
	Forward Direction `json:"forward"`
}

type Direction struct {
	Date           string `json:"date"`
	DepStationCode string `json:"depStationCode"`
	ArvStationCode string `json:"arvStationCode"`
}

///////////// Response

type SearchResponseDto struct {
	Data  *SearchResponseDataDto `json:"data"`
	Error *json.RawMessage       `json:"error"`
}

type SearchResponseDataDto struct {
	Directions ResponseDirectionsDto `json:"directions"`
}

type ResponseDirectionsDto struct {
	Forward *ResponseForwardDto `json:"forward"`
}

type ResponseForwardDto struct {
	Trains []Train `json:"trains"`
}

type Train struct {
	Type          string      `json:"type"`
	Number        string      `json:"number"`
	DepartureDate string      `json:"departureDate"`
	ArrivalDate   string      `json:"arrivalDate"`
	TimeOnWay     string      `json:"timeOnWay"`
	Brand         string      `json:"brand"`
	ProviderType  string      `json:"providerType"`
	OriginRoute   OriginRoute `json:"originRoute"`
	SubRoute      SubRoute    `json:"subRoute"`
	Cars          []Car       `json:"cars"`
}

type OriginRoute struct {
	DepStationName string `json:"depStationName"`
	ArvStationName string `json:"arvStationName"`
}

type SubRoute struct {
	DepStationName string `json:"depStationName"`
	DepStationCode string `json:"depStationCode"`
	ArvStationName string `json:"arvStationName"`
	ArvStationCode string `json:"arvStationCode"`
}

type Car struct {
	Type       string     `json:"type"`
	FreeSeats  int        `json:"freeSeats"`
	Tariffs    []Tariff   `json:"tariffs"`
	SeatDetail SeatDetail `json:"seatDetail"`
}

func (c Car) WagonType() (dto.WagonType, bool) {
	switch c.Type {
	case "Plaskartli":
		return dto.WagonTypePlatzkart, true
	case "Kupe":
		return dto.WagonTypeCoupe, true
	default:
		return "", false
	}
}

type Tariff struct {
	ClassServiceType string `json:"classServiceType"`
	FreeSeats        int    `json:"freeSeats"`
	Tariff           int64  `json:"tariff"`
}

type SeatDetail struct {
	Undef     int `json:"undef"`
	LateralDn int `json:"lateralDn"`
	LateralUp int `json:"lateralUp"`
	FreeComp  int `json:"freeComp"`
	Down      int `json:"down"`
	Up        int `json:"up"`
}
