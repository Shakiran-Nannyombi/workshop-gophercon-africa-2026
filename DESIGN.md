# Karibu Guide design notes

Direction: "signwritten safari": hand-painted sign lettering, warm sand paper, printed cards
(thick ink outline + offset shadow), kanga borders, layered parallax landscapes. No remote images.

- **Colours**: CSS custom properties at the top of `static/css/app.css` (`--paper`, `--surface`, `--ink`,
  `--green` Rift Valley, `--red` laterite, `--gold` savannah, `--night` indigo, `--lake`, `--muted`).
  Dark mode swaps them under `[data-theme="dark"]`. Gold is a fill colour only; text on it is ink.
- **Lettering**: Bricolage Grotesque (headings, 800), Figtree (body), JetBrains Mono (numbers, eyebrows).
  Self-hosted in `static/fonts/`.
- **Styling**: plain CSS, no build step (Tailwind was dropped). Semantic class names.
- **Motion**: `.reveal` / `.reveal-stagger` on scroll, `[data-parallax] [data-depth]` hero layers.
  Everything is off under `prefers-reduced-motion`.
- **Avoid**: emoji as icons, remote images, stereotyping safety copy, unmarked approximate figures.
