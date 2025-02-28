/16 bits and counting bills

white:0xffffff
black:0x000000
red:  0xff0000
green:0x00ff00
blue: 0x0000ff

pen:    {r[0x02fa;0,x,0],y,0x00}      /width;color
brush:  {0xfc020000,x,o}              /color
select: {r[0x012d;,x]}                /object-index
move:   {r[0x0214;|x]}                /x0,y0
line:   {r[0x0213;|x]}                /x1,y1
rect:   {r[0x041b;|x]}                /left,top,right,bottom
ellipse:{r[0x0418;|x]}                /same
polygon:{r[0x0324;(#x),,/+2^x,y]}     /x;y
polylin:{r[0x0325;(#x),,/+2^x,y]}     /x;y
text:   {r[0x0521;,#y],(y,(2!#y)#0x00),r["";|x]}    /x,y;str  
align:  {r[0x012e;,x]}                /2(right) 6(center) 8(bottom) 24(base)
font:   {0xfb02,(h x),(16#0x00),y,(32-#y)#0x00}     /height;name
h:      2#`c@,
r:      {(|x),,/2#'-4^`c y}
o:0x0000

wmf:{
 x,:,o                         /eof
 t:+/+/0xfafb=\*'x             /object table size
 x:(-4^`c 2+(#'x)%2),'x        /count each record
 n:(18+#x)%2                   /total
 m:(|/#'x)%2                   /max record size
 0x010009000003,(`c@,n),(h t),(`c@,m),o,,/x}


/example
wmf(pen[4;red]
    select 0
    rect 0 0 200 100
    font[50;"Courier New"]
    align 6+24
    text[100 60;_101+|0 20 15 6])
