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

  // Converter: swap the two currencies and re-run the conversion
  document.addEventListener("click", function (e) {
    var btn = e.target.closest("[data-swap-currencies]");
    if (!btn) return;
    var form = btn.closest("form");
    var a = form.querySelector("[name=from]"), b = form.querySelector("[name=to]");
    var t = a.value; a.value = b.value; b.value = t;
    a.dispatchEvent(new Event("change", { bubbles: true }));
  });

  // In-page links: scroll to the section without leaving "#section" in the address bar.
  var marked = null;
  function goTo(el) {
    if (marked) marked.classList.remove("is-target");
    el.classList.add("is-target");
    marked = el;
    if (!el.hasAttribute("tabindex") && !el.matches("a[href], button, input, select, textarea, summary")) el.setAttribute("tabindex", "-1");
    el.scrollIntoView({ behavior: reduce.matches ? "auto" : "smooth", block: "start" });
    el.focus({ preventScroll: true });
  }
  function hashTarget(hash) {
    if (!hash || hash.length < 2) return null;
    try { return document.getElementById(decodeURIComponent(hash.slice(1))); } catch (err) { return null; }
  }
  document.addEventListener("click", function (e) {
    if (e.defaultPrevented || e.button !== 0 || e.metaKey || e.ctrlKey || e.shiftKey || e.altKey) return;
    var a = e.target.closest("a[href*='#']");
    if (!a || a.target || a.hasAttribute("download") || a.hasAttribute("hx-get")) return;
    var url = new URL(a.href, location.href);
    if (url.origin !== location.origin || url.pathname !== location.pathname || url.search !== location.search) return;
    var el = hashTarget(url.hash);
    if (!el) return;
    e.preventDefault();
    goTo(el);
  });
  // Arriving from another page with "#section": jump there, then tidy the URL.
  if (location.hash) {
    var arrived = hashTarget(location.hash);
    history.replaceState(history.state, "", location.pathname + location.search);
    if (arrived) {
      arrived.classList.add("is-target");
      marked = arrived;
      arrived.scrollIntoView({ block: "start" });
      var landed = window.scrollY;
      window.addEventListener("load", function () { if (Math.abs(window.scrollY - landed) < 4) arrived.scrollIntoView({ block: "start" }); });
    }
  }

  // Re-run enhancements on content htmx swaps in.
  document.addEventListener("htmx:afterSettle", function (e) {
    reveal(e.target);
    collect();
    document.dispatchEvent(new CustomEvent("karibu:swap", { detail: { target: e.target } }));
  });
  window.Karibu = { reveal: reveal, reduceMotion: function () { return reduce.matches; } };
})();
