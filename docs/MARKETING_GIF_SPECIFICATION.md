# Marketing GIF Specification: TUI Visual Effects Showcase

## Overview

This document specifies requirements for the marketing GIF that showcases Ticketr's visual effects system. The GIF will be used in README.md, release notes, and social media to demonstrate the "enchanting" TUI experience.

## Recording Specifications

### Technical Requirements

- **Resolution**: 1920x1080 (Full HD)
- **Frame Rate**: 60 FPS (for smooth animations)
- **Duration**: 15-20 seconds
- **Output Format**: GIF (optimized) or MP4 (with GIF fallback)
- **Color Depth**: 256 colors minimum (true color preferred)
- **File Size**: Target under 10MB for GIF, under 5MB for MP4

### Terminal Setup

- **Terminal Emulator**: iTerm2, Alacritty, or Windows Terminal (modern, true-color support)
- **Terminal Size**: 120 columns x 30 rows minimum
- **Font**: Monospace with good Unicode support (Fira Code, JetBrains Mono, Cascadia Code)
- **Font Size**: Large enough for readability when scaled (14-16pt recommended)
- **Theme**: Dark theme with hyperspace effects enabled

### Recording Tools

Recommended tools for recording:
- **macOS**: Kap, ScreenFlow, or QuickTime with terminal recording
- **Linux**: Peek, SimpleScreenRecorder, or asciinema + agg
- **Windows**: ScreenToGif, OBS Studio
- **Cross-platform**: asciinema (record) + agg (convert to GIF)

## Scene Breakdown

### Scene 1: Launch with Ambient Effects (0-3s)

**Setup**:
```bash
export TICKETR_THEME=dark
export TICKETR_EFFECTS_AMBIENT=true
export TICKETR_EFFECTS_MOTION=true
export TICKETR_EFFECTS_SHADOWS=true
ticketr tui
```

**Action**: Launch TUI, show initial screen with hyperspace background

**Showcase**:
- Hyperspace stars streaming left to right
- Clean three-panel layout
- Action bar visible at bottom

**Camera**: Static, full terminal window

---

### Scene 2: Navigation with Focus Borders (3-6s)

**Action**: Navigate between panels using Tab

**Showcase**:
- Double-line focused border on active panel
- Single-line unfocused borders on inactive panels
- Smooth border transitions
- Action bar updates showing context-aware keybindings

**Camera**: Static, full terminal window

---

### Scene 3: Modal with Drop Shadow (6-9s)

**Action**: Press 'n' to open "Create Workspace" modal

**Showcase**:
- Modal fade-in transition (░→▒→█)
- Drop shadow (▒) offset by 1 row, 2 columns
- Modal centered over background with hyperspace effect still visible
- Clean form layout

**Camera**: Static, focused on modal

---

### Scene 4: Async Operation with Progress (9-15s)

**Action**: Trigger pull operation (F2 or 'P')

**Showcase**:
- Active spinner (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏) in status bar
- Progress bar with shimmer effect: [█████░░░░░] 50% (45/90)
- Real-time percentage, count, elapsed time, ETA updates
- TUI remains responsive (show navigation during pull)
- Hyperspace background continues animating

**Camera**: Split view showing both status bar and ticket tree

---

### Scene 5: Success Completion (15-17s)

**Action**: Operation completes successfully

**Showcase**:
- Progress bar reaches 100%
- Success sparkles (✦✧⋆∗·) particle burst (500ms)
- Status message: "✓ Pull complete: 90 tickets synced"
- Brief pause to show final state

**Camera**: Focus on status bar, then pull back to full view

---

### Scene 6: Theme Transition (17-20s)

**Action**: Quick demonstration of theme switching (if possible in-app, otherwise cut/transition)

**Showcase**:
- Switch from dark (hyperspace) to arctic (snow) theme
- Snow particles falling gently from top
- Rounded borders (arctic theme characteristic)
- Cyan/white color palette

**Camera**: Static, full terminal window

---

## Post-Production

### Editing Requirements

1. **Smooth Transitions**: Use crossfades between scenes if needed
2. **Text Overlays** (optional): Brief captions explaining features
   - "Hyperspace background effect"
   - "Drop shadows on modals"
   - "Real-time progress with shimmer"
   - "Ambient effects: Optional & themeable"
3. **Speed**: Maintain real-time speed, no time-lapse or slow-motion
4. **Looping**: Ensure GIF loops seamlessly (fade out/in if needed)

### Optimization

- **Color Palette**: Reduce to 256 colors for GIF size optimization
- **Dithering**: Use Floyd-Steinberg dithering for smooth gradients
- **Frame Reduction**: If needed, reduce to 30 FPS for smaller file size
- **Compression**: Use gifsicle or similar tools to minimize file size

### Accessibility

- **Alt Text**: Provide detailed description for screen readers
- **Caption File**: Include .vtt or .srt captions describing actions
- **Static Screenshot**: Provide fallback image if GIF doesn't load

## Placeholder Content (Until GIF Recorded)

### Text Description for README

```markdown
## Visual Effects Demo

Experience Ticketr's enchanting TUI with optional visual effects:

**Subtle Motion**: Active spinners, smooth modal transitions, and focus animations bring the interface to life without distraction.

**Depth & Shadow**: Double-line focused borders, drop shadows on modals, and gradient title bars create visual hierarchy that guides your attention naturally.

**Atmospheric Effects**: Choose from themeable backgrounds—hyperspace stars streaming across dark terminals, or gentle snowfall in arctic themes. All effects are disabled by default for maximum accessibility and performance.

**Polished Details**: Success sparkles, animated progress bars with shimmer effects, and smooth transitions show craftsmanship in every interaction.

*[Marketing GIF showcasing visual effects will be added in release]*
```

### Static Screenshot Alternative

If GIF recording is delayed, capture 2-3 high-quality screenshots:

1. **Screenshot 1**: TUI main view with hyperspace background, focused panel with double-line border
2. **Screenshot 2**: Modal dialog with drop shadow visible
3. **Screenshot 3**: Progress bar with shimmer during async operation

Upload to `docs/images/tui-visual-effects-*.png`

## File Locations

Once recorded:

- **GIF**: `/home/karol/dev/private/ticktr/docs/images/tui-visual-effects-demo.gif`
- **MP4**: `/home/karol/dev/private/ticktr/docs/images/tui-visual-effects-demo.mp4`
- **Screenshots**: `/home/karol/dev/private/ticktr/docs/images/tui-effects-*.png`
- **Alt Text**: `/home/karol/dev/private/ticktr/docs/images/tui-visual-effects-alt.txt`

## README Integration

Add to README.md after "Experience" section:

```markdown
### See It In Action

![Ticketr TUI Visual Effects](docs/images/tui-visual-effects-demo.gif)

*Visual effects showcase: Hyperspace background, drop shadows, shimmer progress, and responsive async operations*
```

## Release Notes Integration

Add to release notes:

```markdown
### Visual Experience

Ticketr v3.1.1 introduces optional visual effects that transform the TUI from functional to enchanting:

![Visual Effects Demo](../docs/images/tui-visual-effects-demo.gif)

- Subtle motion and smooth animations
- Drop shadows and visual depth
- Themeable atmospheric backgrounds (hyperspace, snow)
- Polished progress indicators with shimmer

All effects are **disabled by default** and fully optional. See [TUI Guide](docs/TUI-GUIDE.md#visual-effects) for configuration.
```

## Social Media Assets

### Twitter/X Post

```
🎨 Ticketr v3.1.1: Not just functional. Beautiful.

✨ Subtle animations
🌌 Hyperspace backgrounds
📊 Shimmer progress bars
🎯 Drop shadow modals

All optional. All accessible. All enchanting.

[GIF]

#CLI #TUI #DeveloperTools #JIRA
```

### LinkedIn Post

```
Excited to share Ticketr v3.1.1—where CLI meets visual excellence!

We've transformed the terminal interface with:
• Subtle motion that feels alive
• Visual depth through shadows and gradients
• Optional atmospheric effects (hyperspace stars, snowfall)
• Polished details showing craftsmanship

All effects are disabled by default for accessibility, but when enabled, they create an experience that respects your time and rewards your attention.

Most CLI tools stop at "works." We went further.

[GIF]

#OpenSource #DeveloperExperience #TerminalUI
```

## Quality Checklist

Before publishing GIF:

- [ ] All effects visible and smooth
- [ ] No terminal artifacts or rendering glitches
- [ ] File size acceptable (under target limits)
- [ ] Loops seamlessly
- [ ] Colors accurate to actual terminal display
- [ ] Frame rate smooth (no stuttering)
- [ ] Text readable at various scales
- [ ] Alt text descriptive and accurate
- [ ] Static fallback image available
- [ ] Tested in multiple contexts (GitHub, README viewer, social media)

## Future Enhancements

Potential additions for future marketing materials:

- **Video Tutorial**: 2-3 minute walkthrough with voiceover
- **Comparison Video**: Effects disabled vs. enabled side-by-side
- **Theme Showcase**: All three themes demonstrated
- **Interactive Demo**: asciinema player embedded in docs site
- **Social Media Cuts**: 6-second clips for Twitter/Instagram

---

**Status**: Awaiting recording
**Priority**: Medium (can ship v3.1.1 with placeholder, add GIF in patch release)
**Assigned**: Marketing/Community team or volunteer contributor
**Estimated Time**: 1-2 hours for recording and editing

For questions or to volunteer for GIF creation, contact the Ticketr team via GitHub Discussions.
