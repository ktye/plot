/* $Id: main.c,v 1.4 2006/04/19 12:15:03 elmar Exp $ */

/*
 *	plot - new xrootplot with more features
 *		o time frames (simulations)
 *		o multiple columns in input stream
 *	#	o multiple screens
 *		o command line arguments
 *
 *	ARGUMENTS
 *		-f		fastplot
 *		-b		black background
 *		-d delay 	delay (in usecs)
 *		-X x0:xN:Nx	x grid
 *		-Y y0:yN:Ny	y grid
 *		-X Nx		x grid
 *		-Y Ny		y grid
 *
 *	FILE SYNTAX (stdin)
 *		:1:2:green :1:2:red ...
 *		x y1 y2 ...
 *		:frame
 *	#	:screen:1
 */

#include <stdio.h>
#include <stdlib.h>
#include <ctype.h>
#include <unistd.h>
#include <string.h>
#include <X11/Xlib.h>
#include <X11/Xutil.h>

#define	MAXLINES	100

double	 xmin=-1.0, xmax=1.0, ymin=-1.0, ymax=1.0;
char 	*colors[]={"green","yellow","red","blue","lightblue","brown"};
int	 Lines = 0;
int	 Delay;
char	 *Xgrid, *Ygrid;

void (*plot)();
void (*background)();

struct linestruct {
	unsigned int	 width;
	unsigned int	 style;	
	char	 color[30];
	int	 oldX, oldY;
	int	 new;
} line[MAXLINES];

Display	*display;	
Pixmap	 pix, backpix, pix0, pix1;
GC 	 gc;
int	 screen;
Window	 win;		
unsigned int	display_width, display_height;
Colormap 	screen_colormap;

double Xx (int x) {
	return (xmax-xmin)/((double)display_width-44.0) * x + xmin - 22.0 * (xmax-xmin)/((double)display_width-44.0);
}  
double Yy (int y) {
	return (ymin-ymax)/((double)display_height-44.0) * y + ymax - 22.0 * (ymin-ymax)/((double)display_height-44.0);
}  
int xX (double x) {
	double w = (double)display_width - 44.0;
	/*
	if (isinf(x)) return display_width-22;
	if (isnan(x)) return 0;
	*/
	if (x<xmin) return 22;
	if (x>xmax) return display_width-22;
	return (int)(22.0 + w*x/(xmax-xmin) - w*xmin/(xmax-xmin));
}
int yY (double y) {
	double h = (double)display_height - 44.0;
	/*
	if (isinf(y)) return 0;
	if (isnan(y)) return display_height-22;
	*/
	if (y<ymin) return display_height-22;
	if (y>ymax) return 22;
	return (int)(22.0 + (y-ymax)*h/(ymin-ymax));
}

void plotexit(void) {
	XCloseDisplay(display);
}

void slowplot(double x, double y, int n) {
	int	 X, Y;
	int	 style = LineSolid;
	int	 join_style = JoinBevel;
	int	 cap_style = CapButt;
	XColor	 fg;
	X = xX(x);
	Y = yY(y);
	if (!line[n].new) {
		switch(line[n].style) {
			case 2: style = LineOnOffDash; break;
			case 3: style = LineDoubleDash; break;
			default: style = LineSolid; break;
		}
		XAllocNamedColor (display, screen_colormap, line[n].color, &fg, &fg);
		XSetForeground(display, gc, fg.pixel); 
		if (line[n].style) {
			XSetLineAttributes(display, gc, line[n].width, style, cap_style, join_style);
			XDrawLine(display, pix, gc, line[n].oldX, line[n].oldY, X, Y);
		}
	}
	if (!line[n].style) {
		XAllocNamedColor (display, screen_colormap, line[n].color, &fg, &fg);
		XSetForeground(display, gc, fg.pixel); 
		if (line[n].width>0)
			XFillArc(display, pix, gc, X-line[n].width/2, Y-line[n].width/2, line[n].width, line[n].width, 0, 360*64);
		else
			XDrawPoint(display, pix, gc, X, Y);
	}
	line[n].oldX = X;
	line[n].oldY = Y;
}


void fastplot(double x, double y, int n) {
	int	X, Y;
	X = xX(x);
	Y = yY(y);
	if (!line[n].new) {
		XDrawLine(display, pix, gc, line[n].oldX, line[n].oldY, X, Y);
	}
	line[n].oldX = X;
	line[n].oldY = Y;
}


void frame(void) {
	XSetWindowBackgroundPixmap (display, win, pix);
	XClearWindow(display,win);
	XFlush(display);

	if (pix == pix0) 
		pix = pix1;
	else
		pix = pix0;
	XCopyArea(display, backpix, pix, gc, 0, 0, display_width, display_height, 0, 0);
}

int split(char **argv, int max, char *input, char *delim) {
	int	 i;
	char	**ap;
	for (ap = argv; (*ap = strsep(&input, delim)) != NULL;)
		if (**ap != '\0') {
			if (++ap >= &argv[max]) break;
		}
	i = 0;
	ap = argv;
	while (*ap++ != NULL)
		i++;
	return i;
}

void grid(int h) {
	int	 i, n;
	char	*v[3];
	double	 x, min=0.0, max=1.0;

	if (h)
		i = split(v, 3, Xgrid, ":");
	else
		i = split(v, 3, Ygrid, ":");
	if ((i!=3)&&(i!=1))
		return;
	n = 0;
	if (i==3) {
		min = atof(v[0]);
		max = atof(v[1]);
		n = atoi(v[2]);
	} else if (i==1) {
		if (h) {
			min = xmin;
			max = xmax;
		} else {
			min = ymin;
			max = ymax;
		}
		n = atof(v[0]);
	}
	for (i=0; i<n; i++) {
		x = min + (double)i*(max-min)/(double)(n-1);
		if (h)
			XDrawLine(display, backpix, gc, xX(x), 22, xX(x), display_height-22);
		else
			XDrawLine(display, backpix, gc, 22, yY(x), display_width-22, yY(x));
	}
}

void blackbg(void) {
}
void defbg(void) {
	XColor	 c;
	XFontStruct *font_info;
	char *font_name = "10x20";
	char	 xrange[100]; 
	char	 yrange[100]; 
	unsigned int	 line_width = 1;
	int	 line_style = LineSolid;	
	int	 cap_style = CapButt;
	int	 join_style = JoinBevel; 

	XAllocNamedColor (display, screen_colormap, "green", &c, &c);
	XSetForeground(display, gc, BlackPixel(display, screen)); 
 	XFillRectangle(display, backpix, gc, 0, 0, display_width, display_height);

	XAllocNamedColor (display, screen_colormap, "darkgreen", &c, &c);
	XSetForeground(display, gc, c.pixel); 
	grid(0);
	grid(1);

	XAllocNamedColor (display, screen_colormap, "green", &c, &c);
	XSetForeground(display, gc, c.pixel); 
	XDrawRectangle(display,backpix,gc,22,22,display_width-44,display_height-44);

	if ((xmin<0.0) && (xmax>0.0)) {
		line_style = LineOnOffDash;
		XSetLineAttributes(display, gc, line_width, line_style, cap_style, join_style); 
		XDrawLine(display, backpix, gc, xX(0.0), 22, xX(0.0), display_height-22);
	}
	if ((ymin<0.0) && (ymax>0.0)) {
		line_style = LineOnOffDash;
		XSetLineAttributes(display, gc, line_width, line_style, cap_style, join_style); 
		XDrawLine(display, backpix, gc, 22, yY(0.0), display_width-22, yY(0.0));
	}

	font_info = XLoadQueryFont(display, font_name);
	XSetFont(display, gc, font_info->fid);

	XDrawString(display, backpix, gc, display_width-10-22, display_height, "X", 1);
	XDrawString(display, backpix, gc, 22, 18, "Y", 1);
	XDrawString(display, backpix, gc, display_width/2-5*8, 18, "- plot -", 8);
	snprintf(xrange, 99, "[%g, %g]", xmin, xmax);
	snprintf(yrange, 99, "[%g, %g]", ymin, ymax);
	XDrawString(display, backpix, gc, 22, display_height, xrange, strlen(xrange));
	XDrawString(display, backpix, gc, display_width-10*strlen(yrange)-22, 18, yrange, strlen(yrange));
	XSetLineAttributes(display, gc, 2, LineSolid, cap_style, join_style); 

}

void plotinit(void) {
	char	*display_name = getenv("DISPLAY");
	unsigned long	 valuemask = 0;
	XGCValues	 values;

	display = XOpenDisplay(display_name);
	if (!display) {
		fprintf(stderr, "cannot connect to X server '%s'\n", display_name);
		exit(1);
	}

	screen = DefaultScreen(display);
	display_width = DisplayWidth(display, screen);
	display_height = DisplayHeight(display, screen);

	win = RootWindow(display, screen);
	gc = XCreateGC(display, win, valuemask, &values);

	XSetForeground(display, gc, WhitePixel(display, screen));
	XSetBackground(display, gc, BlackPixel(display, screen)); 

	pix0 = XCreatePixmap(display, win, display_width, display_height, (unsigned int)DefaultDepth(display, screen));
	pix1 = XCreatePixmap(display, win, display_width, display_height, (unsigned int)DefaultDepth(display, screen));
	backpix = XCreatePixmap(display, win, display_width, display_height, (unsigned int)DefaultDepth(display, screen));
	pix = pix0;
	XSync(display, False);

	screen_colormap = DefaultColormap(display, screen);

	background();
	XClearWindow(display,win);
	XCopyArea(display, backpix, pix, gc, 0, 0, display_width, display_height, 0, 0);

	frame();
}

void loop(void) {
	int	 i;
	double	 x=0.0, y;
	char	 buf[1024];
	char	*vec[MAXLINES], *wec[5];

	while (fgets(buf, sizeof(buf), stdin)) {
		if (buf[strlen(buf)-1]=='\n')
			buf[strlen(buf)-1] = '\0';
		if (!strcmp(buf,":frame")) {
			frame();
			if (Delay)
				usleep(Delay);
			for (i=0; i<Lines; i++)
				line[i].new = 1;
		} else {
			Lines = split(vec, MAXLINES-1, buf, " \t");
			if (!Lines)
				for (i=0; i<Lines; i++)
					line[i].new = 1;
			for (i=0; i<Lines; i++) {
				if (split(wec, 4, vec[i], ":") == 3) {
						line[i].style = atoi(wec[0]);
						line[i].width = atoi(wec[1]);
						strncpy(line[i].color, wec[2], sizeof(line[i].color));
				} else if (!i) {
					x = atof(vec[0]);
				} else {
					y = atof(vec[i]);
					plot(x,y,i-1);
					line[i-1].new = 0;
				}
			}
			if (!Lines)
				for (i=0; i<MAXLINES; i++)
					line[i].new = 1;
		}
	}
	frame();
	plotexit();
}

int main(int args, char **argv) {
	int	 c;
	int	 optstyle=-1, optwidth=2;
	char	*var;
	const char *options = "bfhd:X:Y:s:w:";

	if ((var = getenv("xmin"))) xmin = atof(var);
	if ((var = getenv("ymin"))) ymin = atof(var);
	if ((var = getenv("xmax"))) xmax = atof(var);
	if ((var = getenv("ymax"))) ymax = atof(var);


	background = defbg;
	plot = slowplot;
	Delay = 0;
	Xgrid = (char *)NULL;
	Ygrid = (char *)NULL;
	while ((c = getopt(args, argv, options)) != -1) {
		switch(c) {
		case 'b': background = blackbg; break;
		case 'f': plot = fastplot; break;
		case 'd': Delay = atoi(optarg); break;
		case 'h': printf("OPTIONS: %s\n",options); return 0;
		case 'X': Xgrid = optarg; break;
		case 'Y': Ygrid = optarg; break;
		case 's': optstyle=atoi(optarg); break;
		case 'w': optwidth=atoi(optarg); break;
		default:
			  return 1;
		}
	}
	args -= optind;
	argv += optind;

	for (c=0; c<MAXLINES; c++) {
		line[c].width = (unsigned int)optwidth;
		strncpy(line[c].color, colors[c%6], sizeof(line[c].color));
		if (optstyle>=0)
			line[c].style = (unsigned int)optstyle;
		else
			line[c].style = (unsigned int)(1+c/6%3);
		line[c].new = 1;
	}

	if (args == 4) {
		xmin = atof(argv[0]);
		xmax = atof(argv[1]);
		ymin = atof(argv[2]);
		ymax = atof(argv[3]);
		args -= 4;
	}
	if (args)
		return 1;

	plotinit();
	loop();
	return 0;
}
