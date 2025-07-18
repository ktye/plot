/* plot(x)
   plot(x,y)            ampang if (#y)=2*#x
   plot(x,y,x,y,..)
   plot(z,z,z,"polar")  no x if polar
   plot(..,"xy|polar|ampang|ring|foto|text")
   plot({Type:"polar",Lines:[..]},{Type:..})
   plot([{Type:..},{Type:..}])
   plot(..,"width",800,"height",400,"cols":2) */

let plot=(...a)=>{
 let p=[],d=[],w=800,h=600,c=0,t="xy"
 for(let i=0;i<a.length;i++){let x=a[i];x.constructor==Float64Array?d.push(x):Array.isArray(x)?d.push(new Float64Array(x)):x.Type?p.push(x):"width"==x?w=a[++i]:"height"==x?h=a[++i]:"cols"==x?c=a[++i]:"string"==x?t=a[++i]:0}
 if(1==d.length)d.splice(0,0,new Float64Array(d[0].length).map((_,i)=>i))
 if(d.length){p={Lines:[]}
  if(t=="polar")d.forEach((z,i)=>p.Lines.push({Id:i,C:z}))
  else{for(let i=0;i<d.length;i+=2){let l={Id:i/2,x:d[i]};l[(d[1+i].length==2*d[i].length)?(t="ampang","C"):"Y"]=d[1+i];p.Lines.push(l)}}
  p.Type=t;p=[p]}

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

 return plots(p,w,h,c)}
