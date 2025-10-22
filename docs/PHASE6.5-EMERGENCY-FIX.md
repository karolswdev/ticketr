# Phase 6.5: Emergency Fix Sprint

**Status:** INITIATED
**Priority:** CRITICAL - Release Blocker
**Duration:** 2-3 days estimated
**Trigger:** Human UAT revealed catastrophic failures

---

## 🎯 For the Director: Read This First

**You are the Director.** This phase depends on your orchestration skills. The whole project depends on this.

**Your Mission:**
Fix 7 critical failures in the TUI, restore quality standards, and get human approval before release.

**You have 2-3 days.** The user is watching. Make it count.

**Required Reading:**
1. This entire document (you're reading it now)
2. `docs/TUI-WIREFRAMES-SPEC.md` - Visual design authority
3. `UAT.md` - All 7 failures documented

**Your Superpowers:**
- 6 specialized agents (Builder, Verifier, TUIUX, Scribe, Steward, Director)
- Human feedback loop
- Clear acceptance criteria
- Visual wireframes defining "done"

**Your Responsibility:**
- Orchestrate agents to fix all blockers
- Ensure quality matches wireframes
- Get human UAT approval at each milestone
- Don't proceed to Steward without human sign-off

**Let's begin.**

---

## 📚 Orchestration Methodology - Director's Guide

### The 6-Agent System

You command 6 specialized agents. Each has specific capabilities and constraints.

#### 🔨 Builder - Feature Developer
**Use for:** Implementing code, fixing bugs, integrating systems
**Tools:** Read, Write, Edit, Bash, Grep, Glob
**Strengths:** Go code, TUI implementation, integration work
**Output:** Working code, test coverage, integration summary
**Time estimate:** 4-8 hours per major feature

**When to call Builder:**
- "Investigate why pull operation freezes UI"
- "Fix ticket list population bug"
- "Wire 'n' key to workspace creation modal"
- "Integrate visual effects animator into app.go"

**What to ask for:**
- Root cause analysis
- Implementation of fix
- Self-testing results
- Handoff notes for Verifier

**Example prompt:**
```
You are the Builder agent for Phase 6.5 Day 1.

Context: Human UAT found that pressing 'P' freezes the UI during pull.

Your tasks:
1. Investigate why handlePull() causes UI freeze
2. Root cause: Is it Jira adapter blocking? Something else?
3. Implement fix to make pull truly async
4. Self-test: Run TUI and verify pull doesn't freeze
5. Report: What was broken, what you fixed, how to test

Deliverables:
- Root cause analysis document
- Fixed code (list all modified files)
- Test verification (how did you confirm it works?)
- Handoff to Verifier (what to test)

Deadline: 6 hours
```

---

#### ✅ Verifier - Quality Engineer
**Use for:** Testing implementations, finding bugs, performance validation
**Tools:** Read, Bash (tests), profiling tools
**Strengths:** Test coverage, manual testing, performance analysis
**Output:** Test reports, bug lists, sign-off decisions
**Time estimate:** 2-6 hours per verification cycle

**When to call Verifier:**
- After Builder fixes a blocker
- For integration testing
- For performance validation
- For final sign-off before human UAT

**What to ask for:**
- Manual TUI testing with real Jira
- Verification of specific fixes
- Performance measurements (CPU, FPS)
- Bug reports
- GO/NO-GO decision

**Example prompt:**
```
You are the Verifier agent for Phase 6.5 Day 1 Evening.

Context: Builder claims to have fixed UI freeze during pull operation.

Your tasks:
1. Test the fix manually with TUI + real Jira workspace
2. Verify: Press 'P', does UI freeze? (should NOT freeze)
3. Verify: Can you Tab between panes during pull? (should YES)
4. Verify: Can you cancel with ESC? (should YES)
5. Measure: CPU usage, animation FPS
6. Sign-off: Is blocker #1 FIXED or still BROKEN?

Deliverables:
- Test report (PASS/FAIL for each item)
- Any bugs found
- Performance metrics
- GO/NO-GO for proceeding to next blocker

Test against: Real Jira workspace 'tbct' with EPM project
```

---

#### ✨ TUIUX - Visual/UX Specialist
**Use for:** Animations, visual effects, UX polish, design decisions
**Tools:** Read, Write, Edit (effects code, themes, animations)
**Strengths:** Visual design, animations, user experience, TUI frameworks
**Output:** Animation code, visual specs, UX improvements
**Time estimate:** 4-8 hours per major visual feature

**When to call TUIUX:**
- Animations not working (blocker #4)
- Visual quality issues (modals, theming)
- UX improvements (keybindings overflow)
- Framework evaluation (tview vs bubbletea)

**What to ask for:**
- Investigation of animation failures
- Implementation of smooth effects
- Visual polish work
- Framework recommendations

**Example prompt:**
```
You are the TUIUX agent for Phase 6.5 Day 1.

Context: Human UAT reports NO animations visible - spinner is static, no effects.

Your tasks:
1. Investigate: Why is spinner not animating?
2. Check: Is animator started in app.go?
3. Check: Are effects enabled or disabled by default?
4. Check: Does tview support our animation requirements?
5. Fix: Make spinner animate at 80ms intervals per wireframe spec
6. Fix: Ensure progress bar animates smoothly
7. Evaluate: Can tview achieve 30-60 FPS animations?

Reference: docs/TUI-WIREFRAMES-SPEC.md (authoritative)

Deliverables:
- Root cause of animation failure
- Fixed animation code
- FPS measurements
- Recommendation: Keep tview or switch to bubbletea?

Deadline: 6 hours
```

---

#### 📝 Scribe - Documentation Specialist
**Use for:** Updating docs, fixing help text, release notes, user guides
**Tools:** Read, Write, Edit (markdown files)
**Strengths:** Clear writing, user-facing docs, technical accuracy
**Output:** Updated documentation, help screens, guides
**Time estimate:** 2-4 hours per doc update

**When to call Scribe:**
- Help documentation is wrong (issue #5)
- Need to update TUI guide
- Release notes need updates
- User-facing content

**What to ask for:**
- Audit keybindings documentation
- Update help screen
- Fix inaccurate docs
- User-focused language

---

#### 🏛️ Steward - Final Approver
**Use for:** Final quality gate, architecture review, release approval
**Tools:** Read, analysis tools
**Strengths:** Big picture view, quality standards, veto power
**Output:** GO/NO-GO decision, architectural feedback
**Time estimate:** 2-4 hours for final review

**When to call Steward:**
- ONLY after all fixes complete
- ONLY after human UAT passes
- ONLY for final release approval

**CRITICAL:** Never call Steward before human UAT passes. They approve based on human feedback.

---

### Director's Decision Framework

**You are Director.** Here's how to make decisions:

#### Decision: Which agent to call next?

**Ask yourself:**
1. What's the next blocker to fix?
2. Who can fix it? (Builder for code, TUIUX for visuals)
3. Can I parallelize? (Builder + TUIUX can work simultaneously on different blockers)
4. Have I gotten human feedback recently? (If no, get UAT before proceeding)

**Example decision tree:**
```
Current state: Day 1 morning, 4 blockers to fix

Decision: Call Builder + TUIUX in PARALLEL
- Builder: Fix blockers #1, #2, #3 (code issues)
- TUIUX: Fix blocker #4 (animations)

Rationale: These are independent, can work simultaneously
Estimated: 6 hours
Next: Verifier tests both sets of fixes
Then: Human UAT before proceeding to Day 2
```

#### Decision: When to get human feedback?

**Get human UAT:**
- ✅ After each blocker is "fixed" (verify it actually works)
- ✅ End of each day (show daily progress)
- ✅ Before calling Steward (mandatory gate)
- ✅ When uncertain about approach (ask before wasting time)

**Example:**
```
Builder reports: "Blocker #1 fixed, UI no longer freezes"

Director action:
1. Call Verifier to test fix
2. If Verifier approves: Request human UAT
3. User tests: Press 'P', confirm no freeze
4. If human approves: Move to next blocker
5. If human finds issues: Builder fixes again, repeat
```

#### Decision: Should I parallelize agents?

**Parallelize when:**
- ✅ Tasks are independent (Builder fixing code + TUIUX fixing visuals)
- ✅ No dependencies between tasks
- ✅ Both tasks well-defined

**Sequential when:**
- ✅ Tasks depend on each other (Builder must finish before Verifier can test)
- ✅ Uncertain about approach (investigate first, then fix)
- ✅ Need human decision before proceeding

**Example parallel:**
```
Day 2 Morning:
- Builder: Fix issue #5 (help docs integration)
- TUIUX: Fix issue #6 (keybindings overflow)
- Scribe: Fix issue #5 (help docs content)

All independent, can run simultaneously.
```

---

### Director's Quality Standards

**You enforce quality.** Don't accept subpar work.

#### What "DONE" means

**Code is NOT done until:**
- ✅ Builder self-tested it
- ✅ Verifier tested it
- ✅ Human UAT approved it
- ✅ Matches wireframe spec (for visual features)
- ✅ Performance budgets met (≤15% CPU)

**Documentation is NOT done until:**
- ✅ Matches actual implementation
- ✅ Human confirms it's accurate
- ✅ No false claims

**A blocker is NOT fixed until:**
- ✅ Human re-tests and confirms fix
- ✅ No regression in other areas
- ✅ Verifier approves

#### When to reject agent work

**Reject Builder work if:**
- "Fixed" but didn't actually test it
- "Fixed" but broke something else
- "Fixed" but doesn't match wireframe spec
- "Fixed" but performance is terrible (>15% CPU)

**Reject Verifier work if:**
- Only tested unit tests, not real TUI
- Only tested mocks, not real Jira
- Approved without manual testing

**Reject TUIUX work if:**
- Animations still choppy (<30 FPS)
- Doesn't match wireframe visual quality
- "Looks like the 80s" instead of modern

**How to reject:**
```
Builder, your fix for blocker #1 is NOT accepted.

Reason: You report "UI doesn't freeze" but didn't provide evidence.

Required:
1. Video or detailed log showing:
   - Press 'P'
   - UI remains responsive (Tab works during pull)
   - ESC cancels operation
2. Performance measurement (CPU usage)
3. Re-submit with evidence

Do not proceed to next blocker until this is truly fixed.
```

---

### Director's Communication Templates

#### Template: Calling an agent

```
You are the [AGENT] agent for Phase 6.5 Day [N] [Morning/Afternoon].

Context:
[What's the situation? What just happened?]

Your tasks:
1. [Specific actionable task]
2. [Specific actionable task]
3. [Specific actionable task]

Reference materials:
- docs/TUI-WIREFRAMES-SPEC.md (visual authority)
- UAT.md (all issues documented)
- [Any other relevant files]

Deliverables:
- [Specific output expected]
- [Specific output expected]
- Handoff notes for [next agent]

Acceptance criteria:
- [How I'll judge if this is done]

Deadline: [N hours]

Please begin.
```

#### Template: Requesting human UAT

```
Human, I need your UAT feedback on [specific fix/feature].

What was fixed:
- Blocker #[N]: [description]
- Agent: [Builder/TUIUX]
- Changes: [summary]

How to test:
1. Rebuild: go build -o ticketr ./cmd/ticketr
2. Launch: ./ticketr tui
3. Test: [specific steps]
4. Expected: [what should happen]

Please test and report:
- ✅ PASS or ❌ FAIL
- If FAIL: What's still broken?
- Video/screenshot if possible

This gates whether we proceed to next blocker.
```

#### Template: Daily status update

```
Phase 6.5 Day [N] - End of Day Status

Completed today:
- ✅ Blocker #[N]: [description] - FIXED, human approved
- ✅ Issue #[N]: [description] - FIXED, human approved

Still in progress:
- 🔄 Blocker #[N]: [description] - Builder working, needs testing

Blocked/Issues:
- ❌ [Any problems encountered]

Tomorrow's plan:
- [Agent]: [Task]
- [Agent]: [Task]
- Human UAT: [What to test]

Estimated completion: [On track / 1 day delay / etc]

Questions for you:
- [Any decisions needed from human?]
```

---

### Director's Anti-Patterns (DON'T DO THIS)

#### ❌ Anti-pattern: Trust without verification

**Wrong:**
```
Builder: "I fixed blocker #1, UI doesn't freeze anymore."
Director: "Great! Moving to blocker #2."
```

**Right:**
```
Builder: "I fixed blocker #1, UI doesn't freeze anymore."
Director: "Calling Verifier to test..."
Verifier: "Tested - UI still freezes. NOT FIXED."
Director: "Builder, fix is rejected. Try again."
```

#### ❌ Anti-pattern: Skipping human UAT

**Wrong:**
```
Day 3: All agents report fixes complete
Director: "Calling Steward for final approval"
Steward: "GO for release"
[Release happens, user finds bugs]
```

**Right:**
```
Day 3: All agents report fixes complete
Director: "Requesting human UAT..."
Human: "Still broken - animations choppy, modal ugly"
Director: "NO release yet. TUIUX, fix this."
[Fix, re-test, get human approval]
Director: "NOW calling Steward..."
```

#### ❌ Anti-pattern: Accepting "almost done"

**Wrong:**
```
TUIUX: "Animations work but only at 15 FPS, spec says 30-60"
Director: "Close enough, approved"
Human: "This looks choppy and bad"
```

**Right:**
```
TUIUX: "Animations work but only at 15 FPS, spec says 30-60"
Director: "NOT accepted. Spec says 30-60 FPS. Fix or explain why impossible."
TUIUX: "Investigated - tview can't do >20 FPS smoothly"
Director: "Decision: Switch to bubbletea or accept 20 FPS?"
Human: [Makes decision]
```

#### ❌ Anti-pattern: Rushing through phases

**Wrong:**
```
Day 1 Morning: Fix all 4 blockers simultaneously
Day 1 Afternoon: Release
```

**Right:**
```
Day 1: Fix 2 critical blockers
- Human UAT: Approve/reject each one
Day 2: Fix remaining blockers
- Human UAT: Full test of all fixes
Day 3: Polish + final UAT
- Steward approval only after human sign-off
```

---

### Director's Success Metrics

**You succeed when:**
- ✅ All 4 critical blockers fixed and human-approved
- ✅ Visual quality matches wireframe spec
- ✅ Performance meets budgets (≤15% CPU, 30-60 FPS)
- ✅ Human UAT passes
- ✅ Steward approves
- ✅ v3.1.1 released without major bugs

**You fail if:**
- ❌ Release happens with known bugs
- ❌ Human UAT finds same issues after "fix"
- ❌ Steward approves without human sign-off
- ❌ Visual quality still "looks like the 80s"
- ❌ Deadline extended beyond 3 days without good reason

---

### Director's Escalation Paths

#### When you're stuck

**Scenario: Builder can't fix a blocker**
1. Ask Builder to explain why it's hard
2. Research alternative approaches
3. Escalate to human: "Blocker #1 seems unfixable in tview, switch to bubbletea?"
4. Get human decision before proceeding

**Scenario: Deadline slipping**
1. Assess: Which blockers are truly critical?
2. Propose to human: "Fix 2 critical now, defer 2 to v3.1.2?"
3. Get approval before reducing scope

**Scenario: Agent gives contradictory info**
1. Cross-check with code/tests yourself
2. Call Verifier to test independently
3. Trust evidence over agent claims

---

## Executive Summary

User acceptance testing revealed **7 major failures** including 4 critical blockers that make the TUI unusable. Despite Phase 6 Week 2 being marked "complete" by all agents, the actual implementation is fundamentally broken.

**This is a complete orchestration system failure.**

**Director: You are now equipped to fix this. Use the methodology above. Make the user proud.**

---

## Critical Blockers (Must Fix)

### BLOCKER #1: UI Freezes During Pull Operation
**Severity:** CRITICAL
**Impact:** Defeats entire purpose of Phase 6 Day 6-7 async job queue
**User Experience:** Press 'P', UI freezes indefinitely, static spinner, no keyboard input

**Root Cause:** Jira adapter's `SearchTickets()` is synchronous with no context support

**Fix Required:**
- Make Jira adapter context-aware
- Add streaming progress during HTTP calls
- Implement proper cancellation
- Verify UI remains responsive during Jira operations

**Estimated Time:** 4-6 hours (Builder)

---

### BLOCKER #2: Ticket List Never Populates
**Severity:** CRITICAL
**Impact:** Core TUI feature non-functional
**User Experience:** Valid workspace selected, middle pane empty, no tickets ever load

**Root Cause:** Unknown - requires investigation

**Fix Required:**
- Investigate why LoadTickets() not populating view
- Check if tickets are being queried
- Verify view refresh logic
- Test with real workspace data

**Estimated Time:** 2-4 hours (Builder + Verifier)

---

### BLOCKER #3: 'n' Key Does Nothing
**Severity:** CRITICAL
**Impact:** Cannot create workspaces from TUI
**User Experience:** Press 'n' in workspace pane, no response, no modal

**Root Cause:** Keybinding not wired or handler broken

**Fix Required:**
- Check if 'n' is bound in workspace view
- Verify modal creation handler exists
- Wire up proper keybinding
- Test workspace creation flow end-to-end

**Estimated Time:** 1-2 hours (Builder)

---

### BLOCKER #4: Visual Effects Completely Missing
**Severity:** CRITICAL
**Impact:** Entire Day 12.5 work (2,260 lines) appears non-functional
**User Experience:** No animations, static spinners, no polish, "back in the 80s" design

**Root Cause:** Visual effects system not integrated or not enabled

**Possible Causes:**
- Animator never started
- Effects default to OFF and user didn't enable
- Integration incomplete despite agent reports
- tview rendering issues

**Fix Required:**
- Investigate why animations not visible
- Check animator initialization in app.go
- Verify effects are enabled (or enable by default)
- Test spinner animations
- Test modal shadows
- Test progress bar shimmer
- Consider if tview is the right framework

**Estimated Time:** 4-8 hours (TUIUX + Builder)

**Critical Question:** Do we need to switch to bubbletea for proper animation support?

---

## High Priority Issues (Should Fix)

### ISSUE #5: Outdated/Incorrect Help Documentation
**Severity:** HIGH
**Impact:** User trust eroded by false documentation

**Problems:**
- Claims 'W' opens credential management (doesn't exist or broken)
- Keybindings don't match implementation
- Help view is out of sync with code

**Fix Required:**
- Audit ALL keybindings in help view
- Remove non-existent features
- Update to match actual implementation
- Add missing keybindings

**Estimated Time:** 2 hours (Scribe + Builder)

---

### ISSUE #6: Keybindings Section Overflow
**Severity:** HIGH
**Impact:** Help system unusable when list too long

**Problem:**
- Keybindings section cuts off when too many actions
- No scrolling, no marquee
- User cannot see all available commands

**Fix Required:**
- Implement scrolling action bar
- Or implement marquee animation
- Or shorten keybinding list
- Prioritize most important actions

**Estimated Time:** 2-3 hours (TUIUX + Builder)

---

### ISSUE #7: Terrible Modal UX
**Severity:** MEDIUM
**Impact:** "Back in the 80s" - embarrassing quality

**Problem:**
- Modals look like raw tview widgets
- No visual polish
- No theming applied
- Shadows not visible (despite Day 12.5 work)

**Fix Required:**
- Apply theme to all modals
- Integrate ShadowBox from effects system
- Add proper borders and styling
- Make modals visually appealing

**Estimated Time:** 2-3 hours (TUIUX)

---

## Framework Question: tview vs. Bubbletea

**User Question:** "Do we need to switch to bubbletea for smooth animations?"

**Analysis Required:**

**Current:** rivo/tview
- Mature, stable
- Widget-based
- BUT: May not support smooth animations easily
- Our effects system built on top of tview

**Alternative:** charmbracelet/bubbletea
- Modern, animation-first
- Elm architecture (model-update-view)
- Built-in animation support
- Better for smooth effects
- **BUT:** Complete rewrite required

**Recommendation:**
1. **First:** Try to make animations work with tview
2. **If impossible:** Consider bubbletea migration for v3.2.0
3. **For now:** Focus on fixing blockers with current stack

**Decision Point:** Builder + TUIUX to evaluate if tview can support required animations.

---

## Orchestration Methodology for Phase 6.5

### Agent Assignment

**Day 1: Investigation & Critical Fixes**
- **Builder (8 hours):** Investigate all 4 blockers, implement fixes for #1 (async), #2 (ticket list), #3 (keybinding)
- **TUIUX (4 hours):** Investigate blocker #4 (visual effects), evaluate animation framework

**Day 2: Visual Polish & Testing**
- **Builder (4 hours):** Continue blocker fixes, integrate TUIUX fixes
- **TUIUX (6 hours):** Fix visual effects, improve modal UX, implement scrolling/marquee
- **Scribe (2 hours):** Update help documentation
- **Verifier (4 hours):** Test all fixes against real Jira, manual TUI testing

**Day 3: Integration & Human UAT**
- **Builder (2 hours):** Final integration of all fixes
- **Verifier (4 hours):** Full regression suite + manual testing
- **Human (2 hours):** Re-run UAT with fixed build
- **Steward (2 hours):** Final approval (only if human UAT passes)

---

## Acceptance Criteria

### Visual Design Reference
**All implementations must match:** `docs/TUI-WIREFRAMES-SPEC.md`

This specification defines:
- Target FPS (30-60 for animations)
- Performance budgets (≤15% CPU)
- Exact animation timings (spinner: 80ms, fade: 150ms)
- Wireframes for all states
- Color themes
- Character sets for animations

### Must Pass (Critical - Functional)
- [ ] Press 'P' → UI remains responsive during pull
- [ ] Pull operation shows animated spinner (not static!)
- [ ] Pull can be cancelled with ESC
- [ ] Ticket list populates with workspace tickets
- [ ] Press 'n' → Workspace creation modal appears

### Must Pass (Critical - Visual)
- [ ] Spinner animates smoothly at 80ms intervals (⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏)
- [ ] Progress bar fills smoothly (not jumpy)
- [ ] Modals have drop shadows (▒ characters)
- [ ] Focused borders distinct (double-line ╔═╗)
- [ ] UI looks "smooth, modern, professional" NOT "back in the 80s"

### Should Pass (High Priority)
- [ ] Help documentation matches implementation
- [ ] Keybindings section uses priority filter (no overflow)
- [ ] Modal fade-in animation (150ms)
- [ ] Progress bar shimmer effect
- [ ] Status icons visible (✓ ✗ ⚠)

### Performance Requirements
- [ ] Idle CPU ≤1%
- [ ] Active animations CPU ≤5%
- [ ] Full effects CPU ≤15%
- [ ] Frame rate 30-60 FPS
- [ ] No flickering or stuttering

### Human UAT Gate
- [ ] Human re-tests and confirms all blockers fixed
- [ ] Human approves basic usability
- [ ] Human confirms "smooth, modern, professional" visual quality
- [ ] Human confirms animations work per wireframe spec
- [ ] Human approves FPS and performance

**No Steward approval until Human UAT passes.**
**All visual implementations must match TUI-WIREFRAMES-SPEC.md wireframes.**

---

## Process Improvements for Future Phases

### What Went Wrong

1. **No Manual UAT Before Steward** - Fatal oversight
2. **Agents tested mocks, not real system** - Integration gaps
3. **Verifier only ran unit tests** - Missed end-to-end failures
4. **Builder never ran TUI against live Jira** - Critical gap
5. **TUIUX reported integration complete, but not wired** - False completion
6. **Director rushed through phases** - Should have mandated UAT
7. **Steward approved based on test reports** - Trusted automation too much

### New Mandatory Checkpoints

**For ANY Phase Touching TUI:**
1. **Manual TUI Testing Required** - Builder must run against real data
2. **Human UAT Before Steward** - Director orchestrates UAT before Day 14
3. **End-to-End Tests Required** - Verifier tests real flows, not just units
4. **Integration Validation** - Builder proves features work in real TUI
5. **Screenshot/Video Evidence** - Visual features require visual proof

---

## Phase 6.5 Execution Plan

### Day 1 Morning (4 hours)

**Builder:** Investigate blockers
1. Add debug logging to trace pull execution path
2. Identify why UI freezes (Jira adapter vs. other cause)
3. Identify why ticket list doesn't populate
4. Identify why 'n' key doesn't work
5. Report findings

**TUIUX:** Investigate visual effects
1. Check if animator started in app.go
2. Check if effects enabled/disabled
3. Test spinner animation manually
4. Identify integration gaps
5. Evaluate if tview supports smooth animations

### Day 1 Afternoon (4 hours)

**Builder:** Implement fixes
1. Fix async pull operation (context-aware Jira adapter)
2. Fix ticket list population
3. Fix 'n' key workspace creation
4. Build and self-test each fix

**TUIUX:** Fix visual effects
1. Integrate animator if not started
2. Enable effects by default (or make obvious how to enable)
3. Fix spinner animations
4. Test progress bar shimmer

### Day 2 Morning (4 hours)

**Builder:** Continue integration
1. Merge TUIUX fixes
2. Fix modal UX with theming
3. Build clean binary

**TUIUX:** Polish work
1. Implement keybindings scrolling/marquee
2. Apply theme to all modals
3. Verify visual consistency

**Scribe:** Documentation
1. Update help view with correct keybindings
2. Remove non-existent features
3. Add any missing keybindings

### Day 2 Afternoon (4 hours)

**Verifier:** Testing
1. Manual TUI testing against real Jira
2. Test all 7 issues from UAT
3. Verify async operations work
4. Verify animations visible
5. Document any remaining issues

**Builder:** Fix any issues Verifier finds

### Day 3 Morning (4 hours)

**Verifier:** Full regression
1. Run complete test suite
2. Manual smoke tests
3. Performance check
4. Sign off on fixes

**Human:** UAT Round 2
1. Test all 7 original issues
2. Test additional flows
3. Approve or reject fixes

### Day 3 Afternoon (2 hours)

**IF HUMAN UAT PASSES:**
- **Steward:** Final review and GO/NO-GO decision
- **Director:** Proceed to Day 15 release if GO

**IF HUMAN UAT FAILS:**
- **Director:** Extend Phase 6.5 by 1 day
- **Builder:** Fix remaining issues
- **Repeat:** Day 3 morning process

---

## Success Metrics

**Phase 6.5 is complete when:**
1. All 4 critical blockers fixed and verified
2. All 3 high priority issues fixed (or deferred with user approval)
3. Human UAT passes all acceptance criteria
4. Verifier confirms no regressions
5. Steward approves for release

**If these metrics not met:** Extend Phase 6.5 or consider v3.1.1 scope reduction.

---

## Communication

**Director to User:**
- Daily updates on fix progress
- Transparent about what's broken and what's fixed
- Request UAT at end of each day for early feedback
- No surprises - if we can't fix something, say so immediately

---

## Conclusion

Phase 6.5 is a **mandatory emergency fix sprint** before any v3.1.1 release. The orchestration system failed to catch these issues, and we must fix both the code and the process.

**This is our quality mark at stake.**

**Status:** Ready to orchestrate
**Next Step:** Director initiates Day 1 Morning tasks

---

**Created:** 2025-10-20
**Priority:** CRITICAL
**Blocking:** v3.1.1 Release (Phase 6 Day 15)
