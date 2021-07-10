package main

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const xml1 = `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" version="1.1" creator="BBB">
	<metadata>
		<time>2021-06-03T14:00:59Z</time>
	</metadata>
	<trk>
		<name>2021-06-03T14:00:59Z</name>
		<trkseg>
			<trkpt lat="59.907581" lon="30.256245">
				<time>2021-06-03T14:00:59Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="59.907581" lon="30.256245">
				<time>2021-06-03T14:01:00Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="59.90762" lon="30.256319">
				<time>2021-06-03T14:01:01Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="59.907591" lon="30.256423">
				<time>2021-06-03T14:01:05Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
		</trkseg>
	</trk>
</gpx>
`

func Example_decodeGpxXml1() {
	r := strings.NewReader(xml1)
	points, _ := decodeGpxXml(r)
	outputVars(points, "")
	// Output:
	// tracks['1622728859']={"ll":[[59.907581,30.256245],[59.907620,30.256319],[59.907591,30.256423]],"dt":[0,2,4],"dd":[0.00,5.99,6.64]}
}

func Example_decodeGpxXml2() {
	r, err := os.Open(".test/2_июня_2021 г.,_15_19.gpx")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		points, err := decodeGpxXml(r)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			outputVars(points, ".test")
			fmt.Println(len(points))
		}
	}
	// Output:
	// 1244
}

func Example_main() {
	os.Args = append([]string{os.Args[0]}, "-i=.test\\*.gpx", "-o=.test")
	main()
	// Output:
	// .test\2_июня_2021 г.,_15_19.gpx -> .test/1622636398.js
	// .test\3_июня_2021 г.,_17_00.gpx -> .test/1622728859.js
}

func Test_distance(t *testing.T) {
	tst := func(lat1 float64, lon1 float64, lat2 float64, lon2 float64, goal float64) {
		d := distance(lat1, lon1, lat2, lon2)
		if math.Abs(d-goal) > 1 {
			t.Fail()
		}
	}
	tst(77.1539, -139.398, -77.1804, -139.55, 17166029)
	tst(77.1539, 120.398, 77.1804, 129.55, 225883)
	tst(77.1539, -120.398, 77.1804, 129.55, 2332669)
}

func Example_parseArgs() {
	tst := func(args []string) {
		files, dest, html, err := parseArgs(args)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(files)
			fmt.Println(dest)
			fmt.Println(html)
		}
	}
	tst([]string{})
	tst([]string{"-h"})
	tst([]string{"-i=.test/2_июня_2021 г.,_15_19.gpx", "-o=.test"})
	tst([]string{"-i=.test/2_июня_2021 г.,_15_19.gpx", "-o=.test", "-s=.test\\index.html"})
	tst([]string{"-i=.test\\*.gpx"})
	tst([]string{"-i=.test/2_июня_2021 г.,_15_19.gpx", "-o=.test\\xxx"})
	tst([]string{"-i=.test/2_июня_2021 г.,_15_19.gpx", "-o=.test\\1622636398.js"})

	// Output:
	// не указан входной файл
	// flag: help requested
	// [.test/2_июня_2021 г.,_15_19.gpx]
	// .test

	// [.test/2_июня_2021 г.,_15_19.gpx]
	// .test
	// .test\index.html
	// [.test\2_июня_2021 г.,_15_19.gpx .test\3_июня_2021 г.,_17_00.gpx]
	// .

	// выходной каталог '.test\xxx' не является каталогом
	// выходной каталог '.test\1622636398.js' не является каталогом
}

func Example_readHtml() {
	prevLines, postLines, err := readHtml(".test\\index.html")
	if err != nil {
		return
	}
	for _, line := range prevLines {
		fmt.Println(line)
	}
	for _, line := range postLines {
		fmt.Println(line)
	}

	// Output:
	// <!DOCTYPE html>
	// <html>
	// <head>
	//     <title>Велопоездки</title>
	//     <meta charset="utf-8" />
	//     <link rel="stylesheet" href="css/leaflet.css">
	//     <link rel="stylesheet" href="css/cycling-gpx.css">
	//     <script src="js/jquery-3.6.0.min.js"></script>
	//     <script src="js/leaflet.js"></script>
	//     <script>var tracks = {}</script>
	//     <!-- begin of routers -->
	//     <!-- end of routers -->
	// </head>
	// <body>
	//     <div id="mapid"></div>
	//     <select id="tracks"></select>
	//     <script src="js/cycling-gpx.js"></script>
	// </body>
	// </html>
}

func Example_writeHtml() {
	html := ".test\\index.html"
	dest := ".test"
	prevLines, postLines, err := readHtml(html)
	if err != nil {
		return
	}
	w := os.Stdout
	err = writeHtml(w, filepath.Dir(html), dest, prevLines, postLines)
	if err != nil {
		return
	}

	// Output:
	// <!DOCTYPE html>
	// <html>
	// <head>
	//     <title>Велопоездки</title>
	//     <meta charset="utf-8" />
	//     <link rel="stylesheet" href="css/leaflet.css">
	//     <link rel="stylesheet" href="css/cycling-gpx.css">
	//     <script src="js/jquery-3.6.0.min.js"></script>
	//     <script src="js/leaflet.js"></script>
	//     <script>var tracks = {}</script>
	//     <!-- begin of routers -->
	//     <script src="1622636398.js"></script>
	//     <script src="1622728859.js"></script>
	//     <!-- end of routers -->
	// </head>
	// <body>
	//     <div id="mapid"></div>
	//     <select id="tracks"></select>
	//     <script src="js/cycling-gpx.js"></script>
	// </body>
	// </html>
}
