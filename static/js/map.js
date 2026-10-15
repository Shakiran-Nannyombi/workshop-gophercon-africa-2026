// Zoom and pan for [data-zoom-map] SVGs: buttons, drag, wheel (with Ctrl/⌘) and keyboard.
(function () {
  "use strict";
  function init(root) {
    if (root.dataset.ready) return;
    root.dataset.ready = "1";
    var svg = root.querySelector("svg");
    var port = root.querySelector(".zoom-map__viewport");
    var home = svg.getAttribute("data-home").split(" ").map(Number);
    var full = svg.getAttribute("data-full").split(" ").map(Number);
    var vb = home.slice();
    var reduce = window.matchMedia("(prefers-reduced-motion: reduce)").matches;

    function clamp() {
      var minW = full[2] / 8, maxW = full[2];
      vb[2] = Math.min(maxW, Math.max(minW, vb[2]));
      vb[3] = vb[2] * (home[3] / home[2]);
      vb[0] = Math.min(full[0] + full[2] - vb[2], Math.max(full[0], vb[0]));
      vb[1] = Math.min(full[1] + full[3] - vb[3], Math.max(full[1], vb[1]));
    }
    function apply() {
      clamp();
      svg.setAttribute("viewBox", vb.map(function (v) { return v.toFixed(1); }).join(" "));
      var z = home[2] / vb[2];
      svg.classList.toggle("is-zoomed", z > 1.4);
      svg.classList.toggle("is-zoomed-2", z > 2.2);
    }
    function zoom(f, cx, cy) {
      cx = cx == null ? vb[0] + vb[2] / 2 : cx;
      cy = cy == null ? vb[1] + vb[3] / 2 : cy;
      var w = vb[2] / f, h = vb[3] / f;
      vb = [cx - (cx - vb[0]) / f, cy - (cy - vb[1]) / f, w, h];
      apply();
    }
    function animateTo(target) {
      if (reduce) { vb = target.slice(); apply(); return; }
      var from = vb.slice(), t0 = performance.now();
      (function step(t) {
        var k = Math.min(1, (t - t0) / 260), e = 1 - Math.pow(1 - k, 3);
        vb = from.map(function (v, i) { return v + (target[i] - v) * e; });
        apply();
        if (k < 1) requestAnimationFrame(step);
      })(t0);
    }
    root.addEventListener("click", function (e) {
      var b = e.target.closest("[data-zoom]");
      if (!b) return;
      var c = [vb[0] + vb[2] / 2, vb[1] + vb[3] / 2];
      if (b.dataset.zoom === "reset") return animateTo(home);
      var f = b.dataset.zoom === "in" ? 1.5 : 1 / 1.5;
      var w = vb[2] / f;
      animateTo([c[0] - w / 2, c[1] - (w * home[3] / home[2]) / 2, w, w * home[3] / home[2]]);
    });
    function toSvg(clientX, clientY) {
      var r = svg.getBoundingClientRect();
      return [vb[0] + (clientX - r.left) / r.width * vb[2], vb[1] + (clientY - r.top) / r.height * vb[3]];
    }
    port.addEventListener("wheel", function (e) {
      if (!e.ctrlKey && !e.metaKey) return;
      e.preventDefault();
      var p = toSvg(e.clientX, e.clientY);
      zoom(e.deltaY < 0 ? 1.15 : 1 / 1.15, p[0], p[1]);
    }, { passive: false });
    var drag = null, moved = false;
    port.addEventListener("pointerdown", function (e) {
      drag = { x: e.clientX, y: e.clientY, vb: vb.slice() };
      moved = false;
    });
    window.addEventListener("pointermove", function (e) {
      if (!drag) return;
      var r = svg.getBoundingClientRect();
      var dx = (e.clientX - drag.x) / r.width * vb[2], dy = (e.clientY - drag.y) / r.height * vb[3];
      if (Math.abs(e.clientX - drag.x) + Math.abs(e.clientY - drag.y) > 4) moved = true;
      if (!moved) return;
      vb = [drag.vb[0] - dx, drag.vb[1] - dy, vb[2], vb[3]];
      apply();
    });
    window.addEventListener("pointerup", function () { drag = null; });
    port.addEventListener("click", function (e) { if (moved) { e.preventDefault(); e.stopPropagation(); moved = false; } }, true);
    port.addEventListener("keydown", function (e) {
      var step = vb[2] * 0.1;
      var k = { ArrowLeft: [-step, 0], ArrowRight: [step, 0], ArrowUp: [0, -step], ArrowDown: [0, step] }[e.key];
      if (k) { e.preventDefault(); vb[0] += k[0]; vb[1] += k[1]; apply(); }
      else if (e.key === "+" || e.key === "=") { e.preventDefault(); zoom(1.3); }
      else if (e.key === "-") { e.preventDefault(); zoom(1 / 1.3); }
      else if (e.key === "0") { e.preventDefault(); animateTo(home); }
    });
    apply();
  }
  function all() { document.querySelectorAll("[data-zoom-map]").forEach(init); }
  all();
  document.addEventListener("karibu:swap", all);
})();
