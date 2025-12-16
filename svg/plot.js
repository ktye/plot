/* plot(x)
   plot(x,y)            ampang if (#y)=2*#x
   plot(x,y,x,y,..)
   plot(z,z,z,"polar")  no x if polar
   plot(..,"xy|polar|ampang|ring|foto|text")
   plot({Type:"polar",Lines:[..]},{Type:..})
   plot([{Type:..},{Type:..}])
   plot(..,"width",800,"height",400,"cols",2) */
   

let plot=(...a)=>{
 let p=[],d=[],w=800,h=600,c=0,t="xy",FA=x=>new Float64Array(x)
 let font1=20,font2=16,border=5,ticLength=6
 if(Array.isArray(a[0])&&a[0][0].Type)a=[...a[0],...a.slice(1)];//plot([{Type:..},{Type:..}],"width",..) => plot({Type:..},{Type:..},"width",..)
 for(let i=0;i<a.length;i++){let x=a[i];x.constructor==Float64Array?d.push(x):Array.isArray(x)?d.push(FA(x)):x.Type?p.push(x):"width"==x?w=a[++i]:"height"==x?h=a[++i]:"cols"==x?c=a[++i]:"font1"==x?font1=a[++i]:"font2"==x?font1=a[++i]:"string"==typeof x?t=a[++i]:0}
 let titleHeight=t=>t?2+Math.ceil(font1):2,xlabelHeight=(l,u)=>2+(l||u)?font1:0,ylabelWidth=2+Math.ceil(font2)/*h for rotated y-axis*/,ticLabelWidth=yl=>max(...yl.map(x=>font2*x.length/2)),ticLabelHeight=2+font1,rightXYWidth=l=>7+font2*l.length/4
 if(1==d.length)d.splice(0,0,FA(d[0].length).map((_,i)=>i))
 if(d.length){p={Lines:[]}
  if(t=="polar")d.forEach((z,i)=>p.Lines.push({Id:i,C:z}))
  else{for(let i=0;i<d.length;i+=2){let l={Id:i/2,x:d[i]};l[(d[1+i].length==2*d[i].length)?(t="ampang","C"):"Y"]=d[1+i];p.Lines.push(l)}}
  p.Type=t;p=[p]}

 let err=x=>{throw new Error(x)}
 let min=Math.min,max=Math.max,hypot=Math.hypot,atan2=Math.atan2,floor=Math.floor,ceil=Math.ceil,round=Math.round
 let Abs=x=>{let r=FA(x.length/2);for(let i=0;i<r.length;i++)r[i]=hypot(r[2*i],r[2*i+1]);return r}
 let mima=a=>{let mi=Infinity,ma=-Infinity;a.forEach(x=>x.forEach(x=>(mi=min(mi,isNaN(x)?mi:x),ma=max(ma,isNaN(x)?ma:x))));return[mi,ma]}

 let nicenum=(ext,rnd)=>{let e=floor(Math.log10(ext)),f=ext/(10**e),r;return(rnd?((f<1.5)?1:(f<3)?2:(f<7)?5:10):((f<=1)?1:(f<=2)?2:(f<=5)?5:10))*10**e}
 let nicelim=(x,y)=>{let e=nicenum(y-x,false),s=nicenum(e/4,true);return[s*floor(x/s),s*ceil(y/s),s]}
 let ticstr=x=>{let s=String(x),t=x.toPrecision(5);return s.length<t.length?s:t}
 let nicetics=(x,y)=>{let [p,_,s]=nicelim(x,y),r=[],i=0;while(p+i*s<=y){if(p+i*s>=x)r.push(p+i*s);i++};return{Pos:r,Labels:r.map(ticstr)}}   
 let autoscale=a=>nicelim(...mima(a));
 let polarlimits=(p,ring)=>{console.log("polarlimits",p);let l=p.Limits;if(ring)console.log("todo ring-limits");let y0,y1;if(p.Limits.Ymax==0)[y0,y1]=autoscale(p.Lines.map(Abs));[l.Xmin,l.Xmax,l.Ymin,l.Ymax]=[-y1,y1,-y1,y1];return l}
 let xylimits=p=>{let l=p.Limits;if(l.Xmin==l.Xmax)[l.Xmin,l.Xmax]=autoscale(p.Lines.map(l=>l.X));if(l.Ymin==l.Ymax)[l.Ymin,l.Ymax]=autoscale(p.Lines.map(l=>l.Y));return l} //todo envelope/raster
 let deflimits=l=>{if("undefined"==typeof l)l={};"Equal Xmin Xmax Ymin Ymax Zmin Zmax".split(" ").forEach(s=>{if(!(s in l))l[s]=0});return l}
 let limits=p=>{for(let i=0;i<p.length;i++){p[i].Limits=deflimits(p[i].Limits);let t=p[i].Type;p[i].Limits="xy"==t?xylimits(p[i]):"ampang"==t?aalimits(p[i]):"polar"==t?polarlimits(p[i],0):"ring"==t?polarlimits(p[i],1):{}};if(p[0].Limits.equal)console.log("todo equal-limits")}
 let axes=(x,y,w,h)=>({x:x,y:y,w:w,h:h})
 let hs=s=>{const m={'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'};return s.replace(/[&<>"]/g,c=>m[c])}
 
 let text=(x,y,s,al,ver,fs)=>`<text x="${x}" y="${y-font1*([0,0,0,0.5,1,1,1,0.5,0.5][al])}" text-anchor="${['start','middle','end'][0,1,2,2,2,1,0,0,1][al]}" font-size="${fs}"px >${hs(t)}</text>`
 let drawTitle=(a,t,vs)=>t?text(a.x+a.w/2,a.y-vs-1,t,1,0,font1):""
 
 let empty=(p,w,h)=>""
 let xy=(p,w,h)=>{let r="",xt=nicetics(p.Limits.Xmin,p.Limits.Xmax),yt=nicetics(p.Limits.Ymin,p.Limits.Ymax);
  let hfix=2*border+3*ticLength+ylabelWidth+ticLabelWidth(yt.Labels)+rightXYWidth(yt.Labels.length?yt.Labels[yt.Labels.length-1]:"")
  let vfix=2*border+2*ticLength+titleHeight(p.Title)+ticLabelHeight+xlabelHeight(p.Xlabel,p.Unit)
  let hs=w-hfix,vs=h-vfix,x0=0,y0=0;if(vs>2*hs){y0=(vs-2*hs)/2;vs=2*hs;};x0+=ylabelWidth+ticLabelWidth(yt.Labels)+2*ticLength+border;y0+=titleHeight(p.Title)+ticLength+border;
  let ax=axes(x0,y0,hs,vs); //limits?
  console.log("todo drawxy");
  //drawXYTics(ax,xt,yt);
  r+=drawTitle(ax,p.Title,ticLength);
  
  console.log("x0/y0",x0,y0,"w/h",hs,vs,"h",h,"vfix",vfix);
  return r+`<rect x="${x0}" y="${y0}" width="${hs}" height="${vs}" stroke="black" fill="none"/><rect width="${w}" height="${h}" stroke-width="4" stroke="black" fill="none"/>`}
 let polar=(p,w,h)=>(console.log("polar:",p),`<circle cx="${w/2}" cy="${h/2}" r=${min(w,h)/2} stroke-width="4" stroke="black" fill="none"/>`)
 let ring=(p,w,h)=>""
 let ampang=(p,w,h)=>""
 let foto=(p,w,h)=>""
 let textplot=(p,w,h)=>""
 
 
 let grid=(n,w,h,c, g)=>{g={n:n}
  c<0?(g.colmajor=1,-c):(!c)?c=((n<13)?[4,4,4,4,4,3,3,4,4,5,5,4,4][n]:5):0
  g.r=1;g.c=(n<c?n:(g.r=0|n/c,c));
  console.log("g",g)
  g.r=(g.r*g.c<n)?1+g.r:g.r;
  g.w=w/g.c;g.h=h/g.r;g.width=w;return g}
 let xyi=(g,n, i,k,x,y,m)=>{x=0;i=0|n/g.c;k=n%g.c;if(g.colmajor)[i,k]=[k,i]             //{n: 2, r: 1, c: 2, w: 600, h: 500}
  if(i==0|(g.n-1)/g.c){m=1+((g.n-1)%g.c);x=(g.width-m*g.w)/2}
  x+=k*g.w;y=i*g.h; console.log("xyi",g,n,i,k,x,y); return[x,y]}
 let plots=(p,w,h,c)=>{ 
  limits(p);let g=grid(p.length,w,h,c),r=`<svg viewBox="0 0 ${w} ${h}" width="${w}" height="${h}" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink"><clipPath id="B"><rect width="${g.w}" height="${g.h}"/></clipPath>`
  let P={"":empty,"xy":xy,"raster":xy,"polar":polar,"ring":ring,"ampang":ampang,"foto":foto,"text":textplot}
  p.forEach((p,i, x,y)=>{[x,y]=xyi(g,i);r+=`<g transform="translate(${x},${y})" clip-path="url(#B)">`+((p.Type in P)?(console.log("p.Type",p.Type,P[P.Type]),P[p.Type](p,g.w,g.h)):err("no such plot type:"+p.Type))+"</g>"});return r+"</svg>"} 
 return plots(p,w,h,c)}
