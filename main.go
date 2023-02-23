package main

func main() {
	c := make(chan *Client)
	NewClient(c)
}
