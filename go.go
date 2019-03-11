package plot

// Result is returned by Routine.
type Result struct {
	Plot  Plot
	Index int
	Err   error
}

// Routine should be called as a go routine to do multiple plots in parallel.
func Routine(p Plotter, index int, res chan Result) {
	plt, err := p.Plot()
	res <- Result{
		Plot:  plt,
		Index: index,
		Err:   err,
	}
}
