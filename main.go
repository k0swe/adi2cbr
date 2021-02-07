package main

import (
	"fmt"
	"github.com/Matir/adifparser"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

func main() {
	reader := adifparser.NewADIFReader(os.Stdin)
	records := make([]adifparser.ADIFRecord, 0)
	qsoPoints := 0
	for rec, err := reader.ReadRecord(); err != io.EOF; {
		if err != nil {
			log.Fatalf("%v", err)
		}
		if rec == nil {
			continue
		}
		records = append(records, rec)
		call := safeGet(rec, "call")
		mode := safeGet(rec, "mode")
		if mode == "SSB" {
			log.Printf("%v QSO with %v is worth 1 point", mode, call)
			qsoPoints += 1
		} else if mode == "MFSK" {
			log.Printf("%v QSO with %v is worth 2 points", mode, call)
			qsoPoints += 2
		}
		rec, err = reader.ReadRecord()
	}
	log.Printf("Read %d records", reader.RecordCount())

	bandModeMultiplier := 2
	powerMultiplier := 4
	bonusPoints := 4500
	totalScore := qsoPoints*powerMultiplier*bandModeMultiplier + bonusPoints

	outputCabrillo(totalScore, records)
}

func outputCabrillo(totalScore int, records []adifparser.ADIFRecord) {
	fmt.Println("START-OF-LOG: 3.0")
	fmt.Println("Created-By: K0SWE adi2cbr v0.1")
	fmt.Println("CONTEST: WFD")
	fmt.Println("CALLSIGN: K0SWE")
	fmt.Println("LOCATION: Westminster, CO, USA")
	fmt.Println("ARRL-SECTION: CO")
	fmt.Println("CATEGORY: 1O")
	fmt.Println("CATEGORY-POWER: QRP")
	fmt.Println("SOAPBOX: 1,500 points for not using commercial power")
	fmt.Println("SOAPBOX: 1,500 points for setting up outdoors")
	fmt.Println("SOAPBOX: 1,500 points for setting up away from home")
	fmt.Println("SOAPBOX: BONUS Total 4,500")
	fmt.Printf("CLAIMED-SCORE: %d\n", totalScore)
	fmt.Println("OPERATORS: K0SWE")
	fmt.Println("NAME: Chris Keller")
	fmt.Println("ADDRESS: 5526 W 118th Ave")
	fmt.Println("ADDRESS-CITY: Westminster")
	fmt.Println("ADDRESS-STATE: CO")
	fmt.Println("ADDRESS-POSTALCODE: 80020")
	fmt.Println("ADDRESS-COUNTRY: United States")
	fmt.Println("EMAIL: xylo04@gmail.com")
	for _, rec := range records {
		if rec == nil {
			log.Print("Skipping over a nil record")
			continue
		}
		freqStr := safeGet(rec, "freq")
		mode := safeGet(rec, "mode")
		call := safeGet(rec, "call")
		qsoDate := safeGet(rec, "qso_date")
		timeOn := safeGet(rec, "time_on")
		class := safeGet(rec, "class")
		arrlSect := safeGet(rec, "arrl_sect")

		freq, err := strconv.ParseFloat(freqStr, 64)
		if err != nil {
			log.Fatalf("error converting %v to a float: %v", freqStr, err)
		}
		freq = math.Round(freq * 1000)
		switch mode {
		case "SSB":
			mode = "PH"
		case "MFSK":
			mode = "DI"
		}
		qsoDate = qsoDate[0:4] + "-" + qsoDate[4:6] + "-" + qsoDate[6:8]
		timeOn = timeOn[0:4]
		fmt.Printf("QSO: %d %v %v %v K0SWE 1O CO %v %v %v\n",
			int(freq), mode, qsoDate, timeOn, call, class, arrlSect)
	}
	fmt.Println("END-OF-LOG:")
}

func safeGet(rec adifparser.ADIFRecord, key string) string {
	value, err := rec.GetValue(key)
	if err != nil {
		log.Fatalf("error getting %v from adif: %v", key, err)
	}
	return value
}
