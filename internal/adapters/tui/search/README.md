# TUI Search Package

This package provides fuzzy text matching and filtering logic for the Ticketr TUI interface.

## Architecture

- **Pure Go**: No external dependencies except standard library `regexp`
- **Stateless**: All functions are pure (no global state)
- **Thread-safe**: No shared mutable state
- **Performance**: < 50ms for 1000 tickets (actual: ~24ms)

## Files

- `matcher.go` (137 lines) - Fuzzy text matching and relevance scoring
- `filter.go` (231 lines) - Query parsing and ticket filtering
- `matcher_test.go` (343 lines) - Matcher tests
- `filter_test.go` (472 lines) - Filter tests
- `example_test.go` (227 lines) - Usage examples and documentation

## API Overview

### Fuzzy Matching

```go
// FuzzyMatch returns a Match with score 0-100 and match locations
match, found := search.FuzzyMatch(ticket, "authentication")
// Returns: Match{Score: 80, MatchedIn: ["title"]}
```

**Scoring Algorithm:**
- Exact title match: 100 points
- Title contains query: 80 points
- JiraID exact match: 90 points
- JiraID contains query: 70 points
- Description contains query: 60 points
- CustomFields contain query: 40 points
- AcceptanceCriteria contain query: 30 points

All matching is **case-insensitive** and uses **partial word matching**.

### Search

```go
// SearchTickets returns matches sorted by relevance (descending)
results := search.SearchTickets(tickets, "auth")
// Returns: []*Match sorted by Score (highest first)
```

### Query Parsing

```go
query, _ := search.ParseQuery("@john #BACK-123 !high ~sprint23 /auth.*/ login bug")

// Result:
// query.Text = "login bug"
// query.Assignee = "john"
// query.JiraID = "BACK-123"
// query.Priority = "high"
// query.Sprint = "sprint23"
// query.Regex = "auth.*"
```

**Filter Syntax:**
- `@username` - Filter by Assignee custom field
- `#JIRA-ID` - Filter by JiraID (partial match, case-insensitive)
- `!priority` - Filter by Priority custom field
- `~sprint` - Filter by Sprint custom field
- `/pattern/` - Regex filter on Title and Description
- `%XX` - Completion filter (not implemented, silently ignored)

### Filtering

```go
// ApplyFilters applies parsed query filters to ticket list
filtered := search.ApplyFilters(tickets, query)
```

Filters are applied in this order:
1. JiraID filter
2. Assignee filter
3. Priority filter
4. Sprint filter
5. Regex filter

All filters use **case-insensitive partial matching** except regex (uses regex semantics).

## Full Search Pipeline

The typical usage pattern is:

```go
// 1. Parse user input
query, _ := search.ParseQuery("@john !high authentication")

// 2. Apply filters
filtered := search.ApplyFilters(tickets, query)

// 3. Fuzzy search on remaining text
results := search.SearchTickets(filtered, query.Text)

// 4. Display results (sorted by relevance)
for _, match := range results {
    fmt.Printf("%s: %s (score: %d)\n",
        match.Ticket.JiraID,
        match.Ticket.Title,
        match.Score)
}
```

## Performance Characteristics

Benchmarked on AMD Ryzen AI 7 350 (16 threads):

| Operation | Time | Throughput |
|-----------|------|------------|
| ParseQuery | 1.77 μs | 565k ops/sec |
| FuzzyMatch (single) | 0.52 μs | 1.9M ops/sec |
| ApplyFilters (1000 tickets, 2 filters) | 23 μs | 43k ops/sec |
| SearchTickets (1000 tickets) | 368 μs | 2.7k ops/sec |
| Full pipeline (1000 tickets) | 24 μs | 41k ops/sec |

**Memory:**
- ParseQuery: 2.3 KB, 33 allocations
- FuzzyMatch: 368 B, 11 allocations
- Full search (1000 tickets): 14 KB, 369 allocations

## Edge Cases Handled

- ✅ Nil tickets → returns empty slice
- ✅ Nil ticket entries → skipped
- ✅ Empty query → matches all with neutral score (50)
- ✅ Malformed filters (@, #, !, ~) → ignored
- ✅ Invalid regex → returns all tickets unchanged
- ✅ Missing custom fields → ticket excluded from filter
- ✅ Unicode and special characters → handled correctly
- ✅ Very long queries → no panics or timeouts

## Error Handling

All functions are designed to be **panic-free** and return sensible defaults:

- Invalid regex: Returns all tickets (no filtering)
- Nil inputs: Returns empty slices
- Malformed queries: Skips malformed tokens
- Missing fields: Excludes tickets from filter results

The `ParseQuery` function returns an error parameter but currently never returns an error (designed for future validation).

## Integration with TUI

This package provides the **logic layer** only. TUI views should:

1. Collect user input from search field
2. Parse with `ParseQuery()`
3. Apply filters with `ApplyFilters()`
4. Search with `SearchTickets()`
5. Display results sorted by `Match.Score`

The TUI can also use `Match.MatchedIn` to highlight where the match was found.

## Testing

```bash
# Run all tests
go test ./internal/adapters/tui/search/...

# Run with coverage
go test ./internal/adapters/tui/search/... -cover

# Run benchmarks
go test ./internal/adapters/tui/search/... -bench=. -benchmem

# Run examples (validates output)
go test ./internal/adapters/tui/search/... -run Example
```

## Future Enhancements

Potential improvements (not in Week 14 scope):

- [ ] Weighted scoring configuration
- [ ] Trigram matching for better fuzzy search
- [ ] Completion percentage filter (`%50`)
- [ ] Boolean operators (AND, OR, NOT)
- [ ] Field-specific search (`title:auth`, `desc:bug`)
- [ ] Date range filters
- [ ] Caching/indexing for very large datasets
