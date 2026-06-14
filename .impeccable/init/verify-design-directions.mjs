import { chromium } from '../../web/playtest/node_modules/playwright/index.mjs';

const fileUrl = new URL('./stonks-design-directions.html', import.meta.url).href;

function parseColor(value) {
  if (value.startsWith('rgb(') || value.startsWith('rgba(')) {
    const [r, g, b, a = 1] = value.match(/[\d.]+/g).map(Number);
    return { r, g, b, a };
  }

  if (value.startsWith('oklch(')) {
    const [, rawL, rawC, rawH, rawA] = value.match(/oklch\(([\d.]+)\s+([\d.]+)\s+([\d.]+)(?:\s*\/\s*([\d.]+))?\)/) ?? [];
    const rgb = oklchToRgb(Number(rawL), Number(rawC), Number(rawH));
    return { ...rgb, a: rawA === undefined ? 1 : Number(rawA) };
  }

  throw new Error(`Unsupported color format: ${value}`);
}

function oklchToRgb(l, c, h) {
  const radians = (h * Math.PI) / 180;
  const a = c * Math.cos(radians);
  const b = c * Math.sin(radians);

  const lPrime = l + 0.3963377774 * a + 0.2158037573 * b;
  const mPrime = l - 0.1055613458 * a - 0.0638541728 * b;
  const sPrime = l - 0.0894841775 * a - 1.291485548 * b;

  const l3 = lPrime ** 3;
  const m3 = mPrime ** 3;
  const s3 = sPrime ** 3;

  return {
    r: encodeRgb(4.0767416621 * l3 - 3.3077115913 * m3 + 0.2309699292 * s3),
    g: encodeRgb(-1.2684380046 * l3 + 2.6097574011 * m3 - 0.3413193965 * s3),
    b: encodeRgb(-0.0041960863 * l3 - 0.7034186147 * m3 + 1.707614701 * s3)
  };
}

function encodeRgb(value) {
  const bounded = Math.min(1, Math.max(0, value));
  const encoded = bounded <= 0.0031308 ? 12.92 * bounded : 1.055 * bounded ** (1 / 2.4) - 0.055;
  return Math.round(encoded * 255);
}

function relativeLuminance({ r, g, b }) {
  const [sr, sg, sb] = [r, g, b].map((channel) => {
    const normalized = channel / 255;
    return normalized <= 0.03928 ? normalized / 12.92 : ((normalized + 0.055) / 1.055) ** 2.4;
  });
  return 0.2126 * sr + 0.7152 * sg + 0.0722 * sb;
}

function contrast(fg, bg) {
  const lighter = Math.max(relativeLuminance(fg), relativeLuminance(bg));
  const darker = Math.min(relativeLuminance(fg), relativeLuminance(bg));
  return (lighter + 0.05) / (darker + 0.05);
}

function composite(top, bottom) {
  const topAlpha = top.a ?? 1;
  const bottomAlpha = bottom.a ?? 1;
  const alpha = topAlpha + bottomAlpha * (1 - topAlpha);
  if (alpha === 0) return { r: 255, g: 255, b: 255, a: 1 };

  return {
    r: Math.round((top.r * topAlpha + bottom.r * bottomAlpha * (1 - topAlpha)) / alpha),
    g: Math.round((top.g * topAlpha + bottom.g * bottomAlpha * (1 - topAlpha)) / alpha),
    b: Math.round((top.b * topAlpha + bottom.b * bottomAlpha * (1 - topAlpha)) / alpha),
    a: alpha
  };
}

function effectiveBackground(values) {
  return values
    .map(parseColor)
    .reverse()
    .reduce((bottom, top) => composite(top, bottom), { r: 255, g: 255, b: 255, a: 1 });
}

async function main() {
  const browser = await chromium.launch();
  const page = await browser.newPage({ viewport: { width: 1440, height: 1200 }, deviceScaleFactor: 1 });
  const consoleMessages = [];

  page.on('console', (message) => {
    if (message.type() === 'error') consoleMessages.push(message.text());
  });
  page.on('pageerror', (error) => consoleMessages.push(error.message));

  await page.goto(fileUrl);
  await page.screenshot({ path: '/tmp/stonks-design-directions-desktop.png', fullPage: true });

  const desktopOverflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);
  const contrastChecks = await page.evaluate(() => {
    function backgroundStackFor(element) {
      const backgrounds = [];
      let current = element;
      while (current) {
        const color = getComputedStyle(current).backgroundColor;
        if (color !== 'rgba(0, 0, 0, 0)' && color !== 'transparent') backgrounds.push(color);
        current = current.parentElement;
      }
      backgrounds.push(getComputedStyle(document.documentElement).backgroundColor);
      return backgrounds;
    }

    return [
      ['page body', 'body'],
      ['intro prose', '.intro p'],
      ['a heading', '.direction-a h2'],
      ['a muted copy', '.direction-a .panel-sub'],
      ['b dark nav', '.direction-b .nav-item'],
      ['b dark heading', '.direction-b .panel-title'],
      ['c muted copy', '.direction-c .event span'],
      ['d nav copy', '.direction-d .nav-item'],
      ['d panel title', '.direction-d .panel-title']
    ].map(([label, selector]) => {
      const element = document.querySelector(selector);
      const style = getComputedStyle(element);
      return {
        label,
        selector,
        color: style.color,
        backgrounds: backgroundStackFor(element)
      };
    });
  });

  await page.setViewportSize({ width: 390, height: 1200 });
  await page.goto(fileUrl);
  await page.screenshot({ path: '/tmp/stonks-design-directions-mobile.png', fullPage: true });
  const mobileOverflow = await page.evaluate(() => document.documentElement.scrollWidth - window.innerWidth);

  await browser.close();

  const ratios = contrastChecks.map((check) => {
    const background = effectiveBackground(check.backgrounds);
    const ratio = contrast(parseColor(check.color), background);
    return { ...check, ratio: Number(ratio.toFixed(2)) };
  });

  console.log(JSON.stringify({
    fileUrl,
    screenshots: [
      '/tmp/stonks-design-directions-desktop.png',
      '/tmp/stonks-design-directions-mobile.png'
    ],
    consoleErrors: consoleMessages,
    overflowPx: {
      desktop: desktopOverflow,
      mobile: mobileOverflow
    },
    contrast: ratios
  }, null, 2));
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
