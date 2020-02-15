/* 
 * plot - plot on root window
 */

/* 
 * SYNOPSIS
 *
 * #include "plot.h"
 * 
 * void 
 * plotinit (double xmin, double xmax, double ymin, double ymax, 
 * 		const char *fgcolorstr, const char *bgcolorstr, 
 * 		const char *xlabel, const char *ylabel, 
 * 		const char *title,...);
 *
 * void 
 * plot (double x, double y, double xmin, double xmax, double ymin, 
 * 		double ymax, int state);
 *
 * void
 * plotexit(void);
 *
 */

#include <X11/Xlib.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>

#include <math.h>

Display *display;	
Pixmap pix;
GC gc;
int screen;
Window win;		
unsigned int display_width, display_height;
Colormap screen_colormap;
int XY[5];

double Xx (int x, double xmin, double xmax) {
	int width;
	width=display_width;
	return (xmax-xmin)/((double)width-44.0) * x +
	xmin - 22.0 * (xmax-xmin)/((double)width-44.0);
}  
   
double Yy (int y, double ymin, double ymax) {
	int height;
	height=display_height;
	return (ymin-ymax)/((double)height-44.0) * y + 
	ymax - 22.0 * (ymin-ymax)/((double)height-44.0);
}  

int xX (double x, double xmin, double xmax, int width) {
	double w = (double)width - 44.0;
	if (isinf(x)) return width;
	if (isnan(x)) return 0;
	if (x<xmin) return 22;
	if (x>xmax) return width-22;
	return (int)(22.0 + w*x/(xmax-xmin) - w*xmin/(xmax-xmin));
}
int yY (double y, double ymin, double ymax, int height) {
	double h = (double)height - 44.0;
	if (isinf(y)) return 0;
	if (isnan(y)) return height;
	if (y<ymin) return height-22;
	if (y>ymax) return 22;
//	return (int)(22.0 - h*y/(ymax-ymin)-h*ymin/(ymax-ymin));
	return (int)(22.0 + (y-ymax)*h/(ymin-ymax));
}

void plotexit(void) {
	XCloseDisplay(display);
}

/* mode 0: continue 1: newline 
 * style 0: points, 1: line: 2: dashed line, 3: long-dashed */
void plot (double x, double y, double xmin, double xmax, double ymin, 
		double ymax, int mode, int linestyle, int linewidth, 
		const char *colorstr) {
	static int init=1;
	int line_style = LineSolid;	
	int cap_style = CapButt;
	int join_style = JoinBevel; 
	int X,Y;
	static int oldx, oldy;
	XColor color;
	Status rc;
	// LineOnOffDash, LineDoubleDash
	
	line_style = LineOnOffDash;
	switch (linestyle) {
	case 1:
		line_style = LineSolid;
		break;
	case 2:
		line_style = LineOnOffDash;
		break;
	case 3:
		line_style = LineDoubleDash;
		break;
	default:
		line_style = LineSolid; /* 0: Points */
		break;
	}

	rc = XAllocNamedColor (display, screen_colormap, colorstr, &color, &color);
	if (rc == 0) {
		perror("colour fault.");
		exit (1);
	}
	XSetForeground(display, gc, color.pixel);
	
	X = xX(x,xmin,xmax,display_width);
	Y = yY(y,ymin,ymax,display_height);
	if (mode==1) init = 1;
	if (linestyle==0) {
		if (linewidth > 0) {
			//XSetFillStyle(display, gc, FillSolid);
			XFillArc(display, pix, gc, X-linewidth/2, Y-linewidth/2, linewidth, linewidth, 0, 360*64);
		} else {
			XDrawPoint(display, pix, gc, X, Y);
		}
		return;
	}
	if (init) {
		oldx = X;
		oldy = Y;
		init = 0;
	} else {
		XSetLineAttributes(display, gc, linewidth, line_style, 
			cap_style, join_style); 
		XDrawLine(display, pix, gc, oldx, oldy, X, Y);
		oldx = X; 
		oldy = Y;
	}
	
	/* draw two intersecting lines, one horizontal and one vertical,
	 * which intersect at point "50,100".                           
	XDrawLine(display, pix, gc, 50, 0, 50, 200);
	XDrawLine(display, pix, gc, 0, 100, 200, 100);
	XSetWindowBackgroundPixmap (display, win, pix);
	XClearWindow(display,win);
	XFlush(display);
	*/
}

void clearplot(void) {
	XSetWindowBackgroundPixmap (display, win, pix);
	XClearWindow(display,win);
	XFlush(display);

	XSetForeground(display, gc, BlackPixel(display, screen));
 	XFillRectangle(display, pix, gc, 0, 0, display_width, display_height);
	XSetWindowBackgroundPixmap (display, win, pix);
//	XClearWindow(display,win);
//	XFlush(display);
}
void plotbackground() {
	XColor fgcolor, bgcolor;
	Status rc;
}

void plotinit (double xmin, double xmax, double ymin, double ymax, const char *fgcolorstr, 
		const char *bgcolorstr, const char *xlabel, const char *ylabel, 
		const char *title, int width, int height) {


	char *display_name = getenv("DISPLAY");
	char *font_name ="10x20";
	XFontStruct *font_info;
	unsigned long valuemask = 0;
	XGCValues values;
	unsigned int line_width = 1;
	int line_style = LineSolid;	
	int cap_style = CapButt;
	int join_style = JoinBevel; 
	char xrange[100]; 
	char yrange[100]; 
	XColor fgcolor, bgcolor;
	Status rc;

	snprintf(xrange, 99, "[%f, %f]", xmin, xmax);
	snprintf(yrange, 99, "[%f, %f]", ymin, ymax);
	display = XOpenDisplay(display_name);
	if (display == NULL) {
		fprintf(stderr, "cannot connect to X server '%s'\n",
			display_name);
		exit(1);
	}
	/* get the geometry of the default screen for our display. */
	screen = DefaultScreen(display);
	display_width = DisplayWidth(display, screen);
	display_height = DisplayHeight(display, screen);

	if ((width==0)||(height==0)) win = RootWindow(display, screen);
	else { 
		display_width = width;
		display_height = height;
		win = XCreateSimpleWindow(display,RootWindow(display,screen),
			0,0,width,height,0,BlackPixel(display,screen),WhitePixel(display,screen));
		XMapWindow(display,win);
	}
	gc = XCreateGC(display, win, valuemask, &values);
	if (gc < 0) {
		fprintf(stderr, "XCreateGC: \n");
	}
	//XSetForeground(display, gc, WhitePixel(display, screen));
	//XSetBackground(display, gc, BlackPixel(display, screen)); 
	XSetLineAttributes(display, gc, line_width, line_style, cap_style, join_style); 
	//XSetFillStyle(display, gc, FillSolid);

	pix = XCreatePixmap(display, win, display_width, display_height, (unsigned int)DefaultDepth(display, screen));

	XSync(display, False);

	/* Color */
	screen_colormap = DefaultColormap(display, screen);
	rc = XAllocNamedColor (display, screen_colormap, fgcolorstr, &fgcolor, &fgcolor);
	if (rc == 0) {
		perror("colour fault.");
		exit (1);
	}	
	rc = XAllocNamedColor (display, screen_colormap, bgcolorstr, &bgcolor, &bgcolor);
	if (rc == 0) {
		perror("colour fault.");
		exit (1);
	}
	XSetForeground(display, gc, fgcolor.pixel);
	XSetBackground(display, gc, bgcolor.pixel);

	/* Clear Background */
	XSetForeground(display, gc, bgcolor.pixel);
	XFillRectangle(display, pix, gc, 0, 0, display_width, display_height);
	XSetForeground(display, gc, fgcolor.pixel);

	/* Font */
	font_info = XLoadQueryFont(display, font_name);
	XSetFont(display, gc, font_info->fid);

	XDrawString(display, pix, gc, display_width-10*strlen(xlabel)-22, display_height, xlabel, strlen(xlabel));
	XDrawString(display, pix, gc, 22, 18, ylabel, strlen(ylabel));
	XDrawString(display, pix, gc, display_width/2-5*strlen(title), 18, title, strlen(title));
	XDrawString(display, pix, gc, 22, display_height, xrange, strlen(xrange));
	XDrawString(display, pix, gc, display_width-10*strlen(yrange)-22, 18, yrange, strlen(yrange));

	
	XDrawRectangle(display, pix, gc, 22, 22, display_width-44, display_height-44);
	if ((xmin<0.0) && (xmax>0.0)) {
		line_style = LineOnOffDash;
		XSetLineAttributes(display, gc, line_width, line_style, cap_style, join_style); 
		XDrawLine(display, pix, gc, xX(0.0,xmin,xmax,display_width), 
				22, xX(0.0,xmin,xmax,display_width), display_height-22);
	}
	if ((ymin<0.0) && (ymax>0.0)) {
		line_style = LineOnOffDash;
		XSetLineAttributes(display, gc, line_width, line_style, cap_style, join_style); 
		XDrawLine(display, pix, gc, 22, yY(0.0,ymin,ymax,display_height), 
				display_width-22, yY(0.0,ymin,ymax,display_height));
	}

	//XSetForeground(display, gc, WhitePixel(display, screen));
	XSetWindowBackgroundPixmap (display, win, pix);
	XClearWindow(display,win);
	XFlush(display);

}
void flush(void) {
	XSetWindowBackgroundPixmap (display, win, pix);
	XClearWindow(display,win);
	XFlush(display);
}
void rescale(double xmin, double xmax, double ymin, double ymax) {
	int x;
	double Xmin, Xmax, Ymin, Ymax;
	double dx,dy;
	if (XY[0]!=3) {
		if (XY[1]>XY[3]) {
			x=XY[1];
			XY[1]=XY[3];
			XY[3]=x;
		}
		if (XY[2]<XY[4]) {
			x=XY[2];
			XY[2]=XY[4];
			XY[4]=x;
		}
	}
	switch (XY[0]) {
	case (1) : 
		Xmin=Xx(XY[1],xmin,xmax);
		Xmax=Xx(XY[3],xmin,xmax);
		Ymin=Yy(XY[2],ymin,ymax);
		Ymax=Yy(XY[4],ymin,ymax);
		printf("%f %f %f %f\n",Xmin,Xmax,Ymin,Ymax);
		break;
	case(3):
		dx = Xx(XY[3],xmin,xmax)-Xx(XY[1],xmin,xmax);
		dy = Yy(XY[4],ymin,ymax)-Yy(XY[2],ymin,ymax);
		Xmin=xmin-dx;
		Xmax=xmax-dx;
		Ymin=ymin-dy;
		Ymax=ymax-dy;
		printf("%f %f %f %f\n",Xmin,Xmax,Ymin,Ymax);
		break;
	case(4):
		Xmin=xmin+(xmax-xmin)/5.0;	
		Xmax=xmax-(xmax-xmin)/5.0;	
		Ymin=ymin+(ymax-ymin)/5.0;	
		Ymax=ymax-(ymax-ymin)/5.0;	
		printf("%f %f %f %f\n",Xmin,Xmax,Ymin,Ymax);
		break;
	case(5):
		Xmin=xmin-(xmax-xmin)/5.0;	
		Xmax=xmax+(xmax-xmin)/5.0;	
		Ymin=ymin-(ymax-ymin)/5.0;	
		Ymax=ymax+(ymax-ymin)/5.0;	
		printf("%f %f %f %f\n",Xmin,Xmax,Ymin,Ymax);
		break;
	}
	plotexit();
	exit(0);
}
void handle_button_down(XButtonEvent *event, double xmin,double xmax,double ymin, double ymax) {
	int x,y;
	if ((event->button==4)||(event->button==5)) {
		XY[0]=event->button;
		rescale(xmin,xmax,ymin,ymax);
	} else if (event->button==2) {
		fprintf(stderr,"x=%f\ty=%f\n",Xx(event->x,xmin,xmax),Yy(event->y,ymin,ymax));
	} else {
		XY[0]=event->button;
		XY[1]=event->x;
		XY[2]=event->y;
	}
}
void handle_button_release(XButtonEvent *event,double xmin,double xmax, double ymin, double ymax) {
	int x,y;
	XY[3]=event->x;
	XY[4]=event->y;
	if ((event->button==1) || (event->button==3))
		rescale(xmin,xmax,ymin,ymax);
}
void eventloop(double xmin, double xmax, double ymin, double ymax) {
	int done=0;
	XEvent event;
//	XSelectInput(display, win, ExposureMask | KeyPressMask | ButtonPressMask | Button1MotionMask );
	XSelectInput(display, win, ButtonPressMask|ButtonReleaseMask|Button1MotionMask|Button3MotionMask );
	while (!done) {
		XNextEvent(display, &event);
		switch (event.type) {
//		case Expose:
//			printf("expose\n");
//			break;
//		case ConfigureNotify:
//			printf("configureNotify\n");
//			break;
		case ButtonPress:
			handle_button_down((XButtonEvent *)&event.xbutton,xmin,xmax,ymin,ymax);
			break;
		case ButtonRelease:
			handle_button_release((XButtonEvent *)&event.xbutton,xmin,xmax,ymin,ymax);
			break;
		case MotionNotify:
//			printf("motionnotify\n");
			break;
//		case KeyPress:
//			printf("keypress\n");
//			break;
//		default:
//			printf("default\n");
//			break;
		}
	}
}



/*
int main(int argc, char **argv) {
	char *fgcolor = "green";
	char *bgcolor = "black";
	char *xlabel = "X";
	char *ylabel = "Y";
	char *title = "-- (plot) -- ";
	double xmin, xmax, ymin, ymax;
	double x,y;
	int i;

	//x=atof(argv[1]);
	//y=atof(argv[2]);
	xmin = -1.0;
	xmax = 1.0;
	ymin = -1.0;
	ymax = 1.0;

	plotinit (xmin, xmax, ymin, ymax, fgcolor, bgcolor, xlabel, ylabel, title);
	x=0.0;y=0.0;
	for (i=0; i<1000; i++) {
		x+=0.001;
		y=sin(x*2.0*M_PI);
		plot (x, y, xmin, xmax, ymin, ymax, 0);
	}
	plot (x+0.1, y+0.1, xmin, xmax, ymin, ymax, 0);
	plotexit();
	
	return 0;
}*/
