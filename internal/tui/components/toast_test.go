package components

import (
	"strings"
	"testing"

	"github.com/charmbracelet/x/ansi"
	"github.com/stretchr/testify/assert"
)

func TestToastQueue(t *testing.T) {
	t.Run("adds toast within capacity", func(t *testing.T) {
		q := NewToastQueue()
		toast, cmd := q.Add(ActionAdd, "Host added")

		assert.NotEmpty(t, toast.ID)
		assert.Equal(t, ActionAdd, toast.Action)
		assert.Equal(t, "Host added", toast.Message)
		assert.NotNil(t, cmd)
		assert.Len(t, q.Items(), 1)
	})

	t.Run("caps at maximum 3 toasts and prunes oldest", func(t *testing.T) {
		q := NewToastQueue()
		q.Add(ActionAdd, "Toast 1")
		q.Add(ActionEdit, "Toast 2")
		q.Add(ActionDelete, "Toast 3")
		assert.Len(t, q.Items(), 3)

		q.Add(ActionInfo, "Toast 4")
		items := q.Items()
		assert.Len(t, items, 3)
		assert.Equal(t, "Toast 2", items[0].Message)
		assert.Equal(t, "Toast 3", items[1].Message)
		assert.Equal(t, "Toast 4", items[2].Message)
	})

	t.Run("dismisses specific toast by ID", func(t *testing.T) {
		q := NewToastQueue()
		t1, _ := q.Add(ActionAdd, "First")
		t2, _ := q.Add(ActionEdit, "Second")

		q.Dismiss(t1.ID)
		items := q.Items()
		assert.Len(t, items, 1)
		assert.Equal(t, t2.ID, items[0].ID)
	})

	t.Run("maps action colors correctly", func(t *testing.T) {
		tAdd := Toast{Action: ActionAdd}
		tEdit := Toast{Action: ActionEdit}
		tDelete := Toast{Action: ActionDelete}
		tInfo := Toast{Action: ActionInfo}

		assert.Equal(t, colorAdd, tAdd.BorderColor())
		assert.Equal(t, colorEdit, tEdit.BorderColor())
		assert.Equal(t, colorDelete, tDelete.BorderColor())
		assert.Equal(t, colorInfo, tInfo.BorderColor())
	})

	t.Run("overlays without panicking on empty or small buffers", func(t *testing.T) {
		q := NewToastQueue()
		assert.Equal(t, "base", q.Overlay("base", 80, 24))

		q.Add(ActionAdd, "Connected")
		res := q.Overlay("line1\nline2\nline3\nline4\nline5", 40, 5)
		assert.NotEmpty(t, res)
		assert.True(t, strings.Contains(res, "Connected"))
	})

	t.Run("preserves line widths and borders when overlaying", func(t *testing.T) {
		box := "╭───────────────────────────────────────╮\n" +
			"│ line 1                                │\n" +
			"│ line 2                                │\n" +
			"│ line 3                                │\n" +
			"│ line 4                                │\n" +
			"│ line 5                                │\n" +
			"╰───────────────────────────────────────╯"
		q := NewToastQueue()
		q.Add(ActionEdit, `Host "my-cool-server" updated success...`)
		res := q.Overlay(box, 41, 7)
		for line := range strings.SplitSeq(res, "\n") {
			assert.Equal(t, 41, ansi.StringWidth(line))
		}
	})
}
