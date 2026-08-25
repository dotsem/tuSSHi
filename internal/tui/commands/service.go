package commands

// Service returns a command function handling :service subcommands (add, edit, rm, list).
func Service(subcmd, alias string) func(Context) {
	return func(ctx Context) {
		switch subcmd {
		case "add", "a":
			ctx.OpenServiceForm("add", nil)
		case "edit", "e":
			if alias == "" {
				ctx.SetError("Usage: :service edit <alias>")
				return
			}
			ctx.OpenServiceEdit(alias)
		case "rm", "d":
			if alias == "" {
				ctx.SetError("Usage: :service rm <alias>")
				return
			}
			ctx.DeleteService(alias)
		default:
			ctx.OpenServices()
		}
	}
}
