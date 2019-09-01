package gcode

type Gcode struct {
	Cmds        []Command
	LayersCount int
}

func (g *Gcode) Add(cmd Command) {
	g.LayersCount += cmd.LayersCount()
	g.Cmds = append(g.Cmds, cmd)
}
