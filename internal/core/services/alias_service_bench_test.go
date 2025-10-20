package services

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/karolswdev/ticktr/internal/core/domain"
	"github.com/karolswdev/ticktr/internal/core/ports"
)

// mockAliasRepository provides a mock repository for benchmarking.
type mockAliasRepository struct {
	aliases map[string]*domain.JQLAlias
}

func newMockAliasRepository() *mockAliasRepository {
	return &mockAliasRepository{
		aliases: make(map[string]*domain.JQLAlias),
	}
}

func (m *mockAliasRepository) Create(alias *domain.JQLAlias) error {
	key := fmt.Sprintf("%s:%s", alias.Name, alias.WorkspaceID)
	m.aliases[key] = alias
	return nil
}

func (m *mockAliasRepository) Get(id string) (*domain.JQLAlias, error) {
	for _, alias := range m.aliases {
		if alias.ID == id {
			return alias, nil
		}
	}
	return nil, ports.ErrAliasNotFound
}

func (m *mockAliasRepository) GetByName(name, workspaceID string) (*domain.JQLAlias, error) {
	key := fmt.Sprintf("%s:%s", name, workspaceID)
	if alias, ok := m.aliases[key]; ok {
		return alias, nil
	}
	return nil, ports.ErrAliasNotFound
}

func (m *mockAliasRepository) List(workspaceID string) ([]*domain.JQLAlias, error) {
	var result []*domain.JQLAlias
	for _, alias := range m.aliases {
		if alias.WorkspaceID == workspaceID {
			result = append(result, alias)
		}
	}
	return result, nil
}

func (m *mockAliasRepository) ListAll() ([]*domain.JQLAlias, error) {
	var result []*domain.JQLAlias
	for _, alias := range m.aliases {
		result = append(result, alias)
	}
	return result, nil
}

func (m *mockAliasRepository) Update(alias *domain.JQLAlias) error {
	key := fmt.Sprintf("%s:%s", alias.Name, alias.WorkspaceID)
	m.aliases[key] = alias
	return nil
}

func (m *mockAliasRepository) Delete(id string) error {
	for key, alias := range m.aliases {
		if alias.ID == id {
			delete(m.aliases, key)
			return nil
		}
	}
	return ports.ErrAliasNotFound
}

func (m *mockAliasRepository) DeleteByName(name, workspaceID string) error {
	key := fmt.Sprintf("%s:%s", name, workspaceID)
	if _, ok := m.aliases[key]; ok {
		delete(m.aliases, key)
		return nil
	}
	return ports.ErrAliasNotFound
}

// setupAliasesForBenchmark creates a set of test aliases.
func setupAliasesForBenchmark(repo *mockAliasRepository, workspaceID string) {
	now := time.Now()

	// Simple aliases (no references)
	simpleAliases := []struct {
		name string
		jql  string
	}{
		{"mywork", "assignee = currentUser() AND resolution = Unresolved"},
		{"urgent", "priority = Highest AND status != Done"},
		{"blocked", "status = Blocked"},
		{"in-progress", "status = 'In Progress'"},
		{"review", "status = 'In Review'"},
	}

	for _, a := range simpleAliases {
		alias := &domain.JQLAlias{
			ID:          uuid.New().String(),
			Name:        a.name,
			JQL:         a.jql,
			WorkspaceID: workspaceID,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		repo.Create(alias)
	}

	// Aliases with single-level references
	alias1 := &domain.JQLAlias{
		ID:          uuid.New().String(),
		Name:        "my-urgent",
		JQL:         "@mywork AND @urgent",
		WorkspaceID: workspaceID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	repo.Create(alias1)

	// Aliases with two-level references
	alias2 := &domain.JQLAlias{
		ID:          uuid.New().String(),
		Name:        "critical-work",
		JQL:         "@my-urgent AND component = Backend",
		WorkspaceID: workspaceID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	repo.Create(alias2)

	// Complex alias with multiple references
	alias3 := &domain.JQLAlias{
		ID:          uuid.New().String(),
		Name:        "focus",
		JQL:         "(@my-urgent OR @blocked) AND sprint = 'Current Sprint'",
		WorkspaceID: workspaceID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	repo.Create(alias3)
}

// BenchmarkAliasExpansion benchmarks basic alias expansion.
func BenchmarkAliasExpansion(b *testing.B) {
	repo := newMockAliasRepository()
	workspaceID := "test-workspace"
	setupAliasesForBenchmark(repo, workspaceID)

	service := NewAliasService(repo)

	b.Run("SimpleAlias", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.ExpandAlias("mywork", workspaceID)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})

	b.Run("SingleLevelReference", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.ExpandAlias("my-urgent", workspaceID)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})

	b.Run("TwoLevelReference", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.ExpandAlias("critical-work", workspaceID)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})

	b.Run("ComplexMultiReference", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.ExpandAlias("focus", workspaceID)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})
}

// BenchmarkRecursiveAliasExpansion benchmarks deeply nested alias expansion.
func BenchmarkRecursiveAliasExpansion(b *testing.B) {
	repo := newMockAliasRepository()
	workspaceID := "test-workspace"
	now := time.Now()

	// Create a chain of aliases with increasing depth
	depths := []int{3, 5, 10}

	for _, depth := range depths {
		b.Run(fmt.Sprintf("Depth%d", depth), func(b *testing.B) {
			// Clear repository
			repo.aliases = make(map[string]*domain.JQLAlias)

			// Create base alias
			baseAlias := &domain.JQLAlias{
				ID:          uuid.New().String(),
				Name:        "base",
				JQL:         "status = Open",
				WorkspaceID: workspaceID,
				CreatedAt:   now,
				UpdatedAt:   now,
			}
			repo.Create(baseAlias)

			// Create chain of aliases referencing previous level
			for i := 1; i <= depth; i++ {
				prevName := "base"
				if i > 1 {
					prevName = fmt.Sprintf("level%d", i-1)
				}

				alias := &domain.JQLAlias{
					ID:          uuid.New().String(),
					Name:        fmt.Sprintf("level%d", i),
					JQL:         fmt.Sprintf("@%s AND priority = High", prevName),
					WorkspaceID: workspaceID,
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				repo.Create(alias)
			}

			service := NewAliasService(repo)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := service.ExpandAlias(fmt.Sprintf("level%d", depth), workspaceID)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkAliasGet benchmarks alias retrieval operations.
func BenchmarkAliasGet(b *testing.B) {
	repo := newMockAliasRepository()
	workspaceID := "test-workspace"
	setupAliasesForBenchmark(repo, workspaceID)

	service := NewAliasService(repo)

	b.Run("GetUserDefined", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.Get("mywork", workspaceID)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})

	b.Run("GetPredefined", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := service.Get("open", workspaceID)
			if err != nil {
				b.Fatalf("unexpected error: %v", err)
			}
		}
	})
}

// BenchmarkAliasList benchmarks listing aliases with varying counts.
func BenchmarkAliasList(b *testing.B) {
	workspaceID := "test-workspace"
	now := time.Now()

	counts := []int{10, 50, 100, 500}

	for _, count := range counts {
		b.Run(fmt.Sprintf("Count%d", count), func(b *testing.B) {
			repo := newMockAliasRepository()

			// Create aliases
			for i := 0; i < count; i++ {
				alias := &domain.JQLAlias{
					ID:          uuid.New().String(),
					Name:        fmt.Sprintf("alias%d", i),
					JQL:         fmt.Sprintf("project = PROJ AND component = Component%d", i),
					WorkspaceID: workspaceID,
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				repo.Create(alias)
			}

			service := NewAliasService(repo)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, err := service.List(workspaceID)
				if err != nil {
					b.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// BenchmarkAliasCreate benchmarks alias creation.
func BenchmarkAliasCreate(b *testing.B) {
	workspaceID := "test-workspace"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		repo := newMockAliasRepository()
		service := NewAliasService(repo)
		b.StartTimer()

		err := service.Create(
			fmt.Sprintf("test-alias-%d", i),
			"status = Open",
			"Test alias",
			workspaceID,
		)
		if err != nil {
			b.Fatalf("unexpected error: %v", err)
		}
	}
}

// BenchmarkStringOperationsInExpansion benchmarks string operations during expansion.
func BenchmarkStringOperationsInExpansion(b *testing.B) {
	jql := "assignee = currentUser() AND @urgent AND status != Done AND @blocked OR @review"

	b.Run("ContainsCheck", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// This is what the current implementation does
			_ = len(jql) > 0 && jql[0] == '@'
		}
	})

	b.Run("StringsFields", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			// Current implementation uses strings.Fields
			words := make([]string, 0, 10)
			start := 0
			inSpace := true
			for i, r := range jql {
				if r == ' ' || r == '\t' || r == '\n' {
					if !inSpace && i > start {
						words = append(words, jql[start:i])
					}
					inSpace = true
				} else {
					if inSpace {
						start = i
					}
					inSpace = false
				}
			}
			if !inSpace && len(jql) > start {
				words = append(words, jql[start:])
			}
			_ = words
		}
	})
}
