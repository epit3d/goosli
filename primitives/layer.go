package primitives

type Layer struct {
	Order      int
	Norm       Vector
	Paths      []Path
	MiddlePs   []Path //between outer and inner (empty if wallThickness <= 2 layer)
	InnerPs    []Path //inside the body (empty if wallThickness == 1 layer)
	Fill       []Path
	PrintSpeed int // printing speed
	FanOff     bool
}
