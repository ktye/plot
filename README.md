# plot - create interactive raster plots

Plot is a plot package for Go that creates interactive plots on images.

## Planned features

- plot types
  - xy, polar, amplitude phase, foto, raster (heatmap)
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
- 3d plots
  - 3d line plots using `fogleman/ln`, interactive
- animations
  - animated line plots (2d/3d)
  - export to gif
- user interfaces
  - bindings to ktye/ui, lxn/walk
