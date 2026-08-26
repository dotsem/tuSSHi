package commands

type helpCommand struct{}

func (h *helpCommand) Keys() []string {
	return []string{"h", "help", "?"}
}

func (h *helpCommand) Description() string {
	return "Show this help overlay"
}

func (h *helpCommand) Execute(ctx Context, _ []string) {
	ctx.OpenHelp()
}

func init() {
	Register(&helpCommand{})
}
