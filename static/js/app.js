// Karibu Guide: small progressive enhancements. The site works without this file.
(function () {
  "use strict";
  var html = document.documentElement;
  var reduce = window.matchMedia("(prefers-reduced-motion: reduce)");

  // Dark mode toggle
  function syncToggle() {
    var dark = html.getAttribute("data-theme") === "dark";
    document.querySelectorAll("[data-theme-toggle]").forEach(function (b) {
      b.setAttribute("aria-pressed", String(dark));
    });
  }
  document.addEventListener("click", function (e) {
    var btn = e.target.closest("[data-theme-toggle]");
    if (!btn) return;
    var next = html.getAttribute("data-theme") === "dark" ? "light" : "dark";
    html.setAttribute("data-theme", next);
    try { localStorage.setItem("karibu-theme", next); } catch (err) {}
    syncToggle();
  });
  syncToggle();

  // Close the mobile menu on Escape or outside click
  document.addEventListener("keydown", function (e) {
    if (e.key !== "Escape") return;
    document.querySelectorAll("details[data-menu][open]").forEach(function (d) {
      d.removeAttribute("open");
      d.querySelector("summary").focus();
    });
  });
  document.addEventListener("click", function (e) {
    document.querySelectorAll("details[data-menu][open]").forEach(function (d) {
      if (!d.contains(e.target)) d.removeAttribute("open");
    });
  });

  // Scroll reveals
  var io = null;
  function reveal(root) {
    var els = (root || document).querySelectorAll(".reveal:not(.is-in), .reveal-stagger:not(.is-in)");
    if (reduce.matches || !("IntersectionObserver" in window)) {
      els.forEach(function (el) { el.classList.add("is-in"); });
      return;
    }
    if (!io) {
      io = new IntersectionObserver(function (entries) {
        entries.forEach(function (en) {
          if (en.isIntersecting) { en.target.classList.add("is-in"); io.unobserve(en.target); }
        });
      }, { rootMargin: "0px 0px -8% 0px", threshold: 0.08 });
    }
    els.forEach(function (el) { io.observe(el); });
  }
  reveal();

  // Parallax hero: each [data-depth] layer moves at its own speed.
  var layers = [];
  function collect() {
    layers = [];
    document.querySelectorAll("[data-parallax] [data-depth]").forEach(function (el) {
      layers.push({ el: el, depth: parseFloat(el.getAttribute("data-depth")) || 0 });
    });
  }
  var ticking = false;
  function onScroll() {
    if (ticking || reduce.matches) return;
    ticking = true;
    requestAnimationFrame(function () {
      var y = window.scrollY;
      if (y < window.innerHeight * 1.5) {
        layers.forEach(function (l) { l.el.style.transform = "translate3d(0," + (y * l.depth).toFixed(1) + "px,0)"; });
      }
      ticking = false;
    });
  }
  collect();
  window.addEventListener("scroll", onScroll, { passive: true });
  reduce.addEventListener("change", function () {
    if (reduce.matches) layers.forEach(function (l) { l.el.style.transform = ""; });
  });

  // Re-run enhancements on content htmx swaps in.
  document.addEventListener("htmx:afterSettle", function (e) {
    reveal(e.target);
    collect();
    document.dispatchEvent(new CustomEvent("karibu:swap", { detail: { target: e.target } }));
  });
  window.Karibu = { reveal: reveal, reduceMotion: function () { return reduce.matches; } };
})();
