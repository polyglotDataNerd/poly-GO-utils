package scanner

import (
	"bufio"
	"fmt"
	aws "github.com/polyglotDataNerd/zib-Go-utils/aws"
	utils "github.com/polyglotDataNerd/zib-Go-utils/utils"
	"strings"
	"sync"
)

func ProcessObj(line chan string, bucket string, key string) {
	/* close on the sender (producer) and NOT the receiver (consumer) */
	defer close(line)
	obj, objerr := aws.S3Obj{
		bucket,
		key}.S3ReadObjGzip(aws.SessionGenerator("default", "us-west-2"))

	if objerr != nil {
		utils.Error.Fatalln(objerr.Error())
	}

	scan := bufio.NewScanner(strings.NewReader(obj))
	scan.Split(bufio.ScanLines)
	for scan.Scan() {
		line <- scan.Text()
	}

	fmt.Println("sent all lines")
}

func ProcessDir(line chan string, bucket string, key string) {
	var wg sync.WaitGroup
	defer close(line)
	/* close on the sender (producer) and NOT the receiver (consumer) */
	obj, objerr := aws.S3Obj{
		bucket,
		key}.S3ReadObjGZIPDir(aws.SessionGenerator("default", "us-west-2"))

	if objerr != nil {
		utils.Error.Fatalln(objerr.Error())
	}
	/*
	sync.WaitGroup: when spinning up many go routines without corresponding channels there needs to
	be a way to block and wait until all go routines are done. That is what the WaitGroup is doing
	for each go routine we increment a block (ADD) once the go routine is completed the DONE method will
	decrement the WaitGroup while WAIT() waits for all counters to equal 0 which means all go routines are completed
	 */
	for k, v := range obj {
		utils.Info.Println("key:", k)
		wg.Add(1)
		go func() {
			defer wg.Done()
			scan := bufio.NewScanner(strings.NewReader(v))
			scan.Split(bufio.ScanLines)
			for scan.Scan() {
				l := scan.Text()
				line <- l
			}
		}()
	}

	utils.Info.Println("sent all lines")
	wg.Wait()
}
