package catalogs

type StationCode string

const (
	Tashkent      StationCode = "2900000"
	TashkentNorth StationCode = "2900001"
	TashkentSouth StationCode = "2900002"
	Samarkand     StationCode = "2900700"
	Bukhara       StationCode = "2900800"
	Khiva         StationCode = "2900172"
	Urgench       StationCode = "2900790"
	Nukus         StationCode = "2900970"
	Navoi         StationCode = "2900930"
	Andijan       StationCode = "2900680"
	Karshi        StationCode = "2900750"
	Jizzakh       StationCode = "2900720"
	Termez        StationCode = "2900255"
	Gulistan      StationCode = "2900850"
	Kokand        StationCode = "2900880"
	Margilan      StationCode = "2900920"
	Pop           StationCode = "2900693"
	Namangan      StationCode = "2900940"
)

type Station struct {
	Code   StationCode
	NameUZ string
	NameRU string
	NameEN string
}

var Stations = []Station{
	{
		Code:   Tashkent,
		NameUZ: "Toshkent",
		NameRU: "Ташкент",
		NameEN: "Tashkent",
	},
	{
		Code:   TashkentNorth,
		NameUZ: "Toshkent Shimoliy",
		NameRU: "Ташкент Северный",
		NameEN: "Tashkent North",
	},
	{
		Code:   TashkentSouth,
		NameUZ: "Toshkent Janubiy",
		NameRU: "Ташкент Южный",
		NameEN: "Tashkent South",
	},
	{
		Code:   Samarkand,
		NameUZ: "Samarqand",
		NameRU: "Самарканд",
		NameEN: "Samarkand",
	},
	{
		Code:   Bukhara,
		NameUZ: "Buxoro",
		NameRU: "Бухара",
		NameEN: "Bukhara",
	},
	{
		Code:   Khiva,
		NameUZ: "Xiva",
		NameRU: "Хива",
		NameEN: "Khiva",
	},
	{
		Code:   Urgench,
		NameUZ: "Urganch",
		NameRU: "Ургенч",
		NameEN: "Urgench",
	},
	{
		Code:   Nukus,
		NameUZ: "Nukus",
		NameRU: "Нукус",
		NameEN: "Nukus",
	},
	{
		Code:   Navoi,
		NameUZ: "Navoiy",
		NameRU: "Навои",
		NameEN: "Navoi",
	},
	{
		Code:   Andijan,
		NameUZ: "Andijon",
		NameRU: "Андижан",
		NameEN: "Andijan",
	},
	{
		Code:   Karshi,
		NameUZ: "Qarshi",
		NameRU: "Карши",
		NameEN: "Karshi",
	},
	{
		Code:   Jizzakh,
		NameUZ: "Jizzax",
		NameRU: "Джизак",
		NameEN: "Jizzakh",
	},
	{
		Code:   Termez,
		NameUZ: "Termiz",
		NameRU: "Термез",
		NameEN: "Termez",
	},
	{
		Code:   Gulistan,
		NameUZ: "Guliston",
		NameRU: "Гулистан",
		NameEN: "Gulistan",
	},
	{
		Code:   Kokand,
		NameUZ: "Qo'qon",
		NameRU: "Коканд",
		NameEN: "Qo'qon",
	},
	{
		Code:   Margilan,
		NameUZ: "Marg'ilon",
		NameRU: "Маргилан",
		NameEN: "Margilon",
	},
	{
		Code:   Pop,
		NameUZ: "Pop",
		NameRU: "Пап",
		NameEN: "Pop",
	},
	{
		Code:   Namangan,
		NameUZ: "Namangan",
		NameRU: "Наманган",
		NameEN: "Namangan",
	},
}
