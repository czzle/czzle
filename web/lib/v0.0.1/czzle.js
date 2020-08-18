(()=>{class l{constructor(t){this.innerHTML=t}}class h{constructor(t){this.uid=t}create(t,{on:e,attrs:i,style:n,children:r}={}){if(i)for(const a in i){let s=i[a];(a==="id"||a==="class"||a==="for")&&(s=this.prefixAttr(s),i[a]=s)}return{tag:t,style:n||{},on:e||{},attrs:i||{},children:r||[]}}prefixAttr(t){const e=t.toString().split(" ");for(let i=0;i<e.length;i++)e[i]=`${this.uid}-${e[i]}`;return e.join(" ")}renderElement(t){const e=document.createElement(t.tag);for(const n in t.attrs){const r=t.attrs[n];e.setAttribute(n,r.toString())}for(const n of t.children){const r=this.render(n);e.appendChild(r)}let i=[];for(const n in t.style){let r=t.style[n];typeof r=="number"&&(r=`${r}px`),i.push(`${n}: ${t.style[n]}`)}if(i.length&&e.setAttribute("style",i.join("; ")),t.on.click){const n=t.on.click;e.onclick=r=>{n({x:r.x,y:r.y})}}if(t.on.change){const n=t.on.change,r=e.getAttribute("type");r&&(e.onchange=a=>{r==="checkbox"?n({value:e.checked}):n({value:e.value})})}return e}renderText(t){return document.createTextNode(t)}renderHTML(t){const e=document.createElement("template"),i=t.innerHTML.trim();return e.innerHTML=i,e.content.firstChild}render(t){return typeof t=="string"?this.renderText(t):t instanceof l?this.renderHTML(t):this.renderElement(t)}zip(t,e){const i=[];for(let n=0;n<Math.max(t.length,e.length);n++)i.push([t[n],e[n]]);return i}diffAttrs(t,e){const i=[];for(const n in e){const r=e[n];i.push(a=>(a.setAttribute(n,r.toString()),a))}for(const n in t)n in e||i.push(r=>(r.removeAttribute(n),r));return n=>{for(const r of i)r(n)}}diffStyle(t,e){return i=>{let n=[];for(const r in e){let a=e[r];typeof a=="number"&&(a=`${a}px`),n.push(`${r}: ${e[r]}`)}n.length&&i.setAttribute("style",n.join("; "))}}diffEvents(t,e){const i=[];if(e.click){const n=e.click;i.push(r=>(r.onclick=a=>{n({x:a.x,y:a.y})},r))}else t.click&&i.push(n=>{n.onclick=null});if(e.change){const n=e.change;i.push(r=>{const a=r.getAttribute("type");a&&(r.onchange=s=>{a==="checkbox"?n({value:r.checked}):n({value:r.value})})})}else t.change&&i.push(n=>{n.onchange=null});return n=>{for(const r of i)r(n)}}diffChildren(t,e){const i=[];t.forEach((r,a)=>{i.push(this.diff(r,e[a]))});const n=[];for(const r of e.slice(t.length))n.push(a=>(a.appendChild(this.render(r)),a));return r=>{for(const[a,s]of this.zip(i,r.childNodes))a(s);for(const a of n)a(r);return r}}diff(t,e){if(e===void 0)return s=>{s.remove();return};if(typeof t=="string"||typeof e=="string")return t!==e?s=>{const c=this.render(e);return s.replaceWith(c),c}:s=>{};if(t instanceof l&&e instanceof l)return t.innerHTML!==e.innerHTML?s=>{const c=this.render(e);return s.replaceWith(c),c}:s=>{};if(t instanceof l||e instanceof l)return s=>{const c=this.render(e);return s.replaceWith(c),c};if(t.tag!==e.tag)return s=>{const c=this.render(e);return s.replaceWith(c),c};const i=this.diffAttrs(t.attrs,e.attrs),n=this.diffStyle(t.style,e.style),r=this.diffEvents(t.on,e.on),a=this.diffChildren(t.children,e.children);return s=>(i(s),n(s),a(s),r(s),s)}mount(t,e){return e.replaceWith(t),t}bootstrap(t,e,i=!0){let n=e.element(),r=document.createElement("div");return i?t.replaceWith(r):t.append(r),r=this.mount(this.render(n),r),e.registerChange&&e.registerChange(()=>{const a=e.element(),s=this.diff(n,a);r=s(r),n=a}),r}}const o=new h("czzle");function g(t){document.querySelectorAll('input[type="hidden"][czzle]').forEach(e=>{let i;e.hasAttribute("name")?i=e.getAttribute("name"):i="czzle",o.bootstrap(e,new m(t,i))}),o.bootstrap(document.head,new d(),!1),window.addEventListener("load",()=>{console.warn("loaded",t)})}var f={configure:g};let y=0;class m{constructor(t,e){this.cfg=t;this.name=e;this.id=y++;this.checked=!1;this.value="";this.onChange=()=>{}}element(){return o.create("div",{style:{border:"1px solid rgb(184, 184, 184)","max-width":"200px",padding:"16px","border-radius":"3px","background-color":"#fafafa",display:"flex","flex-direction":"row","align-items":"center","justify-content":"start",margin:"10px 0",cursor:this.checked?"default":"pointer","-webkit-touch-callout":"none","-webkit-user-select":"none","-khtml-user-select":"none","-moz-user-select":"none","-ms-user-select":"none","user-select":"none"},children:[o.create("input",{attrs:{type:"hidden",name:this.name,value:this.value}}),o.create("div",{style:{"margin-right":"5px",height:"24px",width:"24px"},children:[this.checked?new l('<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path d="M0 0h24v24H0z" fill="none"/><path fill="#7b8794" d="M19 3H5c-1.11 0-2 .9-2 2v14c0 1.1.89 2 2 2h14c1.11 0 2-.9 2-2V5c0-1.1-.89-2-2-2zm-9 14l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/></svg>'):new l('<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24"><path fill="#9e9e9e" d="M19 5v14H5V5h14m0-2H5c-1.1 0-2 .9-2 2v14c0 1.1.9 2 2 2h14c1.1 0 2-.9 2-2V5c0-1.1-.9-2-2-2z"/><path d="M0 0h24v24H0z" fill="none"/> </svg>')]}),o.create("span",{style:{height:"24px","line-height":"24px","font-size":"18px"},children:["I'm not a robot"]})],on:{click:()=>{if(this.value)return;let t=[];fetch(`${this.cfg.apiUrl}/v1/begin`,{method:"POST",body:JSON.stringify({}),headers:{"Content-Type":"application/json; charset=UTF-8"}}).then(e=>e.ok?e.json():Promise.reject(e)).then(e=>{const i=o.bootstrap(document.body,new x(e.puzzle,n=>{n&&(this.checked=!0,this.value=n,this.onChange()),i.remove()}),!1)})}}})}registerChange(t){this.onChange=t}}class d{element(){return o.create("style",{children:[b]})}}var u;(function(t){t.Unknown="unknown",t.None="none",t.Easy="easy",t.Medium="medium",t.Hard="hard"})(u||(u={}));var p;(function(t){t.Front="front",t.Back="back"})(p||(p={}));const b=`
    @keyframes czzle-puzzle-host{
        from {
            opacity: 0;
        }
        to {
            opacity: 1;
        }
    }
    @keyframes czzle-flip-1 {
        from {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(0deg);
        }
        to {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(90deg);
        }
    }
    @keyframes czzle-flip-2 {
        from {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(90deg);
        }
        to {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(0deg);
        }
    }
    @keyframes czzle-flip-back-1 {
        from {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(0deg);
        }
        to {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(-90deg);
        }
    }
    @keyframes czzle-flip-back-2 {
        from {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(-90deg);
        }
        to {
            border: none;
            margin: 0;
            transform-style: preserve-3d;
            transform: rotateY(0deg);
        }
    }
`;class x{constructor(t,e){this.puzzle=t;this.slovecb=e;this.grid=new v(this.puzzle.tile_map)}element(){return o.create("div",{style:{position:"fixed",display:"flex","flex-direction":"row","align-items":"center","justify-content":"center",top:0,left:0,width:"100vw",height:"100vh","background-color":"rgba(0, 0, 0, 0.5)","z-index":99999999,"animation-name":"czzle-puzzle-host","animation-duration":"0.5s"},on:{click:()=>{}},children:[o.create("div",{style:{padding:"20px","background-color":"#fefefe"},children:[this.grid.element()]})]})}registerChange(t){this.onChange=t,this.grid.registerChange(t)}}class v{constructor(t){this.m=t;this.animate=[];this.animateBack=[];this.animatePrev=[];this.animatePrevBack=[]}element(){const t=this.m.tiles.filter(e=>this.selected?e.type===p.Front?e.pos.x!==this.selected.x||e.pos.y!==this.selected.y:e.pos.x===this.selected.x&&e.pos.y===this.selected.y:e.type===p.Front);return o.create("div",{style:{height:"300px",width:"300px",position:"relative"},children:t.map(e=>{let i={};return this.animate.find(n=>n.x===e.pos.x&&n.y===e.pos.y)&&(i={"animation-name":"czzle-flip-1","animation-duration":"0.25s"}),this.animateBack.find(n=>n.x===e.pos.x&&n.y===e.pos.y)&&(i={"animation-name":"czzle-flip-2","animation-duration":"0.25s"}),this.animatePrev.find(n=>n.x===e.pos.x&&n.y===e.pos.y)&&(i={"animation-name":"czzle-flip-back-1","animation-duration":"0.25s"}),this.animatePrevBack.find(n=>n.x===e.pos.x&&n.y===e.pos.y)&&(i={"animation-name":"czzle-flip-back-2","animation-duration":"0.25s"}),this.selected&&(this.selected.x===e.pos.x&&this.selected.y===e.pos.y&&(i={border:"5px solid rgb(23, 100, 163)",margin:"-5px 0 0 -5px",...i})),o.create("img",{attrs:{src:`data:image/png;base64,${e.data}`},style:{position:"absolute",left:`${e.pos.x*100}px`,top:`${e.pos.y*100}px`,border:"none",width:"100px",height:"100px",outline:"none",cursor:"pointer",...i},on:{click:()=>{if(this.animate.length||this.animateBack.length||this.animatePrev.length||this.animatePrevBack.length)return;let n;this.selected?(this.animatePrev=[this.selected],this.selected.x===e.pos.x&&this.selected.y===e.pos.y?n=null:(n=e.pos,this.animate=[e.pos])):(n=e.pos,this.animate=[e.pos]),setTimeout(()=>{this.selected=n,this.animateBack=this.animate,this.animate=[],this.animatePrevBack=this.animatePrev,this.animatePrev=[],setTimeout(()=>{this.animateBack=[],this.animatePrevBack=[],this.onChange&&this.onChange()},250),this.onChange&&this.onChange()},250),this.onChange&&this.onChange()}}})})})}registerChange(t){this.onChange=t}}Object.defineProperty(window,"czzle",{get(){return f}});})();
