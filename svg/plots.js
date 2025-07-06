let plots=(p,w,h,c, g,r,P)=>{g=grid(p.length,w,h,c);console.log("g",g)
 r=`<svg viewBox="0 0 ${w} ${h}" width="${w}" height="${h}" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">`
 r+=`<clipPath id="B"><rect width="${g.w}" height="${g.h}"/></clipPath>`
 P={"":empty,"XY":xy,"Raster":xy,"Polar":polar,"Ring":ring,"AmpAng":ampang,"Foto":foto,"Text":text}
 p.forEach((p,i, x,y)=>{[x,y]=xyi(g,i);r+=`<g transform="translate(${x},${y})" clip-path="url(#B)">`+((p.Type in P)?P[p.Type](p,g.w,g.h):err("no such plot type:"+p.type))+"</g>"});return r+"</svg>"}
let err=x=>{throw new Error(x)}
let grid=(n,w,h,c, g)=>{g={n:n}
 c<0?(g.colmajor=1,-c):(!c)?c=((n<13)?[4,4,4,4,4,3,3,4,4,5,5,4,4][n]:5):0
 g.r=1;g.c=(n<c?n:(g.r=0|n/c,c));
 console.log("g",g)
 g.r=(g.r*g.c<n)?1+g.r:g.r;
 [g.w,g.h]=[w/g.c,h/g.r];return g}
let xyi=(g,n, i,k,x,y,m)=>{x=0;i=0|n/g.c;k=n%g.c;if(g.colmajor)[i,k]=[k,i]
 if(i==0|(g.n-1)/g.c){m=1+((g.n-1)%g.c);x=(g.w-m*g.w)/2}
 x+=k*g.w;y=i*g.h;return[x,y]}

let empty=(p,w,h)=>""
let xy=(p,w,h)=>""
let polar=(p,w,h)=>""
let ring=(p,w,h)=>""
let ampang=(p,w,h)=>""
let foto=(p,w,h)=>""
let text=(p,w,h)=>""
