package entities

//MySQLから持ってくるデータを割り当てるための構造体（最短経路探索用)
type GraphPoint struct {
	ID int
}

type GraphPath struct {
	StartPointID       int
	DestinationPointID int
	Cost               int
}

//MySQLから持ってくるデータを割り当てるための構造体(Point)
type Point struct {
	ID        int     `json:"id"`
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
	NameJP    string  `json:"title"`
	AddressJP string  `json:"address"`
	OutlineJP string  `json:"outline"`
	Access    string  `json:"access"`
	Tel       string  `json:"tel"`
	Business  string  `json:"business"`
	OffDay    string  `json:"off_day"`
	Fee       string  `json:"fee"`
	HowToBook string  `json:"how_to_book"`
	Parking   string  `json:"parking"`
	Remark    string  `json:"remark"`
	ImageUrl1 string  `json:"image_url_1"`
	ImageUrl2 string  `json:"image_url_2"`
	ImageUrl3 string  `json:"image_url_3"`
	ImageUrl4 string  `json:"image_url_4"`
	ImageUrl5 string  `json:"image_url_5"`
	ViewOrder int     `json:"direction_balloon"`
	SpotVideo string  `json:"video"`
}

type AllSpot struct {
	Spots []Point `json:"spots"`
}

//MySQLから持ってくるデータを割り当てるための構造体(Path)
type Path struct {
	ID                 int
	StartPointID       int
	DestinationPointID int
	Cost               int
	Transport          string
	Polyline           string
}

//MySQLから持ってくるデータを割り当てるための構造体(Edge)
type Edge struct {
	ID               int
	BeforePathID     int
	AfterPathID      int
	Angle            float64
	ConnectionStatus string
}

//JSONにして返却するデータ
type OutputRoutes struct {
	Cost     int         `json:"cost"`
	Routes   []RouteData `json:"routes"`
	Polyline []string    `json:"polyline"`
}

type VideoLocations struct {
	Spots []RouteData `json:"spots"`
}

//返却するデータの一部
type RouteData struct {
	Id             int             `json:"id"`
	StartPoI       int             `json:"startSpotId"`
	DestinationPoI int             `json:"endSpotId"`
	VideoURL       string          `json:"videoURL"`
	Locations      []PolyLocations `json:"locations"`
}

type PolyLocations struct {
	DurationSeconds int     `json:"durationSeconds"`
	Latitude        float64 `json:"lat"`
	Longitude       float64 `json:"lng"`
}

type Video struct {
	Id        int
	PathId    int
	VideoType string
	VideoName string
	VideoUrl  string
}
