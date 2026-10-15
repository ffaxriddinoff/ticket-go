package dto

type Position string

const (
	PositionAll   Position = "all"
	PositionUpper Position = "upper"
	PositionLower Position = "lower"
)

type Zone string

const (
	ZoneAll    Zone = "all"
	ZoneNoSide Zone = "no_side"
)

type WagonType string

const (
	WagonTypeSeated    WagonType = "seated"
	WagonTypePlatzkart WagonType = "platzkart"
	WagonTypeCoupe     WagonType = "coupe"
	WagonTypeLux       WagonType = "lux"
	WagonTypeSoft      WagonType = "soft"
	WagonTypeSV        WagonType = "sv"
	WagonTypeGeneral   WagonType = "general"
)

type WatchFiltersDTO struct {
	Trains     []string    `form:"trains" json:"trains"`
	Position   Position    `form:"position" json:"position"`
	Zone       Zone        `form:"zone" json:"zone"`
	WagonTypes []WagonType `form:"wagon_types" json:"wagon_types"`
	Classes    []string    `form:"classes" json:"classes"`
}
