package snowman

type Topological struct {
	pollNumber int
	head       int
	height     int
	blocks     map[int]*snowmanBlock
}
