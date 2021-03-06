package gocos2d

import "sync"
import gl "github.com/mortdeus/egles/es2"

type Node interface {
	sync.Locker

	Tag() string
	IsVisible() bool
	IsDirty() (bool, bool)

	Update() error
	Draw() error
	OnEnter() error
	OnExit() error
	Cleanup() error
	Visit() error

	Transform(uint) error

	WorldToNodeTransform() AffineTransform
	NodeToWorldTransform() AffineTransform
	ParentToNodeTransform() AffineTransform
	NodeToParentTransform() AffineTransform

	ConvertTo(struct{ X, Y float32 }, uint) (struct{ X, Y float32 }, error)

	GetParent() Node
	SetParent(Node)

	AddChild(string, Node)
	GetChild(string) Node
	RemoveChild(string)
}

const (
	WorldSpace = iota
	NodeSpace
	AnchorRelative
)

type node struct {
	sync.Mutex

	tag string

	dirty,
	transformDirty bool

	visible bool

	ignoreAnchorPointForPosition,

	anchor struct{ uv, pt Point }
	rect   Rect
	rot    float32
	skew   Point
	scale  float32
	zOrder float32
	pos    Point
	size   Size
	//grid

	prog, fsh, vsh     uint
	vtouched, ftouched bool

	parent   Node
	children []Node
	lookup   map[string]int
}

func NewNode(tag string) *node {
	n := &node{
		tag:      tag,
		children: make([]Node, 0),
		lookup:   make(map[string]int, 0),
	}
	return n
}

func (n *node) Cleanup() error {
	n.Lock()
	defer n.Unlock()

	for _, child := range n.children {
		if err := child.Update(); err != nil {
			return err
		}
	}

	return nil
}

func (n *node) Update() error {
	n.Lock()
	defer n.Unlock()

	if yes, _ := n.IsDirty(); yes {
		return nil
	}
	for _, child := range n.children {
		if err := child.Update(); err != nil {
			return err
		}
	}

	return nil
}
func (n *node) Draw() error {
	n.Lock()
	defer n.Unlock()

	//if !n.IsDirty() {
	//	return nil
	//}
	for _, child := range n.children {
		if err := child.Draw(); err != nil {
			return err
		}
	}
	return nil

}

func (n *node) OnEnter() error { return nil }
func (n *node) OnExit() error  { return nil }
func (n *node) Visit() error   { return nil }

func (n *node) Transform(uint) error { return nil }
func (n *node) ConvertTo(pt struct{ X, Y float32 }, spaceType uint) (struct{ X, Y float32 }, error) {
	switch spaceType {
	case NodeSpace:
		//p := Point{pt.X, pt.Y}

		if (spaceType & AnchorRelative) != 0 {

		}
	case WorldSpace:
		if (spaceType & AnchorRelative) != 0 {
			n.ConvertTo(pt, WorldSpace)
		}

	}
	panic("unimplemented")
	return Point{0, 0}, nil
}

func (n *node) Tag() string {
	return n.tag
}

func (n *node) IsDirty() (bool, transform bool) {
	return n.dirty, n.transformDirty
}
func (n *node) IsVisible() bool {
	return n.visible
}

//GetShader returns program unless
//typ == gl.FRAGMENT_SHADER || gl.VERTEX_SHADER
func (n *node) GetShader(typ uint) uint {
	n.Lock()
	defer n.Unlock()

	switch typ {
	case gl.FRAGMENT_SHADER:
		return n.fsh
	case gl.VERTEX_SHADER:
		return n.vsh
	default:
		return n.prog
	}
	return 0
}

//SetShader sets program unless
//typ == gl.FRAGMENT_SHADER || gl.VERTEX_SHADER
func (n *node) SetShader(s, typ uint) {
	n.Lock()
	defer n.Unlock()

	switch typ {
	case gl.FRAGMENT_SHADER:
		n.fsh = s
	case gl.VERTEX_SHADER:
		n.vsh = s
	default:
		n.prog = s
		return
	}
	n.prog = Program(n.fsh, n.vsh)
}

func (n *node) SetParent(p Node) { n.parent = p }
func (n *node) GetParent() Node  { return n.parent }

func (n *node) AddChild(tag string, child Node) {
	n.Lock()
	defer n.Unlock()

	if _, exists := n.lookup[tag]; !exists {
		n.lookup[tag] = len(n.children)
		n.children = append(n.children, child)
	}
}
func (n *node) GetChild(tag string) Node {
	n.Lock()
	defer n.Unlock()
	if i, exists := n.lookup[tag]; exists {
		return n.children[i]
	}
	return nil
}

func (n *node) RemoveChild(tag string) {
	n.Lock()
	defer n.Unlock()

	/*
		Rather than physically remove the child from the list;
		this function just replaces the child with a tombstone.
		The reason we have to tombstone is because we rely on a
		map to quickly lookup the index of our children stored in
		our slice.

		One major setback of tombstoning like this is the fact that
		the node.children only grows longer and longer in algorithms
		that frequently add and remove children. Therefore we at some
		point we should consider implementing more sophisticated memory
		handling logic here.
	*/
	if i, exists := n.lookup[tag]; exists {
		n.children[i] = nil
		delete(n.lookup, tag)
	}
}

func (n *node) NodeToWorldTransform() AffineTransform {
	return nil
}
func (n *node) WorldToNodeTransform() AffineTransform {
	return nil
}
func (n *node) NodeToParentTransform() AffineTransform {
	if n.transformDirty {
		//x, y := n.pos.X, n.pos.X

		//if n.ignoreAnchorPointForPosition {

		//}
	}
	return nil
}
func (n *node) ParentToNodeTransform() AffineTransform {
	return nil
}
