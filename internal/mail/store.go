package mail

// Store supplies mail data to the app. The mock implementation returns the
// design prototype's dataset; a Gmail-backed implementation will satisfy the
// same interface later. All methods return fresh slices so callers cannot
// mutate the store's internal state.
type Store interface {
	Mails() []Mail
	Categories() []Category
	Activity() []Activity
	Rules() []Rule
	Accounts() []Account
}

// mockStore is an in-memory Store seeded with illustrative Swedish data.
type mockStore struct {
	mails      []Mail
	categories []Category
	activity   []Activity
	rules      []Rule
	accounts   []Account
}

// NewMockStore returns a Store backed by the bundled mock dataset.
func NewMockStore() Store {
	return &mockStore{
		mails:      mockMails(),
		categories: defaultCategories(),
		activity:   mockActivity(),
		rules:      mockRules(),
		accounts:   mockAccounts(),
	}
}

func (s *mockStore) Mails() []Mail          { return append([]Mail(nil), s.mails...) }
func (s *mockStore) Categories() []Category { return append([]Category(nil), s.categories...) }
func (s *mockStore) Activity() []Activity   { return append([]Activity(nil), s.activity...) }
func (s *mockStore) Rules() []Rule          { return append([]Rule(nil), s.rules...) }
func (s *mockStore) Accounts() []Account    { return append([]Account(nil), s.accounts...) }
