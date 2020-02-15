/* $Id: space.c,v 1.1.1.1 2003/09/28 22:43:14 elmar Exp $ */
/* space - generate N dimensional space

 * loop - inflexion card wrapper
 *
 * N = (args-1)/3
 *
 * for n in 1:N
 * 	N = [ argv[n+1]; argv[n+2] ]_argv[n+3]
 *
 */

#include <stdio.h>
#include <stdlib.h>
void print(double *z, int N) {
	int i;
	for (i=0; i<N; i++)	printf("%f ", z[i]);
	printf("\n");
}

void loop(double *z, double *XMIN, double *XMAX, int *NX, int n, int N) {
	int i;
	double dx,xmin,xmax;
	int X;

	xmin=XMIN[n];
	xmax=XMAX[n];
	X=NX[n];
	dx=(xmax-xmin)/((double)X-1.0);
	z[n] = xmin;
	for (i=0; i<X; i++) {
		if (n==(N-1))	print(z,N);	 
		else 		loop(z, XMIN, XMAX, NX, n+1, N);
		z[n] += dx;
	}
}

int main (int args, char **argv) {
	int N;
	int i;
	int X;
	double x,xmin,xmax,dx;
	double *z,*XMIN,*XMAX;
	int *NX;

	if ((args<2)||(((args-1)%3)!=0)) {
		fprintf(stderr,"usage %s alpha_min alpha_max N_alpha [[R0min R0max NR0] ... []] \n",argv[0]);
		exit(1);
	}


	N = (args - 1)/3;

	XMIN = (double *)malloc((size_t)N * sizeof(double));
	XMAX = (double *)malloc((size_t)N * sizeof(double));
	NX = (int *)malloc((size_t)N * sizeof(int));
	for (i=0; i<N; i++) {
		XMIN[i]=atof(argv[3*i+1]);
		XMAX[i]=atof(argv[3*i+2]);
		NX[i]=atoi(argv[3*i+3]);
	}
	z = (double *)malloc((size_t)N * sizeof(double));

	/*
	for (i=0; i<N; i++) {
		printf("%d: [%f,%f](%d)\n",i,XMIN[i],XMAX[i],NX[i]);
	}
	*/

	loop(z, XMIN, XMAX, NX, 0, N);

	return 0;
}
