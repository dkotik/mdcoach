package mdcomb

type Comb struct {
	Detectors []Detector
}

// // Walk runs down a node tree, rendering each node.
// func (r *Renderer) Walk(w io.Writer, node *blackfriday.Node) {
// 	// node.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
// 	// 	return r.RenderNode(w, node, entering)
// 	// })
// }
