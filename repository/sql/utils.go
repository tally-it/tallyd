package sql

func boolToString(b bool) string {
	if b {
		return string([]byte{1})
	}

	return string([]byte{0})
}
