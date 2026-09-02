package commands

type serviceCmd struct{}

func (c *serviceCmd) Keys() []string      { return []string{"service", "services", "svc"} }
func (c *serviceCmd) Description() string { return "Manage SSH services (:svc [add|edit|rm])" }
func (c *serviceCmd) Execute(ctx Context, parts []string) {
	subcmd := ""
	alias := ""
	if len(parts) > 1 {
		subcmd = parts[1]
	}
	if len(parts) > 2 {
		alias = parts[2]
	}

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

func init() {
	Register(&serviceCmd{})
}
