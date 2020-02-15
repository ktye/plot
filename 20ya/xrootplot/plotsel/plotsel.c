/* 
 * plotsel - select parameter from plot
 *
 * SYNOPSIS
 *
 * plotsel [XMIN XMAX [YMIN YMAX]]
 *
 * mouse button 1: print coords to stdout
 * mouse button 3: print coords to stderr
 *
 * arguments:
 * 	0:	attempts to read [x-z]{min,max} from env.
 * 	4:	xmin, xmax, ymin, ymax
 * 	2:	xmin, xmax (keeps aspect ratio)
 * 	3:	xmin, xmax, ymedium (keeps aspect ratio)
 *
 * ENVIRONMENT
 * 	xmin, xmax, ymin, ymax
 * 	width of window
 * 	height of window
 *
 * DATAFORMAT:
 * 
 * :linestyle:linewidth:color
 * x y
 *
 */
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include <math.h>

#include "plot_xrootplot.h" 

int WIDTH=640;
int HEIGHT=480;


int main(int args, char **argv) {
	char *fgcolor = "green";
	char *bgcolor = "black";
	char *xlabel = "X";
	char *ylabel = "Y";
	char *title = "-- (plotsel) -- ";
	double xmin, xmax, ymin, ymax;
	int width=0, height=0;
	double x, y;
	int i,n;
	int linemode = 1;
	int linestyle = 1;
	int linewidth = 2;
	char linecolor[100];
	char buf[1024];
	char *var;

	//x=atof(argv[1]);
	//y=atof(argv[2]);
	strncpy (linecolor,"green",99);
	linecolor[99]='\0';

	if (var = getenv("width")) width = atoi(var);
	else width=WIDTH;
	if (var = getenv("height")) height = atoi(var);
	else height=HEIGHT;

	if (args < 3) { 
		xmin=-1.5;
		xmax=1.5;
		ymin=-1.0;
		ymax=1.0;
		if (var = getenv("xmin")) xmin=atof(var);
		if (var = getenv("xmax")) xmax=atof(var);
		if (var = getenv("ymin")) ymin=atof(var);
		if (var = getenv("ymax")) ymax=atof(var);
	} else if (args < 4){
		ymin = 0.5*(xmax+xmin)-
			((double)HEIGHT-44.0)/((double)WIDTH-44.0)*
			0.5*(xmax-xmin);
		ymax = 0.5*(xmax+xmin)+
			((double)HEIGHT-44.0)/((double)WIDTH-44.0)*
			0.5*(xmax-xmin);
	} else if (args < 5){
		xmin = atof(argv[1]);
		xmax = atof(argv[2]);
		ymin = atof(argv[3])-
			((double)HEIGHT-44.0)/((double)WIDTH-44.0)*
			0.5*(xmax-xmin);
		ymax = atof(argv[3])+
			((double)HEIGHT-44.0)/((double)WIDTH-44.0)*
			0.5*(xmax-xmin);
	} else {
		xmin = atof(argv[1]);
		xmax = atof(argv[2]);
		ymin = atof(argv[3]);
		ymax = atof(argv[4]);
	}
	plotinit (xmin, xmax, ymin, ymax, fgcolor, bgcolor, xlabel, 
			ylabel, argv[0], width, height);
	while (fgets(buf, sizeof(buf), stdin)) {
		if (sscanf(buf, ":%d:%d:%s",&linestyle,&linewidth,
					linecolor)==3) {
//			fprintf(stderr, "DEBUG: >%d< >%d< >%s<",linestyle,
//					linewidth,linecolor);
		} else if (sscanf(buf, "%lf %lf",&x, &y) == 2) {
			plot (x, y, xmin, xmax, ymin, ymax, linemode, 
					linestyle, linewidth, linecolor);
			linemode = 0;
		} else {
			linemode = 1;
		}
	}

	flush();
	if ((width!=0)&&(height!=0))
		eventloop(xmin, xmax, ymin,ymax);
	plotexit();

	return 0;
}
