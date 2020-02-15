/* $Id: xt.c,v 1.1 2006/04/18 16:15:12 elmar Exp $ | elmar@abb121 OpenBSD i386 | Sun Jan 26 12:36:13 UTC 2003 */

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <limits.h>
#include <sys/time.h>
char	*progname;

double scale (double a, double b, double A, double B, double x) {
	return (x*(A-B)/(a-b)+(a*B-b*A)/(a-b));
}

void usage(void) {
	fprintf(stderr,"usage: %s\n",progname);
	exit(1);
}

int main(int args, char **argv) {
	int	 i;
	extern char	*optarg;
	extern int	 optind;
	int	 c;
	int	 N=101;
	double 	 xmin=0.0, xmax=1.0, x, t, dt=0.0;
	int	 rt=1;
	char 	*var;
	unsigned long	 sleep=10000;
	struct timeval 	t0, t1;
	struct timeval *t0p, *t1p, *ttp;
	int	 Bout=0;

	if ( (var = getenv("xmin"))!=NULL) xmin=atof(var);
	if ( (var = getenv("xmax"))!=NULL) xmax=atof(var);
	if ( (var = getenv("steps"))!=NULL) N=atoi(var);

	progname = argv[0];
	while ((c = getopt(args, argv, "bd:u:")) != -1) {
		switch (c) {
                case 'b':       /* binary output */
			Bout=1;
			break;
		case 'u':
			sleep = strtoul(optarg,NULL,10);
			break;
		case 'd':
			dt = atof(optarg);
			rt = 0;
			break;
		default:
			usage();
		}
	}
	args -= optind;
	argv += optind;

	if (args>2) {
		xmin=atof(argv[1]);
		xmax=atof(argv[2]);
	}
	if (args==4) {
		N=atoi(argv[4]);
	}

	t0p = &t0;
	t1p = &t1;
	gettimeofday(t0p,NULL);
	t=0.0;
	for (;;) {
		for (i=0; i<N; i++) {
			x=scale(0.0,(double)(N-1),xmin,xmax,i);
			if (Bout) {
				fwrite(&x,sizeof(double),1,stdout);
				fwrite(&t,sizeof(double),1,stdout);
			} else {
				printf("%g %g\n",x,t);
			}
		}
		fflush(stdout);
		if (sleep>0)
			usleep(sleep);
		if (rt) {
			gettimeofday(t1p,NULL);
			dt = (double)(t1p->tv_sec-t0p->tv_sec);
			dt += (double)(t1p->tv_usec-t0p->tv_usec)*0.000001;
			ttp = t0p;
			t0p = t1p;
			t1p = ttp;
		}
		t += dt;
	}

	return 0;
}
