package main

func ChannelDemo(ch chan any) {
	ch <- "juniper"
	ch <- "bruce"
	ch <- []string{"Hello", "World"}
	ch <- 10
}
