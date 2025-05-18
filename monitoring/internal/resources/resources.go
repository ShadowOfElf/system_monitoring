package resources

type Snapshot struct {
	Load       float32
	CPU        float32
	Disk       map[string]float32
	Net        map[string]int64
	TopTalkers []TopTalker
}

type Statistic struct {
	Load       float32
	CPU        float32
	Disk       map[string]float32
	Net        map[string]int64
	TopTalkers []TopTalker
}

type TopTalker struct {
	Name    string
	LoadNet int
}

type CollectorEnable struct {
	Load       bool
	CPU        bool
	Disk       bool
	Net        bool
	TopTalkers bool
}
