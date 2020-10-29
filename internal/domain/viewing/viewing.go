package viewing

import "fmt"

// Viewing ...
type Viewing struct {
	VideosWatched     int `redis:"videos_watched"  json:"videosWatched"`
	LastViewProcessed int `redis:"last_view_processed" json:"-"`
}

// NewViewing ...
func NewViewing() *Viewing {
	return &Viewing{}
}

const viewingFormat string = "<Viewing VideosWatched=%v LastViewProcessed=%v>"

func (v *Viewing) String() string {
	return fmt.Sprintf(viewingFormat, v.VideosWatched, v.LastViewProcessed)
}
