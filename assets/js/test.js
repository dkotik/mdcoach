function randomize() {
  var canvas = document.querySelector('body > main');
  if (canvas.children.length == 0) return;
  for (var i = canvas.children.length; i >= 0; i--) {
    canvas.appendChild(canvas.children[(Math.random() * i) | 0]);
  }
}

function align() {
  var legends = document.querySelectorAll('body > main > section > legend:first-of-type');
  for (var i = 0; i < legends.length; i++) {
    legends[i].textContent = (i + 1).toString();
  }
}

function update() {
  document.querySelector('body').style.fontSize = document.querySelector('select[name=fontsize]').value;
  document.querySelector('main').style.columnCount = document.querySelector('select[name=columns]').value;
}

function eliminate(selector) {
  // node.parentNode.removeChild(node);
  var trash = document.querySelectorAll(selector);
  for (var i = 0; i < trash.length; i++) trash[i].parentNode.removeChild(trash[i]);
  align();
}
