package shared

func Check(e error) {
	if e != nil {
		Sred("%v", e)
	}
}
