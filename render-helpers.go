package main

import (
	"fmt"
	"html/template"
	"math"
	"time"

	"github.com/iamanders/tdr/models"
)

var templateFunctions = template.FuncMap{
	"timeDurationMultiple": models.TimeDurationMultiple,
	"formatDuration":       formatDuration,
}

// Format time duration with format hh:mm
func formatDuration(td time.Duration) string {
	hours := int(math.Floor(td.Hours()))
	minutes := int(int(td.Minutes()) - (hours * 60))
	return fmt.Sprintf("%02d:%02d", hours, minutes)
}
