/* $Id: grid.c,v 1.1 2004/06/22 16:14:37 elmar Exp $ */

/*
 * grid - generate co-odinate system
 */

static char id[]="$Id: grid.c,v 1.1 2004/06/22 16:14:37 elmar Exp $";
#include <stdio.h>
#include <stdlib.h>
#include <math.h>

void usage(const char **argv) {
	fprintf(stderr,"usage: %s [xlevel] [ylevel]\n",argv[0]);
	exit(1);
}

int main (int args, char **argv) {
	int i, xlevel, ylevel;
	char *var;
	double xmin, ymin, xmax, ymax, nx, ny, expx, expy;
	double x, y, dx, dy;
	int Nx, Ny;

	switch (args) {
	case 1:
		xlevel = 1;
		ylevel = 1;
		break;
	case 2:
		xlevel = atoi(argv[1]);
		ylevel = xlevel;
		break;
	case 3:
		xlevel = atoi(argv[1]);
		ylevel = atoi(argv[2]);
		break;
	default:
		xlevel = atoi(argv[1]);
		ylevel = atoi(argv[2]);
		break;
//		usage(argv);
	}

	if ((xlevel < -1)||(xlevel>3))
		xlevel = 1;
	if ((ylevel < -1)||(ylevel>3))
		ylevel = 1;

	xmin = 0.0;
	ymin = 0.0;
	xmax = 1.0;
	ymax = 1.0;
	if (var = getenv("xmin")) xmin=atof(var);
	if (var = getenv("xmax")) xmax=atof(var);
	if (var = getenv("ymin")) ymin=atof(var);
	if (var = getenv("ymax")) ymax=atof(var);

	nx = rint(log10(xmax-xmin)) - xlevel;
	ny = rint(log10(ymax-ymin)) - ylevel;
	expx = pow(10,nx);
	expy = pow(10,ny);

	dx = (xmax-xmin)/100.0;
	dy = (ymax-ymin)/100.0;
	
	printf(":1:1:green\n");
	for (x=ceil(xmin/expx)*expx+expx/2.0; x<xmax; x += expx) {
		printf("%g %g\n%g %g\n\n",x, ymin, x, ymin+dy);	
		printf("%g %g\n%g %g\n\n",x, ymax, x, ymax-dy);	
	}
	for (y=ceil(ymin/expy)*expy+expy/2.0; y<ymax; y += expy) {
		printf("%g %g\n%g %g\n\n",xmin, y, xmin+dx, y);	
		printf("%g %g\n%g %g\n\n",xmax, y, xmax-dx, y);	
	}

	printf(":1:1:green\n");
	for (x=ceil(xmin/expx)*expx; x<xmax; x += expx) {
		printf("%g %g\n",x, ymin);	
		printf("%g %g\n\n",x, ymax);	
	}
	for (y=ceil(ymin/expy)*expy; y<ymax; y += expy) {
		printf("%g %g\n",xmin, y);	
		printf("%g %g\n\n",xmax, y);	
	}

	printf(":1:2:green\n");
	if ((xmin<0.0) && (xmax>0.0))
		printf("%g %g\n%g %g\n\n",0.0,ymin,0.0,ymax);
	if ((ymin<0.0) && (ymax>0.0))
		printf("%g %g\n%g %g\n\n",xmin,0.0,xmax,0.0);
	

	return 0;
}
