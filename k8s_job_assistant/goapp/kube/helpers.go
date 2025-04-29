package kube

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}
