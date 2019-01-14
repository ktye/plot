# plot - high level plot package

Plot is a plot package for Go that is one level above gonum/plot.

## Planned features

- interactive plots
  - zoom, pan, select, lineinfo
  - connected to table as an advanced kind of legend
    - one line per dataset
    - any numer of columns
    - double-link to plot (click on row, or click on line selectes the other)
    - additional plot caption text
  - datatype that holds all data
    - export to disk
    - multiple applications can understand it
      - double-click on a plt file, shows plot, table and menues
    - export to xlsx
    - pptx: export plot with caption table as a slide to an existing presentation
    - html
  - slider: react to a slider widget event (scroll through data) 
- stream plots (lots of data)   
  - same interface but do not hold data
  - plot panel size must be known in advance
  - axis must be known or some kind of rescale algorithm
  - plot is done into an image (rasterized) but rescaling still works kind-of
  - axis ticks and labels are vectorgraphics and can be rescaled
  - online mode: show plot updates at slow rate while still reading data
- 3d plots
  - 3d line plots using `fogleman/ln`, interactive
- animations
  - animated line plots (2d/3d)
  - export to gif
- additional plotters
  - polar: complex numbers in a polar diagram
  - bode: complex numbers amplitude and phase over x
  - images in a xy plot
    - colormap images (spectrograms, ...)
    - photos
- user interfaces
  - bindings to duit and walk
