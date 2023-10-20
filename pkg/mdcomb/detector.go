package mdcomb

import blackfriday "github.com/russross/blackfriday/v2"

type Detector interface {
	Detect(*blackfriday.Node) Renderer
}
