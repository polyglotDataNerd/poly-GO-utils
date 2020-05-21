package scanner

import (
	"bufio"
	"fmt"
	aws "github.com/polyglotDataNerd/zib-Go-utils/aws"
	utils "github.com/polyglotDataNerd/zib-Go-utils/utils"
	"strings"
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

/*
	sync.WaitGroup: when spinning up many go routines without corresponding channels there needs to
	be a way to block and wait until all go routines are done. That is what the WaitGroup is doing
	for each go routine we increment a block (ADD) once the go routine is completed the DONE method will
	decrement the WaitGroup while WAIT() waits for all counters to equal 0 which means all go routines are completed
*/

func ProcessDir(line chan string, bucket string, key string, format string) {
	//var wg sync.WaitGroup
	var source = make(map[string]string)
	defer close(line)
	object := aws.S3Obj{
		Bucket: bucket,
		Key:    key,
	}

	if format == "gzip" {
		/* close on the sender (producer) and NOT the receiver (consumer) */
		obj, objerr := object.S3ReadObjGZIPDir(aws.SessionGenerator("default", "us-west-2"))
		source = obj
		if objerr != nil {
			utils.Error.Fatalln(objerr.Error())
		}
	}
	if format == "flat" {
		/* close on the sender (producer) and NOT the receiver (consumer) */
		obj, objerr := object.S3ReadObjDir(aws.SessionGenerator("default", "us-west-2"))
		source = obj
		if objerr != nil {
			utils.Error.Fatalln(objerr.Error())
		}
	}
	for k, v := range source {
		utils.Info.Println("key:", k)
		scan := bufio.NewScanner(strings.NewReader(v))
		scan.Split(bufio.ScanLines)
		for scan.Scan() {
			l := scan.Text()
			if len(l) > 0 {
				line <- strings.ReplaceAll(l, "\n", "\t")
			}
		}
		//for scan.Scan() {
		//	wg.Add(1)
		//	time.Sleep(500 * time.Nanosecond)
		//	go func() {
		//		defer wg.Done()
		//		l := scan.Text()
		//		if len(l) > 0 {
		//			line <- strings.ReplaceAll(l, "\n", "\t")
		//		}
		//	}()
		//}
		//wg.Wait()
	}
	utils.Info.Println("sent all lines")
}
