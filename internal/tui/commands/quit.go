package commands

type quitCommand struct{}

func (q *quitCommand) Keys() []string {
	return []string{"q", "quit"}
}

func (q *quitCommand) Description() string {
	return "Exit the application"
}

func (q *quitCommand) Execute(ctx Context, _ []string) {
	ctx.Quit()
}

func init() {
	Register(&quitCommand{})
}
