package commands

// Quit executes the application quit request.
func Quit() func(Context) {
	return func(ctx Context) {
		ctx.Quit()
	}
}
