package zservice

// go 协程
func Go(f func()) {
	go f()
}
