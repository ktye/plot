void plotinit (double xmin, double xmax, double ymin, double ymax, const char *fgcolorstr, 
		const char *bgcolorstr, const char *xlabel, const char *ylabel, 
		const char *title, int width, int height)
;
void plotexit(void) ;


void plot (double x, double y, double xmin, double xmax, double ymin, 
		double ymax, int mode, int linestyle, int linewidth, 
		const char *colorstr) ;
void flush(void);
void eventloop(double , double , double , double );
