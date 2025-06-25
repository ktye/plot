package clipboard

type CopyFormat struct {
	Width       int    // image width (px)
	Height      int    // image height (px)
	PlotFont    string // plot font family
	F1          int    // plot font title size
	F2          int    // plot font number label size
	CaptionFont string // caption font family
	F3          int    // caption font size
}
