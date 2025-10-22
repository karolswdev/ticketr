# Visual Effects Configuration Reference

Quick reference for configuring Ticketr's TUI visual effects system.

## Environment Variables

All visual effects are controlled through environment variables set before launching `ticketr tui`.

### Theme Selection

```bash
export TICKETR_THEME=<theme_name>
```

**Available Themes**:
- `default` - Minimal, accessible, fast (no effects)
- `dark` - Blue/cyan palette with hyperspace option
- `arctic` - Cyan/white palette with snow option

**Default**: `default`

### Effect Toggles

Enable or disable specific effect categories:

```bash
export TICKETR_EFFECTS_MOTION=<true|false>    # Animations and spinners
export TICKETR_EFFECTS_SHADOWS=<true|false>   # Drop shadows on modals
export TICKETR_EFFECTS_SHIMMER=<true|false>   # Progress bar shimmer
export TICKETR_EFFECTS_AMBIENT=<true|false>   # Background atmospheric effects
```

**Defaults**:
- `TICKETR_EFFECTS_MOTION`: `true` (spinners are essential feedback)
- `TICKETR_EFFECTS_SHADOWS`: `false`
- `TICKETR_EFFECTS_SHIMMER`: `false`
- `TICKETR_EFFECTS_AMBIENT`: `false`

## Configuration Presets

### Preset 1: Minimal (Default)

Maximum compatibility, fastest performance, best for accessibility.

```bash
ticketr tui
```

No environment variables needed—this is the default.

**Effects**:
- Essential spinners only
- No shadows, shimmer, or ambient effects
- Single-line borders
- Standard colors

---

### Preset 2: Balanced (Recommended)

Good balance of visual polish and performance.

```bash
export TICKETR_THEME=dark
export TICKETR_EFFECTS_MOTION=true
export TICKETR_EFFECTS_SHADOWS=true
ticketr tui
```

**Effects**:
- Active spinners and smooth transitions
- Drop shadows on modals for depth
- Double-line focused borders
- Dark blue/cyan color scheme
- No ambient effects (performance neutral)

---

### Preset 3: Maximum Enchantment

Full visual experience—all effects enabled.

```bash
export TICKETR_THEME=dark
export TICKETR_EFFECTS_MOTION=true
export TICKETR_EFFECTS_SHADOWS=true
export TICKETR_EFFECTS_SHIMMER=true
export TICKETR_EFFECTS_AMBIENT=true
ticketr tui
```

**Effects**:
- All animations and transitions
- Drop shadows on all modals
- Progress bar shimmer sweeps
- Hyperspace background (stars streaming left to right)
- Gradient titles on focused panels

**Performance**: Approximately 3-5% CPU on modern hardware.

---

### Preset 4: Arctic Theme

Clean aesthetic with optional snow effect.

```bash
export TICKETR_THEME=arctic
export TICKETR_EFFECTS_MOTION=true
export TICKETR_EFFECTS_SHADOWS=true
export TICKETR_EFFECTS_AMBIENT=true
ticketr tui
```

**Effects**:
- Cyan/white color palette
- Gentle snowfall background
- Rounded borders
- Smooth animations
- Clean, calming atmosphere

---

### Preset 5: Performance-Constrained

For slow terminals, SSH over high-latency connections, or minimal resource usage.

```bash
export TICKETR_THEME=default
export TICKETR_EFFECTS_MOTION=false
ticketr tui
```

**Effects**:
- No animations (completely static)
- Fastest rendering
- Minimal CPU usage
- Accessibility-optimized

## .env File Configuration

For persistent configuration, add to your `.env` file:

```bash
# .env example - Visual Effects Configuration

# Theme selection
TICKETR_THEME=dark

# Effect toggles
TICKETR_EFFECTS_MOTION=true
TICKETR_EFFECTS_SHADOWS=true
TICKETR_EFFECTS_SHIMMER=false
TICKETR_EFFECTS_AMBIENT=false

# Jira configuration (as usual)
JIRA_URL=https://yourcompany.atlassian.net
JIRA_EMAIL=your.email@company.com
JIRA_API_KEY=your-api-token
JIRA_PROJECT_KEY=PROJ
```

Then load with:
```bash
source .env
ticketr tui
```

Or use `dotenv` tools:
```bash
dotenv ticketr tui
```

## Shell Alias Shortcuts

Add to your `~/.bashrc`, `~/.zshrc`, or equivalent:

```bash
# Minimal Ticketr (default)
alias tt='ticketr tui'

# Ticketr with balanced visual effects
alias tte='TICKETR_THEME=dark TICKETR_EFFECTS_SHADOWS=true ticketr tui'

# Ticketr with maximum enchantment
alias ttmax='TICKETR_THEME=dark TICKETR_EFFECTS_MOTION=true TICKETR_EFFECTS_SHADOWS=true TICKETR_EFFECTS_SHIMMER=true TICKETR_EFFECTS_AMBIENT=true ticketr tui'

# Ticketr arctic theme
alias tta='TICKETR_THEME=arctic TICKETR_EFFECTS_SHADOWS=true TICKETR_EFFECTS_AMBIENT=true ticketr tui'

# Ticketr performance mode (no effects)
alias ttperf='TICKETR_THEME=default TICKETR_EFFECTS_MOTION=false ticketr tui'
```

Then simply run:
```bash
tt        # Default
tte       # Enhanced
ttmax     # Maximum
tta       # Arctic
ttperf    # Performance
```

## Effect Descriptions

### Motion Effects (`TICKETR_EFFECTS_MOTION`)

**When Enabled**:
- Active spinners during async operations (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
- Modal fade-in transitions (░→▒→█ over 100ms)
- Focus pulse on panel borders
- Smooth state changes

**Performance**: Less than 0.5% CPU

**When Disabled**:
- Static status text (no spinners)
- Instant modal appearance
- No animations
- Fastest rendering

---

### Shadow Effects (`TICKETR_EFFECTS_SHADOWS`)

**When Enabled**:
- Drop shadows on modal dialogs
- Shadow character: ▒ (medium shade)
- Offset: 1 row down, 2 columns right
- Creates visual depth and separation

**Performance**: No CPU impact (render-time only)

**When Disabled**:
- No shadows
- Flat visual appearance
- Slightly simpler rendering

---

### Shimmer Effects (`TICKETR_EFFECTS_SHIMMER`)

**When Enabled**:
- Animated shine on progress bars
- Sweeping brightness wave across completed portion
- Subtle sparkles on success messages (✦✧⋆∗·)
- Gradient effects on focused panel titles

**Performance**: Less than 1% CPU

**When Disabled**:
- Static progress bars
- Solid colors without gradients
- No success animations

---

### Ambient Effects (`TICKETR_EFFECTS_AMBIENT`)

**When Enabled** (theme-dependent):
- **Dark theme**: Hyperspace starfield (stars streaming left to right)
- **Arctic theme**: Gentle snowfall (top to bottom)
- **Default theme**: No ambient effect (theme doesn't support it)

**Performance**: Less than 3% CPU (12-20 FPS rate-limited)

**When Disabled**:
- Solid background color
- No particle animations
- Fastest performance

## Troubleshooting

### Effects Not Visible

**Check environment variables**:
```bash
echo $TICKETR_THEME
echo $TICKETR_EFFECTS_MOTION
echo $TICKETR_EFFECTS_SHADOWS
echo $TICKETR_EFFECTS_SHIMMER
echo $TICKETR_EFFECTS_AMBIENT
```

**Verify before launching**:
```bash
env | grep TICKETR
ticketr tui
```

**Common issues**:
- Variables not exported (use `export` not just assignment)
- Variables set after TUI launch (set before running `ticketr tui`)
- Typos in variable names (case-sensitive)

---

### High CPU Usage

**Disable ambient effects**:
```bash
export TICKETR_EFFECTS_AMBIENT=false
```

**Reduce to minimal**:
```bash
export TICKETR_THEME=default
export TICKETR_EFFECTS_MOTION=false
```

**Monitor CPU**:
```bash
top -p $(pgrep ticketr)
```

---

### Rendering Artifacts

**Unicode characters not displaying**:
- Terminal may lack Unicode support
- Switch to default theme (has ASCII fallbacks)
- Update terminal emulator

**Shadow offset issues**:
- Reduce terminal font size
- Try different terminal emulator
- Disable shadows: `export TICKETR_EFFECTS_SHADOWS=false`

**Flickering or tearing**:
- Disable motion: `export TICKETR_EFFECTS_MOTION=false`
- Check terminal emulator performance settings
- Reduce ambient effect density (not yet configurable, file issue)

## Performance Benchmarks

Typical CPU usage on modern hardware (2020+ laptop):

| Configuration | CPU Usage | Notes |
|---------------|-----------|-------|
| Default (minimal) | < 0.1% | Idle |
| Motion only | < 0.5% | Active spinners |
| Motion + Shadows + Shimmer | < 1.5% | No ambient |
| Maximum (all effects) | < 5% | Includes ambient |
| Performance mode (no motion) | < 0.05% | Static rendering |

Memory usage: ~15-25 MB regardless of effects (no significant difference).

## Accessibility

Visual effects are designed with accessibility first:

- **Default OFF**: New users get fast, accessible interface
- **Opt-In**: Users choose their experience level
- **Global Kill Switch**: `TICKETR_EFFECTS_MOTION=false` disables all animations
- **Reduced Motion**: Honor user preferences (future enhancement)
- **Graceful Degradation**: Works on 256-color terminals, limited Unicode support
- **Legibility First**: Background effects never interfere with text readability

## Advanced Configuration (Future)

Planned for future releases:

- YAML config file support (`~/.config/ticketr/visual-effects.yaml`)
- Per-workspace effect preferences
- Custom ambient effect density and speed
- Custom color schemes and palettes
- Animation speed controls
- Keybinding to toggle effects in-app

## Related Documentation

- [TUI-GUIDE.md](TUI-GUIDE.md#visual-effects) - User guide with visual effects section
- [TUI_VISUAL_EFFECTS.md](TUI_VISUAL_EFFECTS.md) - Technical documentation for developers
- [VISUAL_EFFECTS_QUICK_START.md](VISUAL_EFFECTS_QUICK_START.md) - Integration guide for contributors

## Feedback

Help us improve the visual effects system:

- Report bugs: [GitHub Issues](https://github.com/karolswdev/ticketr/issues)
- Suggest themes: [GitHub Discussions](https://github.com/karolswdev/ticketr/discussions)
- Share screenshots: Tag `@ticketr` on social media
- Performance issues: Include terminal emulator, OS, and `TICKETR_*` env vars in report

---

**Version**: 1.0 (Phase 6, Day 12.5)
**Last Updated**: 2025-10-20
