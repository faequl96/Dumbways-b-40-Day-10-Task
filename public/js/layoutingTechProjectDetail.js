let tech1Img = document.querySelector(".techno .techGrid .node img");
let tech1 = document.querySelector(".techno .techGrid .node");
// let tech1Text = document.querySelector(".techno .techGrid .node p");

if(tech1Img.src == "http://localhost:8000/public/img/node.png") {
    tech1.style.display = "flex";
} else {
    tech1.style.display = "none";
}


let tech2Img = document.querySelector(".techno .techGrid .vue img");
let tech2 = document.querySelector(".techno .techGrid .vue");

if(tech2Img.src == "http://localhost:8000/public/img/vue.png") {
    tech2.style.display = "flex";
} else {
    tech2.style.display = "none";
}

let tech3Img = document.querySelector(".techno .techGrid .react img");
let tech3 = document.querySelector(".techno .techGrid .react");

if(tech3Img.src == "http://localhost:8000/public/img/react.png") {
    tech3.style.display = "flex";
} else {
    tech3.style.display = "none";
}


let tech4Img = document.querySelector(".techno .techGrid .typescript img");
let tech4 = document.querySelector(".techno .techGrid .typescript");

if(tech4Img.src == "http://localhost:8000/public/img/typescript.png") {
    tech4.style.display = "flex";
} else {
    tech4.style.display = "none";
}