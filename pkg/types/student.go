package types

type StudentRequestRow struct {
	Name string `json:"name"`
	PAJ  string `json:"paj"`
	DA   string `json:"da"`
	PP   string `json:"pp"`
	MDS  string `json:"mds"`
	SGSC string `json:"sgsc"`
	IBD  string `json:"ibd"`
	BT   string `json:"bt"`
}

type StudentResponseRow struct {
	Name string  `json:"name"`
	GPA  float64 `json:"mean"`
}
