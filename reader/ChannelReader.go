package reader

func ReadObj(line chan string, read chan string) {
	/* read channel is a sender */
	defer close(read)
	for l := range line {
		read <- l
	}

}

func ReadChannel(channelOut chan string, input string) {
	if v := input; v == "" {
		close(channelOut)
	} else {
		channelOut <- v
	}
}
