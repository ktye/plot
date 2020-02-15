/* $Id: lode.c,v 1.1 2006/02/10 16:16:01 elmar Exp $ */

/*
 * lode - solve linear differential equations
 *
 * A: (internal: column major matrix)
 * 	A_ik:		A[k*Order+i];
 * 	Re(A_ik):	A[2*(k*Order+i)];
 * 	Im(A_ik):	A[2*(k*Order+i)+1];
 *
 * A = [A00.r A00.i A01.r A01.i ... A0nr. A0ni A10r, A10i, ... A1... ... ...]
 */


#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>

enum	{EULER=0, RK4} ODE = RK4;
char	*progname;

int	 Order = 0;
int	 Complex = 0;
int	 vflag = 0;
int	 tflag = 0;
int	 N = 32;
double	 Xmin = 0.0;
double	 Xmax = 1.0;

void usage(void) {
	fprintf(stderr,"usage: xmin=... xmax=... steps=... %s OPTIONS A.. b..\n",progname);
	fprintf(stderr,"OPTIONS\n");
	fprintf(stderr,"\t-c\tcomplex\n");
	fprintf(stderr,"\t-a\talgorithm (euler, rk4) default: rk4\n");
	fprintf(stderr,"\t-v\tverbose (print input values)\n");
	fprintf(stderr,"\t-t\ttranspose input matrix (col-major notation)\n");
	fprintf(stderr,"ENVIRONMENT\n");
	fprintf(stderr,"\txmin, xmax, steps\n");
	exit(1);
}

void printmatrix(double *A) {
	int	 i,k;
	for (i=0; i<(Complex+1)*Order; i+=(Complex+1)) {
		for (k=0; k<(Complex+1)*Order; k+=(Complex+1)) {
			if (Complex)
				printf("(%g,%g)",A[k*Order+i],A[k*Order+i+1]);
			else
				printf("(%g)",A[k*Order+i]);
		}
		printf("\n");
	}
}
void printstate(double *x, double t) {
	int	 i;
	printf("%g ",t);
	for (i=0; i<(Complex+1)*Order; i++)
		printf("%g ",x[i]);
	printf("\n");
}

void printvector(double *x) {
	int	 i;
	for (i=0; i<(Complex+1)*Order; i+=Complex+1) {
		if (Complex)
			printf("(%g,%g)",x[i],x[i+1]);
		else
			printf("(%g)",x[i]);
		printf("\n");
	}
}

void matvec(double *A, double *b, double *Ab) {
	int	 i,k;
	for (i=0; i<(Complex+1)*Order; i+=Complex+1) {
		Ab[i] = 0.0;
		if (Complex)
			Ab[i+1] = 0.0;
		for (k=0; k<(Complex+1)*Order; k+=Complex+1) {
			if (Complex) {
				Ab[i] += A[i*Order+k]*b[k] - A[i*Order+k+1]*b[k+1];
				Ab[i+1] += A[i*Order+k+1]*b[k] + A[i*Order+k]*b[k+1];
			} else
				Ab[i] += A[k*Order+i]*b[k];
		}
	}
}

void vecassign(double *y, double *x) {
	int	 i;
	for (i=0; i<(Complex+1)*Order; i++)
		y[i] = x[i];
}

void eulerstep(double *A, double *b, double *x, double dt) {
	int	i;
	matvec(A,b,x);
	for (i=0; i<(Complex+1)*Order; i++) {
		x[i] *= dt;
		x[i] += b[i];
	}
}

/* return size of 5 vectors */
double *rk4init(void) {
	return (double *)malloc(5*(Complex+1)*Order*sizeof(double));
}
void rk4step(double *A, double *b, double *x, double dt, double *rk4work) {
	int	 i;
	double	*k1, *k2, *k3, *k4, *tm;
	k1 = rk4work;
	k2 = &rk4work[(Complex+1)*Order];
	k3 = &rk4work[2*(Complex+1)*Order];
	k4 = &rk4work[3*(Complex+1)*Order];
	tm = &rk4work[4*(Complex+1)*Order];

	matvec(A,b,k1);
	for (i=0; i<(Complex+1)*Order; i++) k1[i] *= dt;

	for (i=0; i<(Complex+1)*Order; i++) tm[i] = b[i]+0.5*k1[i];
	matvec(A,tm,k2);
	for (i=0; i<(Complex+1)*Order; i++) k2[i] *= dt;

	for (i=0; i<(Complex+1)*Order; i++) tm[i] = b[i]+0.5*k2[i];
	matvec(A,tm,k3);
	for (i=0; i<(Complex+1)*Order; i++) k3[i] *= dt;

	for (i=0; i<(Complex+1)*Order; i++) tm[i] = b[i]+k3[i];
	matvec(A,tm,k4);
	for (i=0; i<(Complex+1)*Order; i++) k4[i] *= dt;

	for (i=0; i<(Complex+1)*Order; i++)
		x[i] = b[i] + k1[i]/6.0 + k2[i]/3.0 + k3[i]/3.0 + k4[i]/6.0;
}

int isqrt(int x) {
	int	 y;
	if (x<1)
		return 0;
	y = 1;
	while (y*y<x)
		y++;
	return y;
}

int ckorder(int x) {
	int	 i;
	i = isqrt(4*x+1);
	if (i*i == (4*x+1))
		if (!((--i)%2))
			return i/2;
	return -1;
}

int main(int args, char **argv) {
	int	 i, t;
	extern char	*optarg;
	extern int	 optind;
	int	 c;
	char	*var;
	double	*A, *b, *x, *rk4work = (double *)NULL, dt;

	progname = argv[0];
	while ((c = getopt(args, argv, "a:chtv")) != -1) {
		switch (c) {
		case 'a':
			if (!strcmp(optarg,"euler"))
				ODE = EULER;
			else if (!strcmp(optarg,"rk4"))
				ODE = RK4;
			else {
				fprintf(stderr,"error: algorithm %s unknown\n",optarg);
				return 1;
			}
			break;
		case 'c':
			Complex = 1;
			break;
		case 'h': usage();
		case 't':
			tflag = 1;
			break;
		case 'v':
			vflag = 1;
			break;
		default:
			usage();
		}
	}
	args -= optind;
	argv += optind;

	/*
	 * check for Order: 
	 * 	args==Order(1+Order)	real system
	 * 	args==2Order(1+Order)	complex system
	 */
	if (Complex)
		i = args/2;
	else
		i = args;
	if ((Order = ckorder(i))<=0) {
		fprintf(stderr,"error: wrong number of arguments\n");
		return 1;
	}

	if ((var = getenv("xmin")))	Xmin = atof(var);
	if ((var = getenv("xmax")))	Xmax = atof(var);
	if ((var = getenv("steps")))	N = atof(var);
	A = (double *)malloc((Complex+1)*Order*Order*sizeof(double));
	b = (double *)malloc((Complex+1)*Order*sizeof(double));
	x = (double *)malloc((Complex+1)*Order*sizeof(double));

	i = (Complex+1)*Order*Order + (Complex+1)*Order;
	if (args != i) {
		fprintf(stderr,"error: Wrong number of elements: %d (%d needed)\n",args,i);
		return 1;
	}
	for (i=0; i<(Complex+1)*Order*Order; i+=Complex+1) {
		if (tflag)
			if (Complex) {
				A[i] = atof(argv[i]);
				A[i+1] = atof(argv[i+1]);
			} else
				A[i] = atof(argv[i]);
		else
			if (Complex) {
				A[2*(i/2/Order+Order*((i/2)%Order))] = atof(argv[i]);
				A[2*(i/2/Order+Order*((i/2)%Order))+1] = atof(argv[i+1]);
			} else
				A[i/Order+Order*(i%Order)] = atof(argv[i]);
	}
	for (i=0; i<(Complex+1)*Order; i+=Complex+1) {
		if (Complex) {
			b[i] = atof(argv[(Complex+1)*Order*Order + i]);
			b[i+1] = atof(argv[(Complex+1)*Order*Order + i + 1]);
		} else
			b[i] = atof(argv[(Complex+1)*Order*Order + i]);
	}

	if (vflag) {
		printf("lode: Order=%d Complex=%d ODE=%d\n\n",Order,Complex,ODE);
		printf("A=\n");
		printmatrix(A);
		printf("\nb=\n");
		printvector(b);
		printf("\n");
	}

	if (!N)
		return 0;

	switch (ODE) {
	case EULER:
		break;
	case RK4:
		rk4work = rk4init();
		break;
	default:
		return 2;
	}

	printstate(b,Xmin);
	dt = (Xmax-Xmin)/(double)(N-1);
	for (t=1; t<N; t++) {
		switch (ODE) {
		case EULER:
			eulerstep(A, b, x, dt);
			break;
		case RK4:
			rk4step(A, b, x, dt, rk4work);
			break;
		default:
			return 2;
		}
		vecassign(b,x);
		printstate(x,(double)t*(Xmax-Xmin)/(double)(N-1));
	}

	return 0;
}
