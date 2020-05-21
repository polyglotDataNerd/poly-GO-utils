package reader

import (
	log "github.com/polyglotDataNerd/zib-Go-utils/utils"
	"time"
)
func ReadObj(line chan string, read chan string) {
	/* read channel is a sender */
	defer close(read)
	for l := range line {
		read <- l
	}

}

func ReadSelect(sender chan string, receiver chan string, status chan int) {
	var l string
	ok := true
	for ok {
		select {
		case l := <-sender:
			log.Info.Println(l)
		case receiver <- l:
		case <-status:
			ok = false
			return
		default:
			time.Sleep(500 * time.Nanosecond)
		}
	}

}

func ReadChannel(channelOut chan string, input string) {
	if v := input; v == "" {
		close(channelOut)
	} else {
		channelOut <- v
	}
}

