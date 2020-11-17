package csv

type CSV struct {
	path string
	rows [][]string
}

func New(path string) *CSV {
	return &CSV{
		path: path,
		rows: [][]string{},
	}
}
