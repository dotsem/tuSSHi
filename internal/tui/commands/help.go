package commands

// Help triggers opening the interactive help dialog.
func Help() func(Context) {
	return func(ctx Context) {
		ctx.OpenHelp()
	}
}
