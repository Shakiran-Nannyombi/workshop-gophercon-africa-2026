// Runs before paint: apply the saved theme, or follow the system setting.
(function () {
  var saved = null;
  try { saved = localStorage.getItem("karibu-theme"); } catch (e) {}
  var dark = saved ? saved === "dark" : window.matchMedia("(prefers-color-scheme: dark)").matches;
  var html = document.documentElement;
  html.setAttribute("data-theme", dark ? "dark" : "light");
  html.classList.add("js");
})();
