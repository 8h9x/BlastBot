package common

func BoolToEmoji(b bool) string {
	if b {
		return "<a:thumbs_up:1096490549218906193>"
	}

	return "<a:thumbs_down:1096490737715130388>"
}
