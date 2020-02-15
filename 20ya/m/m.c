/* $Id: m.c,v 1.9 2006/04/20 15:33:20 elmar Exp $ */

/*
 * m
 *
 * USAGE
 * 	m -- STACK			start with no stack
 * 	m -N INPUT_COLS	-- STACK	read from stdin
 * 	m -x -- STACK			start with 1 col-stack
 * 	m -xt -- STACK			start with x,t (2 col stack) (infinite loop over t)
 * 	m -l				print list of functions
 *
 * OPTIONS
 * 	-N			number of input cols
 * 	-b [i|o]		binary in and/or output (no -x)
 * 	-f format		output format (for ascii output), default: "%g"
 * 	-x			x domain in first col (no stdin) 
 * 				environment vars: xmin, xmax, steps
 * 	-t			time loop (on second col)(also needs -x) (on output 2nd col is deleted)
 * 	-d dt			t increment for (-t option); otherwise realtime
 * 	-u delay		delay in microseconds (-t option)
 * 	-F LINES		flush stdout after LINES (-bo only)
 *
 * STACK OPERATORS
 * 	=VAR		assigns last element from stack to named variable VAR
 * 	VAR		push VAR on stack
 * 	$N		Nth element from stack
 * 	-$N		kill Nth element from stack
 * 	x		abbrevation: 1st element from stack
 * 	t		abbrevation: 2st element from stack
 * 	1.2 pi sqrt2	push constants on stack
 * 	+ - * / sin cos ...
 *
 * EXAMPLE
 * 	m -xt 'x t sub sin' | plot
 *
 * 	oscillating pipe flow:
 * 	xmin=-1 xmax=1 ymin=-2 ymax=2 steps=64; ii="sqrt1_2 sqrt1_2"  m -xt 8 \=alpha 1 0 alpha x mul 0 $ii cmul J0 alpha 0 $ii cmul J0 cdiv csub t 0 i cmul cexp cmul | plot
 */

/*
 * TODO:
 * implement roll (postscript)
 */

/*
 * BUG: slow input routine mread()
 */

#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include <math.h>
#include <ctype.h>
#include <sys/types.h>
#include <regex.h>

#include "include"

#define	SSIZE	256

double	 TIME = 0.0;
double	 xmin=0.0, xmax=1.0;
int	 steps=10;
struct timeval	 T0;

enum {
#include "enum"
	m_assign, m_position, m_kill, m_clear, m_roll, m_const, m_var
};

typedef struct {
	char	*name;
	double	 val;
} var;

char	*progname;

void usage(void) {
	fprintf(stderr,"usage: %s\n",progname);
	exit(1);
}

int hasvar(var *v, char *s, int N) {
	int	 i;
	for (i=0; i<N; i++)
		if (!strcmp(v[i].name,s))
			return i;
	return -1;
}

void inctime(double dt) {
	struct timeval	 t1;
	double	 d;

	if (dt!=0.0)
		TIME += dt;
	else {
		gettimeofday(&t1,NULL);
		d = (double)(t1.tv_sec-T0.tv_sec);
		d+= (double)(t1.tv_usec-T0.tv_usec)*0.000001;
		T0 = t1;
		TIME += d;
	}
}
int mread(double *X, int n, int Bin, int xflag, int tflag, double dt, int *N) {
	int	 i;
	static int	 xc = 0, tc = 0;
	static int	 first = 1, firsttime = 1;


	*N = 0;
	if (xflag) {
		if ((!tflag)&&(xc==steps))
			return 0;
		if (xc==steps)
			xc = 0;
		*N = 1;
		X[0] = xmin + (double)xc*(xmax-xmin)/(double)(steps-1);
		if (tflag) {
			*N = 2;
			if ((!(tc%steps))&&(!firsttime)) {
				tc = 0;
				inctime(dt);
			}
			if (firsttime) {
				if (dt==0.0)
					gettimeofday(&T0, NULL);
				firsttime = 0;
			}
			tc++;
			X[1] = TIME;
		}
		xc++;
		return 1;
	} else if (Bin) {
		if (fread(X,sizeof(double),n,stdin)==n) {
			*N = n;
			return 1;
		} else
			return 0;
	} else if (n>0) {
		for (i=0; i<n; i++) {
			if (scanf("%lf",&X[i])!=1)
				return 0;
		}
		*N=n;
		return 1;
	} else {
		if (first) {
			first = 0;
			*N=0;
			return 1;
		} else
			return 0;
	}
	return 0;
}

char *nextstr(const char *str, int *pos, int len) {
        char    *s;
        for (; (*pos<=len)&&(str[*pos]=='\0'); (*pos)++) {
        }
        if (*pos>=len)
                return (char *)NULL;
        else {
                s = (char *)&str[*pos];
                *pos += strlen(str+*pos) + 1;
                return s;
        }
}

void ckstacksize(int N, int i, int o) {
	if (((N-i)<0)||((N+o-i)>SSIZE))
		exit(1);
}

void mwrite(double *X, int N, int Bout, const char *format, int tflag, useconds_t u, int F) {
	int	 i;
	static int	 first=1;
	static double	 t0;
	static int	 c;
	if (Bout) {
		fwrite(X, sizeof(double), N, stdout);
		if (F) {
			c++;
			if (c==F) {
				fflush(stdout);
				c=0;
			}
		}
	} else if (N>0) {
		if (tflag) {
			if (first&&(N>1)) {
				t0 = X[1];
				first = 0;
			} else {
				if (X[1]!=t0) {
					printf(":frame\n");
					fflush(stdout);
					usleep(u);
				}
				t0 = X[1];
			}
		}
		for (i=0; i<N-1; i++) {
			if ((tflag&&(i!=1))||(!tflag)) {
				printf(format,X[i]);
				printf(" ");
			}
		}
		if ((tflag&&(N-1!=1))||(!tflag))
			printf(format,X[N-1]);
		printf("\n");
	}
}

int main(int args, char **argv) {
	int	 i, j, len, pos, n=0, in, out, F=0;
	useconds_t	 u=10000;
	double	 dt = 0.0;
	extern char	*optarg;
	extern int	 optind;
	int	 c;
	char	*format, formatg[]="%g", *cmd;
	int	 Bin=0, Bout=0, tflag=0, xflag=0;
	double	 cstack[SSIZE];	/* constants  */
	int	 fstack[SSIZE];	/* functions */
	int	 pstack[SSIZE];	/* positional constants  */
	int	 vstack[SSIZE];	/* variables */
	int	 ncs=0, nfs=0, nps=0, nvs=0;
	var	 variables[SSIZE];
	int	 nvars=0;
	int	 N=0;
	double	 X[SSIZE], *x, y[5];
	regex_t	 regfloat;


	if ( (cmd = getenv("xmin"))!=NULL) xmin=atof(cmd);
	if ( (cmd = getenv("xmax"))!=NULL) xmax=atof(cmd);
	if ( (cmd = getenv("steps"))!=NULL) steps=atoi(cmd);

	if (regcomp(&regfloat, "^[+-]*[0-9]*\\.*[0-9]*$", REG_BASIC))
		return 1;

	format = formatg;
	progname = argv[0];
	while ((c = getopt(args, argv, "b:d:f:F:lN:tu:x")) != -1) {
		switch (c) {
		case 'b':	/* binary in/out */
			if (strchr(optarg,'i'))
				Bin = 1;
			if (strchr(optarg,'o'))
				Bout = 1;
			break;
		case 'd':	/* dt */
			dt = atof(optarg);
			break;
		case 'f':	/* output format */
			format=optarg;
			break;
		case 'F':
			F = atoi(optarg);
			if (F<1) {
				fprintf(stderr,"error (%s): -F (<1)\n",progname);
				return 1;
			}
			break;
		case 'N':	/* size of input stack */
			n = atoi(optarg);
			if ((n<0)||(n>SSIZE)) {
				fprintf(stderr,"error (%s): -N (<0||>SSIZE)\n",progname);
				return 1;
			}
			break;
		case 'l':
#include "list"
			return 0;
		case 't':	/* time domain */
			tflag = 1;
			break;
		case 'u':	/* sleep nanoseconds */
			u = strtoul(optarg,NULL,10);
			break;
		case 'x':	/* x domain */
			xflag = 1;
			break;
		default:
			usage();
		}
	}
	args -= optind;
	argv += optind;

	for (i=0; i<args; i++) {
		len = strlen(argv[i]);
		for (j=0; j<len; j++)
			if ((argv[i][j]==' ')||(argv[i][j]=='\t')||(argv[i][j]=='\n'))
				argv[i][j] = '\0';
		pos = 0;
		while ((cmd=nextstr(argv[i],&pos,len))!=NULL) {
			if ((ncs==SSIZE)||(nfs==SSIZE)||(nps==SSIZE)||(nvs==SSIZE))
				return 1;
			if (!strcmp("+",cmd)) { fstack[nfs++] = m_add; }
			else if (!strcmp("-",cmd)) { fstack[nfs++] = m_sub; }
			else if (!strcmp("*",cmd)) { fstack[nfs++] = m_mul; }
			else if (!strcmp("/",cmd)) { fstack[nfs++] = m_div; }
			else if (!strcmp("^",cmd)) { fstack[nfs++] = m_pow; }
			else if (!strcmp(">=",cmd)) { fstack[nfs++] = m_ge; }
			else if (!strcmp("<=",cmd)) { fstack[nfs++] = m_le; }
			else if (!strcmp(">",cmd)) { fstack[nfs++] = m_gt; }
			else if (!strcmp(">",cmd)) { fstack[nfs++] = m_lt; }
			else if (!strcmp("=",cmd)) { fstack[nfs++] = m_eq; }
			else if (!strcmp("++",cmd)) { fstack[nfs++] = m_cadd; }
			else if (!strcmp("--",cmd)) { fstack[nfs++] = m_csub; }
			else if (!strcmp("**",cmd)) { fstack[nfs++] = m_cmul; }
			else if (!strcmp("//",cmd)) { fstack[nfs++] = m_cdiv; }
			else if (!strcmp("||",cmd)) { fstack[nfs++] = m_cabs; }
			else if (!strcmp("==",cmd)) { fstack[nfs++] = m_ceq; }
			else if (cmd[0]=='=') { 
				j=hasvar(variables, &cmd[1], nvars); 
				if (j>=0) {
					vstack[nvs++] = j;
				} else {
					variables[nvars].name = &cmd[1];
					vstack[nvs++] = nvars++;
				}
				fstack[nfs++]= m_assign; 
			} else if (!strcmp("clear",cmd)) { fstack[nfs++] = m_clear;
			} else if (!strcmp("roll",cmd)) { fstack[nfs++] = m_roll;
			} else if (cmd[0]=='$') {
				fstack[nfs++]= m_position; 
				pstack[nps++] = atoi(&cmd[1])-1; 
			} else if ((cmd[0]=='-')&&(cmd[1]=='$')) {
				fstack[nfs++] = m_kill;
				pstack[nps++] = atoi(&cmd[2])-1;
			} else if (!strcmp("x",cmd)) { 
				fstack[nfs++] = m_position;
				pstack[nps++] = 0;
			} else if (!strcmp("t",cmd)) { 
				fstack[nfs++] = m_position;
				pstack[nps++] = 1;
			} else if ((j=hasvar(variables, cmd, nvars))>=0) {
				fstack[nfs++] = m_var;
				vstack[nvs++] = j;
			}
#include "init"
			else if (!strcmp("i",cmd)) { 
				fstack[nfs++] = m_const; cstack[ncs++] = 0.0; 
				fstack[nfs++] = m_const; cstack[ncs++] = 1.0; 
			}
			else if (!strcmp("pi",cmd)) { fstack[nfs++] = m_const; cstack[ncs++] = M_PI; }
			else if (!strcmp("pi_2",cmd)) { fstack[nfs++] = m_const; cstack[ncs++] = M_PI_2; }
			else if (!strcmp("pi_4",cmd)) { fstack[nfs++] = m_const; cstack[ncs++] = M_PI_4; }
			else if (!strcmp("1_pi",cmd)) { fstack[nfs++] = m_const; cstack[ncs++] = M_1_PI; }
			else if (!strcmp("e",cmd)) { fstack[nfs++] = m_const; cstack[ncs++] = M_E; }
			else if (!strcmp("sqrt2",cmd)) { fstack[nfs++] = m_const; cstack[ncs++] = M_SQRT2; }
			else if (!strcmp("sqrt1_2",cmd)) { fstack[nfs++] = m_const; cstack[ncs++] = M_SQRT1_2; }
			else if (!regexec(&regfloat,cmd,0,NULL,0)){
				fstack[nfs++] = m_const;
				cstack[ncs++] = atof(cmd);
			} else {
				fprintf(stderr,"error (%s): unknown command: %s\n",progname,cmd);
				return 1;
			}
		}
	}

	/* TEST FSTACK
	printf("f: "); for (i=0; i<nfs; i++) printf("%d ",fstack[i]); printf("\n");
	printf("c: "); for (i=0; i<ncs; i++) printf("%g ",cstack[i]); printf("\n");
	printf("p: "); for (i=0; i<nps; i++) printf("%d ",pstack[i]); printf("\n");
	printf("v: "); for (i=0; i<nvs; i++) printf("%d ",vstack[i]); printf("\n");
	printf("var: "); for (i=0; i<nvars; i++) printf("%s ",variables[i].name); printf("\n");
	*/

	while (mread(X,n,Bin,xflag,tflag,dt,&N)) {
		ncs=0; 
		nps=0; 
		nvs=0;
		for (i=0; i<nfs; i++) {
			switch (fstack[i]) {
			case m_assign:
				in = 1; out = 0; ckstacksize(N,in,out); x = &X[N-in];
				variables[vstack[nvs++]].val = x[0];
				for (j=0; j<out; j++) X[N-in+j] = y[j]; N += out-in;
				break;
			case m_position:
				in = 0; out = 1; ckstacksize(N,in,out); x = &X[N-in];
				y[0] = X[pstack[nps++]];
				for (j=0; j<out; j++) X[N-in+j] = y[j]; N += out-in;
				break;
			case m_clear:
				N=0;
				break;
			case m_roll:
				ckstacksize(N,2,3);
				in = (int)X[N-2];
				out = (int)X[N-1];
				ckstacksize(N,in+2,0);
				for (j=0; j<in; j++) {
					X[N-2-in+(out%in)]=X[N]; /* j... TODO hack on ... */
				}
				N -= 2;
				break;
			case m_kill:
				in = 0; out = 1; ckstacksize(N,in,out);
				for (j=pstack[nps++]; j<N-1; j++)
					X[j] = X[j+1];
				N--;
				break;
			case m_const:
				in = 0; out = 1; ckstacksize(N,in,out); x = &X[N-in];
				y[0] = cstack[ncs++];
				for (j=0; j<out; j++) X[N-in+j] = y[j]; N += out-in;
				break;
			case m_var:
				in = 0; out = 1; ckstacksize(N,in,out); x = &X[N-in];
				y[0] = variables[vstack[nvs++]].val;
				for (j=0; j<out; j++) X[N-in+j] = y[j]; N += out-in;
				break;
#include "cases"
			default:
				return 1;
			}
		}
		mwrite(X,N,Bout,format,tflag,u,F);
	}

	return 0;
}
