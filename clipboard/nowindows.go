//go:build !windows

package clipboard

import "github.com/ktye/plot"

func CopyToClipboard(plots plot.Plots, columns int, hi []plot.HighlightID, f CopyFormat) { return }
