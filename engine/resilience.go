package engine

type ResilienceResult struct {
	ResilientTopology   string  `json:"resilient_topology"`
	CapacityDropPercent float64 `json:"capacity_drop_percent"`
	FailoverAction      string  `json:"failover_action"`
}

// SimulateFailover, çoklu düğüm (multi-node) yapısında bir node arızası 
// durumunda küp kapasite kaybını ve devreye girecek stratejiyi simüle eder.
func SimulateFailover(totalNodes int, tpSize int, ppSize int) ResilienceResult {
	if totalNodes <= 1 {
		return ResilienceResult{
			ResilientTopology:   "Standalone (Tekil Düğüm)",
			CapacityDropPercent: 100.0,
			FailoverAction:      "Yedeksiz mimari: Node arızasında servis tamamen kesilir.",
		}
	}

	dropPercent := (1.0 / float64(totalNodes)) * 100.0

	return ResilienceResult{
		ResilientTopology:   "Multi-Node Cluster (Redundant Ring)",
		CapacityDropPercent: round(dropPercent, 1),
		FailoverAction:      "Arızalı node izole edilir, trafik kalan node kümesine otomatik balance edilir.",
	}
}
