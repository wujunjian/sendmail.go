package main

import (
	"strings"
	"io"
	"fmt"
	"net/http"
	"io/ioutil"
	"log"
	"os"
	"time"
	"encoding/csv"
	"bytes"
	"strconv"
	util "./util"
	config "./config"
	"runtime"
	"sync"
	"regexp"
)

var wg sync.WaitGroup
var alreadySend map[string]int
func init () {
	alreadySend = make(map[string]int)
}

func main () {
	runtime.GOMAXPROCS(runtime.NumCPU())
	run()
}

func run () {
	yestoday := time.Now().Add(-86400* time.Second).Format("20060102")
	datestr := yestoday

	if len(os.Args) > 1 {
		datestr = os.Args[1]
		yestoday = os.Args[1]
	}

	fmt.Println(datestr, "begin...")
	resp, err := http.Get("http://trace-abord.cm.ijinshan.com/Indexx/ExportOpt?thever=174&date="+datestr+"&field=dumpintro&field_content=java.lang.RuntimeException: theme engine error: NULLBITMAP&field1=dumpintro&field_content1=")

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.OpenFile(datestr+".csv", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	file.Write(body)

	csvReader := csv.NewReader(bytes.NewReader(body))
	for line:=1;;line++{
		record, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if line == 1 {
			continue
		}

		if err != nil && len(record) != 6 {
			log.Fatal(err, "len(record)=", len(record))
		}

		errcount, err := strconv.Atoi(record[3])
		if  err != nil {
			log.Fatal(err)
		}

		// send mail
		if errcount>=config.Gf.Threshold {
			//fmt.Println(record[5])

			var themepkgname string
			detailstrs := strings.Fields(record[5])
			for i:=0;;i++ {
				if !strings.Contains(detailstrs[i], "packageName") {
					continue
				}
				tmppkgname := strings.Split(detailstrs[i], ":")
				themepkgname = tmppkgname[1]
				break
			}

			var validID = regexp.MustCompile(`^([a-z|A-Z|0-9]+\.[a-z|A-Z|0-9]+\.?[a-z|A-Z|0-9]?)+[a-z|A-Z|0-9]+$`)
			if ! validID.MatchString(themepkgname) {
				fmt.Println("invalid pkgname", themepkgname)
				continue
			}

			if alreadySend[themepkgname] > 0 {
				fmt.Println (themepkgname," is already send")
				continue
			} else {
				alreadySend[themepkgname] = errcount
			}


			fmt.Println("Begin to send mail about ", themepkgname)
			wg.Add(1)
			go func () {
				defer wg.Done()
				util.Sendcrashmail(record[3], themepkgname, yestoday)
			}()
		}
	}

	wg.Wait()
	fmt.Println(datestr, "end...")
}