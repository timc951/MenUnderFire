#!/usr/bin/env node
/**
 * Image optimization script for Men Under Fire
 *
 * Converts PNGs to WebP, resizes oversized images, and creates
 * mobile variants. Output goes to public/optimized/ — upload
 * the contents to your CDN (content.menunderfire.com).
 *
 * Usage:
 *   cd frontend
 *   npx --yes --package=sharp -- node scripts/optimize-images.cjs
 */

const sharp = require('sharp');
const { mkdir, stat } = require('node:fs/promises');
const { join, basename, extname } = require('node:path');

const PUBLIC_DIR = join(__dirname, '..', 'public');
const OUT_DIR = join(PUBLIC_DIR, 'optimized');

// Per-image optimization rules
const IMAGE_RULES = {
  'poster.png': [
    // Full-size WebP for desktop (background image)
    { suffix: '', format: 'webp', quality: 80 },
    // Smaller mobile version (600px wide)
    { suffix: '-mobile', format: 'webp', quality: 75, resize: { width: 600 } },
    // Keep a PNG fallback but compressed
    { suffix: '', format: 'png', pngQuality: true },
  ],
  'favicon.png': [
    // Favicon only needs to be small — 64x64 is plenty
    { suffix: '', format: 'png', resize: { width: 64, height: 64 }, pngQuality: true },
    { suffix: '', format: 'webp', resize: { width: 64, height: 64 }, quality: 80 },
  ],
  'logo.png': [
    // Logo displayed at 32-48px — 96px source is fine for 2x/3x
    { suffix: '', format: 'webp', resize: { width: 96, height: 96 }, quality: 85 },
    { suffix: '', format: 'png', resize: { width: 96, height: 96 }, pngQuality: true },
  ],
  'MenUnderFireIcon.png': [
    // Icon displayed at 36px — 108px for 3x retina
    { suffix: '', format: 'webp', resize: { width: 108, height: 108 }, quality: 85 },
    { suffix: '', format: 'png', resize: { width: 108, height: 108 }, pngQuality: true },
  ],
  'MenUnderFireStone2.png': [
    // Wordmark — keep full width, just convert format
    { suffix: '', format: 'webp', quality: 85 },
    { suffix: '', format: 'png', pngQuality: true },
  ],
  'MenUnderFireStone.png': [
    { suffix: '', format: 'webp', quality: 85 },
    { suffix: '', format: 'png', pngQuality: true },
  ],
  'MenUnderFireLava.png': [
    { suffix: '', format: 'webp', quality: 85 },
    { suffix: '', format: 'png', pngQuality: true },
  ],
};

function formatBytes(bytes) {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(2)} MB`;
}

async function processImage(filename, rules) {
  const inputPath = join(PUBLIC_DIR, filename);
  const name = basename(filename, extname(filename));
  const originalStat = await stat(inputPath);
  const originalSize = originalStat.size;

  console.log(`\n--- ${filename} (${formatBytes(originalSize)}) ---`);

  for (const rule of rules) {
    const outName = `${name}${rule.suffix}.${rule.format}`;
    const outPath = join(OUT_DIR, outName);

    let pipeline = sharp(inputPath);

    // Resize if specified
    if (rule.resize) {
      pipeline = pipeline.resize(rule.resize.width, rule.resize.height, {
        fit: 'inside',
        withoutEnlargement: true,
      });
    }

    // Format conversion
    if (rule.format === 'webp') {
      pipeline = pipeline.webp({ quality: rule.quality || 80 });
    } else if (rule.format === 'png' && rule.pngQuality) {
      pipeline = pipeline.png({
        compressionLevel: 9,
        palette: true,
        quality: 80,
      });
    }

    await pipeline.toFile(outPath);

    const outStat = await stat(outPath);
    const savings = ((1 - outStat.size / originalSize) * 100).toFixed(1);
    console.log(
      `  -> ${outName}: ${formatBytes(outStat.size)} (${savings}% smaller)`
    );
  }
}

async function main() {
  await mkdir(OUT_DIR, { recursive: true });

  console.log('=== Men Under Fire Image Optimization ===');
  console.log(`Source: ${PUBLIC_DIR}`);
  console.log(`Output: ${OUT_DIR}\n`);

  let totalOriginal = 0;
  let totalOptimized = 0;

  for (const [filename, rules] of Object.entries(IMAGE_RULES)) {
    try {
      const inputPath = join(PUBLIC_DIR, filename);
      const originalStat = await stat(inputPath);
      totalOriginal += originalStat.size;

      await processImage(filename, rules);

      // Count the primary WebP output for total savings
      const name = basename(filename, extname(filename));
      const primaryOut = join(OUT_DIR, `${name}.webp`);
      try {
        const outStat = await stat(primaryOut);
        totalOptimized += outStat.size;
      } catch {
        const pngOut = join(OUT_DIR, `${name}.png`);
        const outStat = await stat(pngOut);
        totalOptimized += outStat.size;
      }
    } catch (err) {
      console.error(`  SKIP ${filename}: ${err.message}`);
    }
  }

  console.log('\n=== Summary ===');
  console.log(`Original total:  ${formatBytes(totalOriginal)}`);
  console.log(`Optimized total: ${formatBytes(totalOptimized)} (WebP primary)`);
  console.log(`Savings:         ${((1 - totalOptimized / totalOriginal) * 100).toFixed(1)}%`);
  console.log(`\nOutput in: ${OUT_DIR}`);
  console.log('Upload these files to your CDN (content.menunderfire.com).');
  console.log('\nNext steps:');
  console.log('  1. Upload the .webp files to your CDN');
  console.log('  2. Update image references in code to use .webp (or use Cloudflare Polish for auto-conversion)');
  console.log('  3. Replace favicon.png on CDN with the optimized 64x64 version');
  console.log('  4. Use poster-mobile.webp for mobile via CSS media query or <picture> element');
}

main().catch(console.error);
