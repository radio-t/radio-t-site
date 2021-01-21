(() => {
  const url = require('../images/icons-sprite.svg')
  const xhr = new XMLHttpRequest();
  const container = document.createElement("div");
  const insert = () => document.body.insertBefore(container, document.body.childNodes[0]);
  
  container.className = 'd-none'
  xhr.open("GET", url, true);
  xhr.send();
  xhr.onload = () => {
    container.innerHTML = xhr.responseText;
    
    if (document.body) {
      insert();
      return;
    }

    document.onreadystatechange = function () {
      if (document.readyState === "interactive") {
        insert();
      }
  }
  }
})()