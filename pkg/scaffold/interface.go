package scaffold

// Scaffolder interface creates files to set up a controller manager
type Scaffolder interface {
	// Scaffold performs the scaffolding
	Scaffold() error
}
