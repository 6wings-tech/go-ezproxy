package git

import (
	"encoding/json"
	"time"
)

type Tag struct {
	Version string    // version string
	Time    time.Time // commit time
}

func (t Tag) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Version string
		Time    string
	}{
		t.Version,
		t.Time.Format(time.RFC3339),
	})
}
