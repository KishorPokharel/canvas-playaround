function getMousePos(canvas, evt) {
  const rect = canvas.getBoundingClientRect();
  return {
    x: evt.clientX - rect.left,
    y: evt.clientY - rect.top,
  };
}

const canvas = document.getElementById('canvas');
const canvas2 = document.getElementById('canvas2');
const saveCanvasBtn = document.querySelector('.saveCanvasBtn');

console.assert(canvas != null);
console.assert(saveCanvasBtn != null);

const ctx = canvas.getContext('2d');
const ctx2 = canvas2.getContext('2d');
canvas.addEventListener('mousemove', e => {
  const {x, y} = getMousePos(canvas, e);
  if (e.shiftKey) {
    ctx.fillStyle = '#000000';
    ctx.fillRect(x, y, 10, 10);
  } else {
    // ctx.fillStyle = '#FFFFFF';
    // ctx.fillRect(x, y, 20, 20);
  }
});

saveCanvasBtn.addEventListener('click', _ => {
  const data = canvas.toDataURL();
  const img = new Image();
  img.onload = () => {
    ctx2.drawImage(img, 0, 0);
  };
  img.src = data;
});
