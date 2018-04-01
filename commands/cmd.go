package commands

type Command interface {
	ToGCode() string
}