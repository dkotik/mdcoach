package mdcoach

type SlideCutPlacement uint

// TODO: dump NotesCut in favor of SlideCutKind node?
// type SlideCutTransformer func(ast.Node) SlideCutPlacement

const (
	SlideCutPlacementSkipNode = iota
	SlideCutPlacementBeforeNode
	SlideCutPlacementAfterNode
)
