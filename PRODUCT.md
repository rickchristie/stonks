# Product

## Register

product

## Users

Stonks is personal-first: Rick is the primary user and the product should optimize for a serious self-directed investor who watches it throughout the day. The product may later become something other serious investors can buy, so workflows should avoid private one-off assumptions when a reusable product model is just as clear.

Users are researching stocks, tracking portfolios across brokers and asset types, defining short-term trade setups, and supervising automated or semi-automated trading logic. They need fast scanning, durable research memory, visible risk, and clear explanations for every proposed or executed action.

## Product Purpose

Stonks helps investors consolidate portfolio truth, research long-term investments, define trade setups, and eventually automate disciplined trades with explicit entry, exit, risk, and audit behavior. Success means Rick can trust the app as a daily command center: useful at market open, during active monitoring, after close, and during deeper research sessions.

## Feature Model

### Portfolio Viewer

Consolidates stocks, purchase dates, cost basis, performance, alerts, and broker or asset-type fragmentation into one readable portfolio view. The portfolio surface should make current exposure, gains/losses, concentration, cash, and risk limits obvious without forcing manual reconciliation.

### Autotrader Engine

Defines rules, setups, and bots that can make short-term trading bets, calculate risk, react to events, and always carry an entry and exit strategy. Trading automation must start with dry-run behavior, keep a visible kill switch, and produce an audit trail for every proposed, skipped, blocked, or executed action.

### Trade Analyst

Defines analyst agents that constantly scan for specific setups and propose trades to the autotrader. Example analysts include MACD crossover scans, ascending-triangle scans, or watchlist-specific candle and price-pattern checks. Each analyst has tracked ideas, executions, win rate, and historical trade results so Rick can measure which analyst definitions actually work.

### Atlas

Atlas is the long-term investing research system. Each research folder is organized around a stock ticker, further categorized by industry and tagging. A ticker can contain many folders and documents, including pasted markdown from AI tools. Every ticker should have a "Thesis and Strategy" area that summarizes the current thesis, validation signals, invalidation signals, and actions to take.

Atlas can use AI agents for research jobs, such as fetching new information on next earnings, checking whether cRPO is above 21.5%, verifying operating margin above 25%, or flagging invalidating evidence. These agents help analyze and maintain research, but the interface must keep sources, assumptions, and next actions clear.

## Brand Personality

Sharp, alive, disciplined.

The product should feel aesthetically pleasing enough to live on screen all day, but never frivolous. It should combine the bright precision of an instrument panel with the energy of a signal room: fun navigation, animated percentage tickers, vivid semantic colors, and dense analytical surfaces that still feel controlled.

## Anti-references

Do not make Stonks look boring, generic, or like a copied admin template.

Avoid beige/serif template styling, dull gray dashboards, washed-out semantic colors, generic fintech navy-and-gold, purple SaaS gradients, crypto-casino neon, glassmorphism as decoration, endless identical card grids, and terminal-only dark mode. The rejected C/D init variants are too boring as a visual direction; Atlas as a product feature is valuable, but it needs the A/B energy and color discipline.

## Design Principles

1. Personal-first, product-ready. Optimize for Rick's daily use now, but model workflows cleanly enough that serious investors could use the same system later.
2. Signal should feel alive. Gains, losses, alerts, ticker movement, analyst activity, and automation state should visibly move and pop, while still respecting reduced-motion preferences.
3. Every trade idea carries its reason. Entry, exit, risk, setup definition, analyst source, and outcome tracking are first-class data.
4. Research and execution stay connected. Atlas should turn research into thesis, thesis into strategy, strategy into alerts, and alerts into disciplined action.
5. Automation is observable and interruptible. Dry-run mode, audit logs, blocked-action explanations, and kill switches are product primitives, not afterthoughts.

## Accessibility & Inclusion

Target WCAG AA or better for text contrast, with body text designed closer to 7:1 whenever practical. Reduced motion is mandatory; ticker and percentage motion should degrade to static changes or subtle crossfades. Semantic colors must pop, but red/green must never be the only signal: include signs, labels, icons, text, or position in every gain/loss and risk state.
