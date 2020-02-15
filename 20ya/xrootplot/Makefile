PROG=	xrootplot
SRCS=	xrootplot.c plot_xrootplot.c

CFLAGS+=	-I/usr/X11R6/include -I/usr/local/include
#LDADD+= -L/usr/X11R6/lib -lX11 -lXpm -L/usr/local/lib -ljpeg -lpng -ltiff -lz -lgif -lImlib
LDADD+= -L/usr/X11R6/lib -lX11 -lXpm -L/usr/local/lib -ljpeg -lpng -ltiff -lz -lImlib

.include <$(HOME)/.elmar.prog.mk>
