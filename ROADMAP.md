# Roadmap

> Updated: 2026-02-05 | Format: **Now / Next / Later**

## Vision

The fastest Unicode width calculation library in the Go ecosystem — correct emoji handling, zero allocations, full Unicode compliance. A drop-in replacement for go-runewidth.

---

## Now (v0.2.0 — Current Release)

- [x] **4-tier O(1) lookup**: ASCII → CJK → Emoji → 3-stage table (3.8KB)
- [x] **SWAR optimization**: ASCII detection and width counting at 8 bytes/iter
- [x] **ZWJ emoji sequences**: 👨‍👩‍👧‍👦 = width 2 (forward-scan state machine)
- [x] **Emoji modifiers**: Skin tones 👍🏽 = width 2
- [x] **Variation selectors**: U+FE0E (text) / U+FE0F (emoji)
- [x] **Regional indicator pairs**: Flag emoji 🇺🇸 = width 2
- [x] **Unicode 16.0** compliance
- [x] **96.4%** test coverage, zero lint issues

---

## Next (v0.3.0)

- [ ] **Profile-Guided Optimization (PGO)** — expected 10-20% speedup on hot paths
- [ ] **Benchmark CI** — automated regression detection on every PR
- [ ] **Unicode 17.0 preparation** — generator pipeline ready for next release
- [ ] **Keycap sequences** — `#️⃣`, `*️⃣`, `0️⃣-9️⃣`
- [ ] **Migration guide** — step-by-step from go-runewidth
- [ ] **API review** — gather feedback from early adopters

---

## Later (v1.0.0+)

### API Freeze
- [ ] Stable API guarantee (no breaking changes until v2.0)
- [ ] Semantic versioning commitment
- [ ] Validated by multiple production projects

### Explicit SIMD
- [ ] **Go assembly** (Plan 9 `.s` files) — SSE2/AVX2/NEON, 16-32 bytes/iter
- [ ] **`archsimd`** (Go 1.26+) — portable SIMD intrinsics ([golang/go#67520](https://github.com/golang/go/issues/67520))
- [ ] **AVX-512** — server-side bulk processing, 64 bytes/iter
- [ ] **ARM NEON** — Apple Silicon / AWS Graviton

### Grapheme Clusters (Conditional)
- [ ] Full [UAX #29](https://unicode.org/reports/tr29/) support (opt-in, not default)
- [ ] Reactivation: when users report incorrect widths for specific scripts

### Ecosystem
- [ ] [Phoenix TUI](https://github.com/phoenix-tui/phoenix) integration
- [ ] **unigrapheme** — companion grapheme segmentation library

---

## Non-Goals

- Automatic locale detection (use Options API)
- Font-specific width variations
- Backward compatibility below Go 1.25
- Full ICU replacement

---

## Contributing

We welcome contributions in these areas:

1. **Bug reports** — width calculation issues, emoji mismatches
2. **Performance testing** — benchmarks on different hardware
3. **Real-world usage** — integrate in your app, report API friction
4. **Unicode edge cases** — sequences we handle incorrectly

See [GitHub Issues](https://github.com/unilibs/uniwidth/issues) for current work items.
