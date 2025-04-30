package resources

type Snapshot struct {
	Load       float32
	CPU        float32
	Disk       float32
	Net        float32
	TopTalkers []TopTalker
}

type Statistic struct {
	Load       float32
	CPU        float32
	Disk       float32
	Net        float32
	TopTalkers []TopTalker
}

type TopTalker struct {
	ID      int
	Name    string
	LoadNet float32
}

type CollectorEnable struct {
	Load       bool
	CPU        bool
	Disk       bool
	Net        bool
	TopTalkers bool
}
