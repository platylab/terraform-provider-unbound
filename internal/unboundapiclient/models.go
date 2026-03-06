package unboundapiclient

// LocalZone -
type LocalZone struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// LocalData -
type LocalData struct {
	Id     int    `json:"id"`
	Domain string `json:"domain"`
	Ttl    int    `json:"ttl"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}
