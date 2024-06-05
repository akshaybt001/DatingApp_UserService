package helperstruct

type InterestHelper struct {
	InterestId   int
	InterestName string
}

type GenderHelper struct {
	GenderId   int
	GenderName string
}

type FetchUser struct {
	Age int
}

type FetchPreference struct {
	MinAge     int
	MaxAge     int
	Gender     int
	DesireCity string
}

type Home struct {
	Id      string
	Name    string
	Age     int
	Gender  string
	City    string
	Country string
	Images  []string
	Interests []string
}


