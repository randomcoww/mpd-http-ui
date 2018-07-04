package main


func main() {
	err := NewDataFeeder()

	if err != nil {
		panic(err)
	}
}
