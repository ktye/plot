/* 
 * xrootplot - plot on root window
 *
 * SYNOPSIS
 *
 * xrootplot [XMIN XMAX [YMIN|YMED] YMAX]]
 *
 * arguments:
 * 	0:	attempts to read [x-z]{min,max} from env.
 * 	4:	xmin, xmax, ymin, ymax
 * 	2:	xmin, xmax (keeps aspect ratio)
 * 	3:	xmin, xmax, ymedium (keeps aspect ratio)
 *
 * ENVIRONMENT
 * 	xmin, xmax, ymin, ymax
 * 	width of window, if unset: rootwindow
 * 	height of window, if unset: rootwindow
 *
 * DATAFORMAT:
 * 
 * :linestyle:linewidth:color
 * x y
 *
 */
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include <math.h>

#include "plot_xrootplot.h" 

/* bad thing to do (should not mess with plot library) */
int WIDTH=1280;
int HEIGHT=960;
char	*filename;


int main(int args, char **argv) {
	char fgcolor[80];
	char bgcolor[80];
	char xlabel[80];
	char ylabel[80];
	char title[80];
	double xmin, xmax, ymin, ymax;
	unsigned int width=0, height=0;
	double x, y;
	int i,n;
	int	 xoff = 10, yoff=10;
	int linemode = 1;
	int linestyle = 1;
	int linewidth = 2;
	char linecolor[100];
	char buf[1024];
	char *var;
	int	 c;
	extern char     *optarg;
	extern int       optind;


	//x=atof(argv[1]);
	//y=atof(argv[2]);
	strncpy (linecolor,"green",99);
	linecolor[99]='\0';

	if (var = getenv("width")) width = (unsigned int)atoi(var);
	else width=0;
	if (var = getenv("height")) height = (unsigned int)atoi(var);
	else height=0;

	strlcpy(fgcolor,"green",sizeof(fgcolor));
	strlcpy(bgcolor,"black",sizeof(bgcolor));
	strlcpy(xlabel,"X",sizeof(xlabel));
	strlcpy(ylabel,"Y",sizeof(ylabel));
	strlcpy(title,argv[0],sizeof(title));

	filename = (char *)NULL;
	while ((c = getopt(args, argv, "g:x:y:t:f:b:F:")) != -1) {
		switch(c) {
			/*
		case 'w': width = atoi(optarg); break;
		case 'h': height = atoi(optarg); break;
		case 'X': xoff = atoi(optarg); break;
		case 'Y': yoff = atoi(optarg); break;
		*/
		case 'g': XParseGeometry(optarg, &xoff, &yoff, &width, &height); break;
		case 't': strlcpy(title,optarg,sizeof(title)); break;
		case 'x': strlcpy(xlabel,optarg,sizeof(xlabel)); break;
		case 'y': strlcpy(ylabel,optarg,sizeof(ylabel)); break;
		case 'f': strlcpy(fgcolor,optarg,sizeof(fgcolor)); break;
		case 'b': strlcpy(bgcolor,optarg,sizeof(bgcolor)); break;
		case 'F': filename = optarg; break;
		default:
			return 1;
		}
			
	}
	args -= optind;
	argv += optind;

	if (args < 2) { 
		xmin=-1.5;
		xmax=1.5;
		ymin=-1.0;
		ymax=1.0;
		if (var = getenv("xmin")) xmin=atof(var);
		if (var = getenv("xmax")) xmax=atof(var);
		if (var = getenv("ymin")) ymin=atof(var);
		if (var = getenv("ymax")) ymax=atof(var);
	} else if (args < 3){
		ymin = 0.5*(xmax+xmin)-
			((double)HEIGHT-44.0)/((double)WIDTH-44.0)*
			0.5*(xmax-xmin);
		ymax = 0.5*(xmax+xmin)+
			((double)HEIGHT-44.0)/((double)WIDTH-44.0)*
			0.5*(xmax-xmin);
	} else if (args < 4){
		xmin = atof(argv[0]);
		xmax = atof(argv[1]);
		ymin = atof(argv[2])-
			((double)HEIGHT-44.0)/((double)WIDTH-44.0)*
			0.5*(xmax-xmin);
		ymax = atof(argv[2])+
			((double)HEIGHT-44.0)/((double)WIDTH-44.0)*
			0.5*(xmax-xmin);
	} else {
		xmin = atof(argv[0]);
		xmax = atof(argv[1]);
		ymin = atof(argv[2]);
		ymax = atof(argv[3]);
	}
	plotinit (xmin, xmax, ymin, ymax, fgcolor, bgcolor, xlabel, 
			ylabel, title, width, height, xoff, yoff);
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
	// Window mode
	if ((width==0)&&(height==0)) sleep(2);
	return 0;
}
