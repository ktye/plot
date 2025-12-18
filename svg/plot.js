/* plot(x)
   plot(x,y)            ampang if (#y)=2*#x
   plot(x,y,x,y,..)
   plot(z,z,z,"polar")  no x if polar
   plot(..,"xy|polar|ampang|ring|foto|text")
   plot({Type:"polar",Lines:[..]},{Type:..})
   plot([{Type:..},{Type:..}])
   plot(..,"width",800,"height",400,"cols",2) */
   

let plot=(...a)=>{
 let p=[],d=[],w=800,h=600,c=0,ncolors,t="xy",fontratio=0.5455/*roughly font tahoma numbers*/,FA=x=>new Float64Array(x)
 let font1=16,font2=12,border=1,ticLength=6
 if(Array.isArray(a[0])&&a[0][0].Type)a=[...a[0],...a.slice(1)];//plot([{Type:..},{Type:..}],"width",..) => plot({Type:..},{Type:..},"width",..)
 for(let i=0;i<a.length;i++){let x=a[i];x.constructor==Float64Array?d.push(x):Array.isArray(x)?d.push(FA(x)):x.Type?p.push(x):"width"==x?w=a[++i]:"height"==x?h=a[++i]:"cols"==x?c=a[++i]:"font1"==x?font1=a[++i]:"font2"==x?font1=a[++i]:"string"==typeof x?t=a[++i]:0}
 let titleHeight=t=>t?2+Math.ceil(font1):2,xlabelHeight=l=>2+(l.length?font1:0),ylabelWidth=2+Math.ceil(font2)/*h for rotated y-axis*/,ticLabelWidth=yl=>max(...yl.map(x=>ceil(fontratio*font2*x.length))),ticLabelHeight=2+font2,rightXYWidth=l=>ceil(7+font2*fontratio*l.length/2)
 if(1==d.length)d.splice(0,0,FA(d[0].length).map((_,i)=>i))
 if(d.length){p={Lines:[]}
  if(t=="polar")d.forEach((z,i)=>p.Lines.push({Id:i,C:z}))
  else{for(let i=0;i<d.length;i+=2){let l={Id:i/2,x:d[i]};l[(d[1+i].length==2*d[i].length)?(t="ampang","C"):"Y"]=d[1+i];p.Lines.push(l)}}
  p.Type=t;p=[p]}

 let err=x=>{throw new Error(x)}
 let min=Math.min,max=Math.max,hypot=Math.hypot,sin=Math.sin,cos=Math.cos,atan2=Math.atan2,floor=Math.floor,ceil=Math.ceil,round=Math.round;const pi=Math.PI
 let Abs=x=>{let r=FA(x.length/2);for(let i=0;i<r.length;i++)r[i]=hypot(x[2*i],x[2*i+1]);return r}
 let ReIm=(x,o)=>{let r=FA(x.length/2);for(let i=o;i<r.length;i++)r[i]=x[2*i];return r},Real=x=>ReIm(x,0),Imag=x=>ReIm(x,1)
 let Ang=x=>{let r=FA(x.length/2);for(let i=0;i<r.length;i++)r[i]=atan2(x[2*i+1],x[2*i])*180/pi;return r}//-180,180
 let mima=a=>{let mi=Infinity,ma=-Infinity;a.forEach(x=>x.forEach(x=>(mi=min(mi,isNaN(x)?mi:x),ma=max(ma,isNaN(x)?ma:x))));return[mi,ma]}

 let scale=(x,x0,x1,y0,y1)=>y0+(x-x0)*(y1-y0)/(x1-x0)
 let axscale=(a,X,Y)=>([X.map(x=>(x=scale(x,a.xmin,a.xmax,0,10000),x<-10000?-10000:x>20000?20000:round(x))),Y.map(y=>(y=scale(y,a.ymax,a.ymin,0,10000),y<-10000?-10000:y>20000?20000:round(y)))])
 let nicenum=(ext,rnd)=>{let e=floor(Math.log10(ext)),f=ext/(10**e),r;return(rnd?((f<1.5)?1:(f<3)?2:(f<7)?5:10):((f<=1)?1:(f<=2)?2:(f<=5)?5:10))*10**e}
 let nicelim=(x,y)=>{let e=nicenum(y-x,false),s=nicenum(e/4,true);return[s*floor(x/s),s*ceil(y/s),s]}
 let ticstr=x=>{let s=String(x),t=x.toPrecision(5);return s.length<t.length?s:t}
 let nicetics=(x,y)=>{let [p,_,s]=nicelim(x,y),r=[],i=0;while(p+i*s<=y){if(p+i*s>=x)r.push(p+i*s);i++};return{Pos:r,Labels:r.map(ticstr)}}   
 let autoscale=a=>nicelim(...mima(a)),autoscalr=a=>{let[x,y]=mima(a);return nicelim(0,y)}
 let polarlimits=(p,ring)=>{let l=p.Limits;if(ring)console.log("todo ring-limits");let y0,y1=l.Ymax;if(p.Limits.Ymax<=0)[y0,y1]=autoscalr(p.Lines.map(l=>Abs(l.C)));[l.Xmin,l.Xmax,l.Ymin,l.Ymax]=[-y1,y1,-y1,y1];return l}
 let xxlimits=p=>{let l=p.Limits;if(l.Xmin==l.Xmax)[p.Limits.Xmin,p.Limits.Xmax]=autoscale(p.Lines.map(l=>l.X))}
 let xylimits=p=>{xxlimits(p);let l=p.Limits;if(l.Ymin==l.Ymax)[l.Ymin,l.Ymax]=autoscale(p.Lines.map(l=>l.Y));return l} //todo envelope/raster
 let aalimits=p=>{xxlimits(p);let l=p.Limits,x_;l.Ymin=0;if(l.Ymax<=0)[x_,l.Ymax]=autoscale(p.Lines.map(l=>Abs(l.C)));return l}
 let deflimits=l=>{if("undefined"==typeof l)l={};"Equal Xmin Xmax Ymin Ymax Zmin Zmax".split(" ").forEach(s=>{if(!(s in l))l[s]=0});return l}
 let limits=p=>{for(let i=0;i<p.length;i++){p[i].Limits=deflimits(p[i].Limits);let t=p[i].Type;p[i].Limits="xy"==t?xylimits(p[i]):"ampang"==t?aalimits(p[i]):"polar"==t?polarlimits(p[i],0):"ring"==t?polarlimits(p[i],1):{}};if(p[0].Limits.equal)console.log("todo equal-limits")}
 let labels=p=>{for(let i=0;i<p.length;i++){"Xlabel Ylabel Xunit Yunit".split(" ").forEach(x=>p[i][x]=x?x:"")}}
 let axes=(x,y,w,h,xmin,xmax,ymin,ymax)=>({x:x,y:y,w:w,h:h,xmin:xmin,xmax:xmax,ymin:ymin,ymax:ymax})
 let hs=s=>{const m={'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;',"'":'&#39;'};return s.replace(/[&<>"]/g,c=>m[c])}
 let xyer=l=>[l.X,l.Y,0/*todo env*/],xyamp=l=>[l.X,Abs(l.C),0],xyang=l=>[l.X,Ang(l.C),0],xypolar=l=>[Imag(l.C),Real(l.C),0]
 
 let text=(x,y,s,a,f2)=>`<text x="${x}" y="${y-2+(f2?font2:font1)*([0,0,0,0.5,1,1,1,0.5,0.5][a])}" class="${(f2?'s ':'')+('185'.includes(a)?'a1':'234'.includes(a)?'a2':'')}">${hs(s)}</text>`
 let vtext=(x,y,s)=>`<g transform="translate(${x} ${y}) rotate(270)"><text class="a1">${hs(s)}</text></g>`
 let line=(x1,y1,x2,y2)=>`<line x1="${x1}" y1="${y1}" x2="${x2}" y2="${y2}"/>`
 let drawTitle=(a,t,yo)=>t?text(a.x+a.w/2,a.y-ticLength-3-(yo?yo:0),t,1,0):""
 let drawXYTics=(a,xp,yp,xl,yl, l)=>(l=ticLength,line(a.x,a.y-l,a.x+a.w,a.y-l)+line(a.x,a.y+a.h+l,a.x+a.w,a.y+a.h+l)+line(a.x-l,a.y,a.x-l,a.y+a.h)+line(a.x+a.w+l,a.y,a.x+a.w+l,a.y+a.h) +htics(a,yp,yl,a.x-l,a.x)+htics(a,yp,[],a.x+a.w,a.x+a.w+l)+vtics(a,xp,[],a.y-l,a.y)+vtics(a,xp,xl,a.y+a.h,a.y+a.h+l))
 let htics=(a,Y,L,x1,x2)=>Y.map((y,i)=>(y=round(scale(y,a.ymax,a.ymin,a.y,a.y+a.h)),line(x1,y,x2,y)+(L.length?text(x1-3,y,L[i],3,1):""))).join("")
 let vtics=(a,X,L,y1,y2)=>X.map((x,i)=>(x=round(scale(x,a.xmin,a.xmax,a.x,a.x+a.w)),line(x,y1,x,y2)+(L.length?text(x,y2+2,L[i],5,1):""))).join("")
 let drawXlabel=(a,l,u)=>text(a.x+round(a.w/2),a.y+a.h+ticLength+ticLabelHeight,(l+" "+u).trim(),5,0)
 let drawYlabel=(a,l,u,ylw)=>vtext(a.x-2*ticLength-ylw,a.y+round(a.h/2),(l+" "+u).trim())
 let drawPolar=(a,rt, r,al)=>(r=floor(a.w/2),cx=a.x+r,cy=a.y+r,r1=r+ticLength/2,r2=r-ticLength/2,al=[1,0,0,7,6,6,5,4,4,3,2,2],p4=x=>x.toPrecision(4),
  Array(12).fill(0).map((_,i)=>30*i).map((p,i)=>{let cs=cos(p*pi/180),sn=sin(p*pi/180);return line(p4(cx+r1*cs),p4(cy+r1*sn),p4(cx+r2*cs),p4(cy+r2*sn))+text(p4(cx+r1*cs),p4(cy+r1*sn),((90+p)%360)+"",al[(3+i)%12],1)}).join("")
  +rt.map(R=>`<circle cx="${cx}" cy="${cy}" r="${R/a.ymax*r}" stroke-width="1" stroke="black" fill="none"/>`).join("")+line(cx-r,cy,cx+r,cy)+line(cx,cy-r,cx,cy+r)+`<circle cx="${cx}" cy="${cy}" r="${r}" stroke-width="2" stroke="black" fill="none"/>`)
 let linestyle=(p,l,i)=>{let lw=l?.Style?.Line?.Width?l.Style.Line.Width:0,ps=l?.Style?.Marker?.Size?l.Style.Marker.Size:0;[lw,ps]=(!(lw||ps))?(p.Type=="polar"?[0,3]:[2,0]):[lw,ps]; console.log("ps",ps); return[lw,ps,l?.Style?.Color?l.Style.Color:l?.Id?1+l.Id:1+i]}
 let lineclass=(lw,c)=>`class="c${1+(c-1%ncolors)}"`+(2!=lw?`stroke-width="${lw}"`:"")
 let drawLines=(a,p,f)=>`<g transform="translate(${a.x} ${a.y}) scale(${a.w/10000} ${a.h/10000})">`+p.Lines.map((l,i)=>/*todo l.Style.Marker.Marker=="bar"*/drawLine(a,p,l,i,f)).join("")+`</g>`
 let scalepoint=(ps,w)=>round(10000*ps/w)
 let drawLine=(a,p,l,i,f)=>{let[lw,ps,c]=linestyle(p,l,i),r="",[x,y]=axscale(a,...f(l));x=Array.from(x);/*todo labels,endmarks*/if(lw>0&&x.length)r+=`<path d="M${x[0]+' '+y[0]}`+x.map((x,i)=>`L${x+' '+y[i]}`).join("")+`" ${lineclass(lw,c)}/>`
  if(ps)r+=`<g class="C${1+(c-1%ncolors)}">`+x.map((x,i)=>`<circle cx="${x}" cy="${y[i]}" r="${scalepoint(ps,a.w)}"/>`).join("")+`</g>`
  return r}  //todo Style.Line.Arrow Style.Line.EndMarks
  
 let empty=(p,w,h)=>""
 let xy=(p,w,h)=>{let xt=nicetics(p.Limits.Xmin,p.Limits.Xmax),yt=nicetics(p.Limits.Ymin,p.Limits.Ymax),ylw=ticLabelWidth(yt.Labels);
  let hfix=2*border+3*ticLength+ylabelWidth+ylw+rightXYWidth(xt.Labels.length?xt.Labels[xt.Labels.length-1]:"")
  let vfix=2*border+2*ticLength+titleHeight(p.Title)+ticLabelHeight+xlabelHeight(p.Xlabel+p.Xunit);
  let hs=w-hfix,vs=h-vfix,x0=0,y0=0;if(vs>2*hs){y0=floor((vs-2*hs)/2);vs=2*hs;};
  x0+=ylabelWidth+ylw+2*ticLength+border;y0+=titleHeight(p.Title)+ticLength+border;
  let ax=axes(x0,y0,hs,vs,p.Limits.Xmin,p.Limits.Xmax,p.Limits.Ymin,p.Limits.Ymax);
  console.log("todo drawxy");return drawLines(ax,p,xyer)+drawXYTics(ax,xt.Pos,yt.Pos,xt.Labels,yt.Labels)+drawTitle(ax,p.Title)+drawXlabel(ax,p.Xlabel,p.Xunit)+drawYlabel(ax,p.Ylabel,p.Yunit,ylw)}
 let polar=(p,w,h)=>{let rt=nicetics(0,p.Limits.Ymax),ylw=ticLabelWidth(["270"]); console.log("limits",p.Limits,"nt",nicetics(0,p.Limits.Ymax) ,nicetics(p.Limits.Ymin,p.Limits.Ymax) );
  let hfix=2*border+2*ylw
  let vfix=2*border+titleHeight(p.Title)+2*ticLabelHeight
  let hs=w-hfix,vs=h-vfix,d=hs<0&&vs<0?0:hs<vs?hs:vs; console.log("h",h,"vfix",vfix,"vs",vs);  d-=1-(1&d);  console.log("d!",d);    if(d<0)return"";
  let x0=floor((w-hfix-d)/2),y0=floor((h-vfix-d)/2),ax=axes(x0+ylw+border,y0+titleHeight(p.Title)+ticLabelHeight+border,d,d,p.Limits.Xmin,p.Limits.Xmax,p.Limits.Ymin,p.Limits.Ymax);
  return drawLines(ax,p,xypolar)+drawTitle(ax,p.Title,ticLabelHeight-ticLength)+drawPolar(ax,rt.Pos)}
 let ring=(p,w,h)=>""
 let ampang=(p,w,h)=>{let xt=nicetics(p.Limits.Xmin,p.Limits.Xmax),yt=nicetics(p.Limits.Ymin,p.Limits.Ymax),ylw=ticLabelWidth(yt.Labels);
  let hfix=2*border+3*ticLength+ylabelWidth+ylw+rightXYWidth(xt.Labels.length?xt.Labels[xt.Labels.length-1]:"")
  let vfix=2*border+4*ticLength+titleHeight(p.Title)+ticLabelHeight+xlabelHeight(p.Xlabel+p.Xunit)
  let x0=0,y0=0,hs=w-hfix,vs=h-vfix;  //if(hs>1.5*vs){x0=floor(hs-(1.5*vs)/2);hs=ceil(1.5*vs)}else if(vs>2*hs){y0=floor(hs-(2*vs)/2);vs=2*hs}
  let aw=hs,h1=ceil(2*vs/3),h2=vs-h1;
  x0+=ylabelWidth+ylw+2*ticLength+border;y0+=titleHeight(p.Title)+ticLength+border;
  let amp=axes(x0,y0,hs,h1,p.Limits.Xmin,p.Limits.Xmax,p.Limits.Ymin,p.Limits.Ymax)
  let ang=axes(x0,y0+h1+2*ticLength,hs,h2,p.Limits.Xmin,p.Limits.Xmax,-180,180),angs="-180 -90 0 90 180".split(" ")
  return drawLines(amp,p,xyamp)+drawLines(ang,p,xyang)+drawXYTics(amp,xt.Pos,yt.Pos,[],yt.Labels)+drawXYTics(ang,xt.Pos,angs.map(Number),xt.Labels,angs)+drawTitle(amp,p.Title)+drawXlabel(ang,p.Xlabel,p.Xunit)+drawYlabel(amp,p.Ylabel,p.Yunit,ylw)
 }
 let foto=(p,w,h)=>""
 let textplot=(p,w,h)=>""
 
 let grid=(n,w,h,c, g)=>{g={n:n}
  c<0?(g.colmajor=1,-c):(!c)?c=((n<13)?[4,4,4,4,4,3,3,4,4,5,5,4,4][n]:5):0
  g.r=1;g.c=(n<c?n:(g.r=0|n/c,c));
  g.r=(g.r*g.c<n)?1+g.r:g.r;
  g.w=w/g.c;g.h=h/g.r;g.width=w;return g}
 let xyi=(g,n, i,k,x,y,m)=>{x=0;i=0|n/g.c;k=n%g.c;if(g.colmajor)[i,k]=[k,i]             //{n: 2, r: 1, c: 2, w: 600, h: 500}
  if(i==0|(g.n-1)/g.c){m=1+((g.n-1)%g.c);x=(g.width-m*g.w)/2}
  x+=k*g.w;y=i*g.h;return[x,y]}
 let plots=(p,w,h,c)=>{let colors=p.length?(p[0]?.Style?.Order?p[0].Style.Order.split(","):[]):[];colors=(colors.length?colors:"#003FFF,#03ED3A,#E8000B,#8A2BE2,#FFC400,#00D7FF".split(","));ncolors=colors.length;
  limits(p);labels(p);let g=grid(p.length,w,h,c),r=`<svg viewBox="0 0 ${w} ${h}" width="${w}" height="${h}" xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink">
  <style>text{font-family:Tahoma,sans-serif;font-size:${font1}px}.a1{text-anchor:middle}.a2{text-anchor:end}.s{font-size:${font2}px}.v{writing-mode:sideways-lr;font-size:${font1}px}
  ${colors.map((x,i)=>`.c${1+i}{stroke:${x}}`).join("")}${colors.map((x,i)=>`.C${1+i}{fill:${x}}`).join("")}
  line{stroke:black}
  path{fill:none;stroke:black;stroke-width:2;vector-effect:non-scaling-stroke}
  </style>
  <clipPath id="B"><rect width="${g.w}" height="${g.h}"/></clipPath>`
  let P={"":empty,"xy":xy,"raster":xy,"polar":polar,"ring":ring,"ampang":ampang,"foto":foto,"text":textplot}
  p.forEach((p,i, x,y)=>{[x,y]=xyi(g,i);r+=`<g transform="translate(${x+0.5},${y+0.5})" clip-path="url(#B)">`+((p.Type in P)?(P[p.Type](p,g.w,g.h)):err("no such plot type:"+p.Type))+"</g>"});return r+"</svg>"} 
 return plots(p,w,h,c)}
