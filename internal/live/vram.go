package live

import (
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

const Disclaimer = "Measured on this box, this run. Not a vendor SLA."

type Device struct {
	Index     int     `json:"index"`
	Name      string  `json:"name"`
	TotalMB   float64 `json:"total_mb"`
	UsedMB    float64 `json:"used_mb"`
	FreeMB    float64 `json:"free_mb"`
}

type Snapshot struct {
	Available  bool     `json:"available"`
	Source     string   `json:"source"`
	Disclaimer string   `json:"disclaimer"`
	Devices    []Device `json:"devices"`
	Note       string   `json:"note"`
}

func SnapshotNow() Snapshot {
	s := Snapshot{Disclaimer: Disclaimer, Source: "nvidia-smi", Devices: []Device{}}
	out, err := exec.Command("nvidia-smi", "--query-gpu=index,name,memory.total,memory.used,memory.free", "--format=csv,noheader,nounits").Output()
	if err != nil {
		s.Note = "no nvidia-smi on this box"
		return s
	}
	for _, line := range strings.Split(string(bytes.TrimSpace(out)), "\n") {
		parts := splitCSV(line)
		if len(parts) < 5 {
			continue
		}
		idx, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
		tot, _ := strconv.ParseFloat(strings.TrimSpace(parts[2]), 64)
		used, _ := strconv.ParseFloat(strings.TrimSpace(parts[3]), 64)
		free, _ := strconv.ParseFloat(strings.TrimSpace(parts[4]), 64)
		s.Devices = append(s.Devices, Device{Index: idx, Name: strings.TrimSpace(parts[1]), TotalMB: tot, UsedMB: used, FreeMB: free})
	}
	s.Available = len(s.Devices) > 0
	if !s.Available {
		s.Note = "nvidia-smi returned no rows"
	}
	return s
}

func splitCSV(line string) []string {
	parts := strings.Split(line, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		out = append(out, strings.TrimSpace(p))
	}
	return out
}
