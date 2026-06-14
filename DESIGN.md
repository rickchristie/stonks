---
name: Stonks
description: A personal-first stock research, portfolio, and trading automation command center.
colors:
  canvas: "oklch(1.000 0.000 0)"
  canvas-soft: "oklch(0.968 0.008 210)"
  ink: "oklch(0.165 0.018 230)"
  muted: "oklch(0.430 0.030 225)"
  line: "oklch(0.875 0.014 220)"
  steel-blue: "oklch(0.510 0.105 209)"
  signal-cyan: "oklch(0.690 0.145 205)"
  gain-pop: "oklch(0.540 0.170 151)"
  loss-pop: "oklch(0.590 0.215 28)"
  brass: "oklch(0.690 0.120 63)"
  signal-black: "oklch(0.085 0.000 0)"
  signal-panel: "oklch(0.108 0.004 235)"
  signal-surface: "oklch(0.130 0.006 235)"
typography:
  display:
    fontFamily: "\"Inter Tight\", \"Arial Black\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: "clamp(3rem, 8vw, 5.5rem)"
    fontWeight: 850
    lineHeight: 0.94
    letterSpacing: "-0.035em"
  headline:
    fontFamily: "\"Inter Tight\", \"Arial Black\", \"Helvetica Neue\", Arial, sans-serif"
    fontSize: "3rem"
    fontWeight: 820
    lineHeight: 0.98
    letterSpacing: "-0.03em"
  title:
    fontFamily: "Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, \"Segoe UI\", sans-serif"
    fontSize: "0.875rem"
    fontWeight: 760
    lineHeight: 1.2
  body:
    fontFamily: "Inter, ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, \"Segoe UI\", sans-serif"
    fontSize: "0.9375rem"
    fontWeight: 450
    lineHeight: 1.55
  mono:
    fontFamily: "Consolas, Monaco, \"SFMono-Regular\", \"Liberation Mono\", monospace"
    fontSize: "0.75rem"
    fontWeight: 650
    lineHeight: 1.35
rounded:
  sm: "8px"
  md: "10px"
  panel: "12px"
  shell: "16px"
  pill: "999px"
spacing:
  xs: "4px"
  sm: "8px"
  md: "12px"
  lg: "16px"
  xl: "24px"
  xxl: "32px"
components:
  button-primary:
    backgroundColor: "{colors.steel-blue}"
    textColor: "{colors.canvas}"
    rounded: "{rounded.pill}"
    padding: "9px 12px"
  button-signal:
    backgroundColor: "{colors.signal-cyan}"
    textColor: "{colors.signal-black}"
    rounded: "{rounded.pill}"
    padding: "9px 12px"
  button-quiet:
    backgroundColor: "transparent"
    textColor: "{colors.muted}"
    rounded: "{rounded.pill}"
    padding: "8px 11px"
  input-command:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink}"
    rounded: "{rounded.md}"
    padding: "9px 10px"
  panel:
    backgroundColor: "{colors.canvas}"
    textColor: "{colors.ink}"
    rounded: "{rounded.panel}"
    padding: "14px"
  ticker-tile-gain:
    backgroundColor: "{colors.gain-pop}"
    textColor: "{colors.canvas}"
    rounded: "{rounded.md}"
    padding: "12px"
  ticker-tile-loss:
    backgroundColor: "{colors.loss-pop}"
    textColor: "{colors.canvas}"
    rounded: "{rounded.md}"
    padding: "12px"
---

# Design System: Stonks

## 1. Overview

**Creative North Star: "The Signal Instrument"**

Stonks is an A+B hybrid from the init directions: the bright precision of Instrument Panel plus the focused energy of Signal Room. The default workspace is light, exact, readable, and calm enough for deep research. Signal-heavy areas such as alerts, analyst scanning, and future automation can shift into a dark surface where glowing cyan, gain green, and loss red feel alive.

The system should not look like a brokerage clone or a generic admin dashboard. It should be beautiful enough to keep open all day, with lively navigation, animated percentage tickers, compact watchlists, and semantic color that actually pops. The fun is in the data behavior and navigation rhythm, not decorative effects.

**Key Characteristics:**
- Heavy grotesque sans display type with tight but readable spacing.
- Consolas/Monaco-first monospace for tickers, percentages, shortcuts, and rule syntax.
- Light analytical canvas by default, dark signal mode for alerting and automation.
- Vivid red/green/cyan semantic tiles that read instantly.
- Dense panels, tables, charts, and timelines with stable dimensions.

## 2. Colors

The palette is blue-steel architecture with bright data colors: calm surfaces for reading, dark signal surfaces for focus, and semantic colors that pop hard enough for all-day monitoring.

### Primary
- **Blue-Steel Instrument** (`steel-blue`): primary navigation selection, light-mode primary action, chart line, and core brand mark. It should feel precise, not corporate.

### Secondary
- **Signal Cyan** (`signal-cyan`): dark-mode action, command focus, active automation, and data tiles that are informational rather than positive or negative.
- **Patina Brass** (`brass`): rare research accent for thesis bullets, document highlights, and "needs review" moments. Use it sparingly.

### Tertiary
- **Gain Pop** (`gain-pop`): positive performance, up moves, successful analyst results, and safe active states.
- **Loss Pop** (`loss-pop`): negative performance, drawdown, failed analyst results, blocked rules, and urgent risk.

### Neutral
- **Canvas White** (`canvas`): primary reading and research surface.
- **Soft Instrument Surface** (`canvas-soft`): navigation rail, panel wells, and subtle grouping.
- **Ink** (`ink`): body text and core data labels on light surfaces.
- **Muted Ink** (`muted`): secondary text, metadata, and inactive nav.
- **Signal Black / Signal Panel / Signal Surface** (`signal-black`, `signal-panel`, `signal-surface`): dark signal-room stack for alerts, automation, and monitoring.

### Named Rules

**The Pop Means Data Rule.** Bright green, red, and cyan are for market movement, risk, analyst state, and automation state. Never use them as random decoration.

**The Two Rooms Rule.** Use the light room for reading and research. Use the dark room for scanning, alerting, and automation supervision. Do not smear dark-mode styling across every screen just because trading tools often do.

## 3. Typography

**Display Font:** Heavy grotesque sans, preferably `Inter Tight`, with `Arial Black`, `Helvetica Neue`, and Arial fallbacks.
**Body Font:** Inter/system sans.
**Label/Mono Font:** Consolas, Monaco, `SFMono-Regular`, `Liberation Mono`, monospace.

**Character:** Big headings should have the same confident heavy shape as the chosen reference image: massive, sans-serif, tight, and clean. Product labels stay smaller and quieter so tables, watchlists, and rule builders remain usable.

### Hierarchy

- **Display** (850, `clamp(3rem, 8vw, 5.5rem)`, 0.94): hero-scale setup screens, empty states, and major product moments only.
- **Headline** (820, `3rem`, 0.98): page headers and major module names.
- **Title** (760, `0.875rem`, 1.2): panel headers, table group labels, and nav section names.
- **Body** (450, `0.9375rem`, 1.55): readable prose, research notes, thesis summaries, and explanations. Cap prose at 65-75ch.
- **Mono** (650, `0.75rem`, 1.35): ticker symbols, percentages, dates, shortcuts, rules, analyst scores, and compact metadata.

### Named Rules

**The Big Type Is Rare Rule.** The heavy display voice is powerful because it is rare. Do not use hero-scale type inside dense dashboards, tables, sidebars, or compact panels.

**The Ticker Type Rule.** Market symbols, percentages, rule syntax, and keyboard hints use Consolas/Monaco-style monospace. This is part of the product identity.

## 4. Elevation

Stonks uses tonal layering and precise borders first. Light-mode panels can use a small structural shadow only for the main shell or hover lift; dark-mode surfaces should feel inset and illuminated through color, not through blurry glass effects.

### Shadow Vocabulary

- **Shell Lift** (`0 8px 8px color-mix(in oklch, var(--steel-blue) 8%, transparent)`): one main app screenshot, main modal, or important floating shell.
- **Interactive Lift** (`0 4px 8px color-mix(in oklch, var(--ink) 8%, transparent)`): temporary hover or drag state only.

### Named Rules

**The Flat Until State Rule.** Surfaces are flat at rest. Elevation appears only for structural shells or interactive state, never as decorative card haze.

## 5. Components

### Buttons

- **Shape:** confident pill actions (`999px`) with compact height for dense toolbars.
- **Primary:** steel-blue fill on light surfaces; signal-cyan fill on dark signal surfaces.
- **Hover / Focus:** brighten the fill slightly, show a clear focus ring, and move no more than 1px. Do not animate layout.
- **Secondary / Quiet:** transparent pill with line color and muted text for non-primary commands.

### Chips

- **Style:** small mono labels with full rounded shape, soft fill, and readable text.
- **State:** selected chips use steel-blue or signal-cyan with white or near-black text based on fill luminance. Inactive chips remain quiet.

### Cards / Containers

- **Corner Style:** restrained product radius (`12px` panels, `16px` outer shell).
- **Background:** canvas white or signal panel, never cream/beige.
- **Shadow Strategy:** follow Elevation. Prefer border and tonal layering.
- **Border:** 1px full border using `line`; never colored side stripes.
- **Internal Padding:** dense panels use 14-16px; research reading panes can use 24px.

### Inputs / Fields

- **Style:** command-search fields are rounded rectangles (`10px`) with subtle border, compact padding, and a mono prefix or shortcut hint.
- **Focus:** shift border to steel-blue or signal-cyan and add a visible focus ring.
- **Error / Disabled:** error uses loss-pop plus text; disabled lowers opacity but remains readable.

### Navigation

- **Style:** navigation should be aesthetically pleasing and fun: active state pills, ticker-like status snippets, small percentage movement, and compact mono shortcuts.
- **Behavior:** active nav can animate a small percentage/ticker readout, but reduced motion must freeze to a static value. Navigation should feel alive without becoming a slot machine.
- **Mobile:** collapse to horizontal or stacked controls with stable sizes. Text must not overflow chips or nav items.

### Ticker Tiles

Ticker tiles are the signature component. They use bold semantic fills, mono labels, and stable dimensions so movement never shifts layout. Gain uses `gain-pop`, loss uses `loss-pop`, informational uses `signal-cyan`.

## 6. Do's and Don'ts

### Do:

- **Do** build from the A+B direction: bright Instrument Panel for research, dark Signal Room for alerting and automation.
- **Do** use the heavy sans display style from the reference image for major moments.
- **Do** use Consolas, Monaco, `SFMono-Regular`, and `Liberation Mono` for ticker and rule language.
- **Do** make semantic gain/loss colors pop, especially in ticker tiles, watchlists, charts, and analyst result states.
- **Do** include non-color cues with red and green states: plus/minus signs, labels, icons, or position.
- **Do** make navigation feel alive with percent tickers, compact status, and small state motion.

### Don't:

- **Don't** make Stonks look boring, generic, or like a copied admin template.
- **Don't** use beige/serif template styling, dull gray dashboards, washed-out semantic colors, generic fintech navy-and-gold, purple SaaS gradients, crypto-casino neon, glassmorphism as decoration, endless identical card grids, or terminal-only dark mode.
- **Don't** use the rejected C/D color direction as the visual baseline; Atlas is a feature, not the muted look from that variant.
- **Don't** use bright semantic colors as decoration when no data state is present.
- **Don't** hide automation risk. Dry-run mode, blocked actions, audit logs, and kill switches must be visible in automation surfaces.
