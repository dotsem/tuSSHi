package commands

// New triggers opening a form for creating a new host connection.
func New() func(Context) {
	return func(ctx Context) {
		ctx.OpenForm("add")
	}
}

// Edit triggers opening a form for editing the selected host connection.
func Edit(hasFiltered bool) func(Context) {
	return func(ctx Context) {
		if hasFiltered {
			ctx.OpenForm("edit")
		}
	}
}
