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

  // Parallax hero: layers drift up as you scroll (down would vanish into overflow).
  var layers = [];
  function collect() {
    layers = [];
    document.querySelectorAll("[data-parallax]").forEach(function (hero) {
      hero.classList.add("is-ready");
      hero.querySelectorAll("[data-depth]").forEach(function (el) {
        layers.push({ el: el, depth: parseFloat(el.getAttribute("data-depth")) || 0 });
      });
    });
  }
  var ticking = false;
  function onScroll() {
    if (ticking || reduce.matches) return;
    ticking = true;
    requestAnimationFrame(function () {
      var y = window.scrollY;
      if (y < window.innerHeight * 1.8) {
        layers.forEach(function (l) {
          l.el.style.setProperty("--parallax", (-y * l.depth * 1.35).toFixed(1) + "px");
        });
      }
      ticking = false;
    });
  }
  requestAnimationFrame(function () {
    collect();
    onScroll();
  });
  window.addEventListener("scroll", onScroll, { passive: true });
  reduce.addEventListener("change", function () {
    if (reduce.matches) {
      layers.forEach(function (l) { l.el.style.removeProperty("--parallax"); });
    } else {
      onScroll();
    }
  });

  // Park the hero's drive while it is off screen so it costs nothing.
  if ("IntersectionObserver" in window) {
    var park = new IntersectionObserver(function (entries) {
      entries.forEach(function (en) { en.target.classList.toggle("is-parked", !en.isIntersecting); });
    });
    document.querySelectorAll(".hero[data-parallax]").forEach(function (h) { park.observe(h); });
  }

  // Ease in hero copy and Money/Map cards that arrive out of band.
  document.addEventListener("htmx:oobAfterSwap", function (e) {
    var el = e.detail.target;
    if (!el) return;
    var animate = el.classList.contains("hero__copy") || el.id === "home-glance";
    if (!animate) return;
    el.querySelectorAll(".reveal-stagger").forEach(function (n) { n.classList.add("is-in"); });
    el.classList.remove("is-new");
    void el.offsetWidth;
    el.classList.add("is-new");
  });

  // Converter: swap the two currencies and re-run the conversion
  document.addEventListener("click", function (e) {
    var btn = e.target.closest("[data-swap-currencies]");
    if (!btn) return;
    var form = btn.closest("form");
    var a = form.querySelector("[name=from]"), b = form.querySelector("[name=to]");
    var t = a.value; a.value = b.value; b.value = t;
    a.dispatchEvent(new Event("change", { bubbles: true }));
  });

  // Re-run enhancements on content htmx swaps in.
  document.addEventListener("htmx:afterSettle", function (e) {
    reveal(e.target);
    collect();
    document.dispatchEvent(new CustomEvent("karibu:swap", { detail: { target: e.target } }));
  });
  window.Karibu = { reveal: reveal, reduceMotion: function () { return reduce.matches; } };
})();
