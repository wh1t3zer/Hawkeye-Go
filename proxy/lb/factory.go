package lb

// LbType ...
type LbType int

const (
	// LbRandom ...
	LbRandom LbType = iota
	// LbRoundRobin ...
	LbRoundRobin
	// LbWeightRoundRobin ...
	LbWeightRoundRobin
	// LbConsistentHash ...
	LbConsistentHash
)

// LoadBalance ...
type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	Update()
}

// LoadBalanceFactory ...
func LoadBalanceFactory(lbtype LbType) LoadBalance {
	switch lbtype {
	case LbRandom:
		return &RandomBalance{}
	case LbRoundRobin:
		return &RoundRobinBalance{}
	case LbWeightRoundRobin:
		return &WeightRoundRobinBalance{}
	case LbConsistentHash:
		return NewConsistentHashBanlance(10, nil)
	default:
		return &RandomBalance{}
	}
}
