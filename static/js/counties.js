// Counties page enhancements: map tooltips, map dimming for the filtered grid, legend filters.
// Everything works without this file: counties are links and filters are a GET form.
(function () {
  "use strict";

  function reduce() { return window.Karibu && window.Karibu.reduceMotion(); }

  // Tooltip + highlight outline for each Kenya map
  function setupMap(fig) {
    if (fig.dataset.ready) return;
    fig.dataset.ready = "1";
    var frame = fig.querySelector(".kmap__frame");
    var tip = fig.querySelector(".kmap__tip");
    var hl = fig.querySelector(".kmap__hl");
    var nameEl = tip.querySelector(".kmap__tip-name");
    var metaEl = tip.querySelector(".kmap__tip-meta");
    var current = null;

    function show(a, x, y, focus) {
      var path = a.querySelector("path");
      if (current !== a) {
        current = a;
        nameEl.textContent = a.getAttribute("data-name");
        metaEl.textContent = "Seat: " + a.getAttribute("data-seat") + " · " + a.getAttribute("data-region") + " region";
        hl.setAttribute("d", path.getAttribute("d"));
      }
      hl.classList.toggle("is-focus", !!focus);
      var box = frame.getBoundingClientRect();
      if (x == null) {
        var r = path.getBoundingClientRect();
        x = r.left + r.width / 2;
        y = r.top + Math.min(r.height / 2, 20);
      }
      var px = x - box.left, py = y - box.top;
      tip.hidden = false;
      var half = tip.offsetWidth / 2 + 4;
      px = Math.max(half, Math.min(box.width - half, px));
      var below = py < tip.offsetHeight + 24;
      tip.classList.toggle("is-below", below);
      tip.style.left = px.toFixed(0) + "px";
      tip.style.top = py.toFixed(0) + "px";
    }
    function hide() {
      current = null;
      tip.hidden = true;
      hl.setAttribute("d", "");
    }

    fig.addEventListener("pointermove", function (e) {
      var a = e.target.closest("a.kmap__county");
      if (!a) { if (current && document.activeElement !== current) hide(); return; }
      show(a, e.clientX, e.clientY, false);
    });
    fig.addEventListener("pointerleave", function () {
      if (current && document.activeElement === current) show(current, null, null, true);
      else hide();
    });
    fig.addEventListener("focusin", function (e) {
      var a = e.target.closest("a.kmap__county");
      if (a) show(a, null, null, true);
    });
    fig.addEventListener("focusout", function (e) {
      if (!fig.contains(e.relatedTarget)) hide();
    });
    fig.addEventListener("keydown", function (e) {
      if (e.key === "Escape" && !tip.hidden) { hide(); }
    });
  }

  // Dim map counties that aren't in the current results, and mark the active legend region.
  function syncResults() {
    var res = document.getElementById("county-results");
    if (!res) return;
    var filtered = res.hasAttribute("data-filtered");
    var keep = {};
    (res.getAttribute("data-slugs") || "").split(" ").forEach(function (s) { if (s) keep[s] = true; });
    document.querySelectorAll(".kmap:not(.kmap--locator) a.kmap__county").forEach(function (a) {
      a.classList.toggle("is-dim", filtered && !keep[a.getAttribute("data-slug")]);
    });
    var region = res.getAttribute("data-region") || "";
    document.querySelectorAll(".legend-region[data-region-filter]").forEach(function (l) {
      if (l.getAttribute("data-region-filter") === region && region) l.setAttribute("aria-current", "true");
      else l.removeAttribute("aria-current");
    });
  }

  // Legend links drive the region pills instead of reloading the page.
  document.addEventListener("click", function (e) {
    var link = e.target.closest("[data-region-filter]");
    if (!link || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey || e.button !== 0) return;
    var input = document.querySelector('#all-counties input[name="region"][value="' + link.getAttribute("data-region-filter") + '"]');
    if (!input || !window.htmx) return;
    e.preventDefault();
    input.checked = true;
    input.dispatchEvent(new Event("change", { bubbles: true }));
    var grid = document.getElementById("all-counties");
    if (grid && !window.matchMedia("(min-width: 56rem)").matches) grid.scrollIntoView({ behavior: reduce() ? "auto" : "smooth", block: "start" });
  });

  function init() {
    document.querySelectorAll("[data-kenya-map]").forEach(setupMap);
    syncResults();
  }
  if (document.readyState === "loading") document.addEventListener("DOMContentLoaded", init);
  else init();
  document.addEventListener("karibu:swap", init);
})();
