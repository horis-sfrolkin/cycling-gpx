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
		<time>2021-06-03T14:00:58Z</time>
	</metadata>
	<trk>
		<name>2021-06-03T14:00:58Z</name>
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

const xml2 = `<?xml version="1.0" encoding="UTF-8"?>
<gpx xmlns="http://www.topografix.com/GPX/1/1" xsi:schemaLocation="http://www.topografix.com/GPX/1/1 http://www.topografix.com/GPX/1/1/gpx.xsd http://www.garmin.com/xmlschemas/GpxExtensions/v3 http://www.garmin.com/xmlschemas/GpxExtensionsv3.xsd http://www.garmin.com/xmlschemas/TrackPointExtension/v1 http://www.garmin.com/xmlschemas/TrackPointExtensionv1.xsd" xmlns:gpxtpx="http://www.garmin.com/xmlschemas/TrackPointExtension/v1" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance" xmlns:gpxx="http://www.garmin.com/xmlschemas/GpxExtensions/v3" version="1.1" creator="BBB">
	<metadata>
		<time>2022-05-01T10:57:18Z</time>
	</metadata>
	<trk>
		<name>2022-05-01T10:57:18Z</name>
		<trkseg>
			<trkpt lat="42.486525" lon="18.700998">
				<time>2022-05-01T10:57:18Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.48651" lon="18.701077">
				<time>2022-05-01T10:57:19Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486495" lon="18.701159">
				<time>2022-05-01T10:57:20Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486467" lon="18.701269">
				<time>2022-05-01T10:57:21Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486466" lon="18.701354">
				<time>2022-05-01T10:57:22Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486471" lon="18.701434">
				<time>2022-05-01T10:57:23Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486465" lon="18.701532">
				<time>2022-05-01T10:57:24Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486464" lon="18.701634">
				<time>2022-05-01T10:57:25Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486449" lon="18.701734">
				<time>2022-05-01T10:57:26Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486446" lon="18.701846">
				<time>2022-05-01T10:57:27Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486438" lon="18.701959">
				<time>2022-05-01T10:57:28Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486442" lon="18.702078">
				<time>2022-05-01T10:57:29Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486425" lon="18.702218">
				<time>2022-05-01T10:57:30Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486402" lon="18.702357">
				<time>2022-05-01T10:57:31Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486388" lon="18.702496">
				<time>2022-05-01T10:57:32Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486368" lon="18.70263">
				<time>2022-05-01T10:57:33Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486339" lon="18.702776">
				<time>2022-05-01T10:57:34Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486312" lon="18.702919">
				<time>2022-05-01T10:57:35Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486274" lon="18.703041">
				<time>2022-05-01T10:57:36Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486218" lon="18.703182">
				<time>2022-05-01T10:57:37Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486164" lon="18.70333">
				<time>2022-05-01T10:57:38Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486106" lon="18.703476">
				<time>2022-05-01T10:57:39Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.486031" lon="18.703612">
				<time>2022-05-01T10:57:40Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.48595" lon="18.703752">
				<time>2022-05-01T10:57:41Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.485813" lon="18.70387">
				<time>2022-05-01T10:57:42Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.485669" lon="18.703943">
				<time>2022-05-01T10:57:43Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.485557" lon="18.704045">
				<time>2022-05-01T10:57:44Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.485428" lon="18.704121">
				<time>2022-05-01T10:57:45Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.485314" lon="18.704214">
				<time>2022-05-01T10:57:46Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.485187" lon="18.704295">
				<time>2022-05-01T10:57:47Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.485067" lon="18.70439">
				<time>2022-05-01T10:57:48Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484947" lon="18.704499">
				<time>2022-05-01T10:57:49Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484871" lon="18.704696">
				<time>2022-05-01T10:57:50Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.48479" lon="18.704864">
				<time>2022-05-01T10:57:51Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484702" lon="18.705025">
				<time>2022-05-01T10:57:52Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484595" lon="18.705179">
				<time>2022-05-01T10:57:53Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484531" lon="18.705357">
				<time>2022-05-01T10:57:54Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484557" lon="18.705585">
				<time>2022-05-01T10:57:55Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484542" lon="18.705785">
				<time>2022-05-01T10:57:56Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484487" lon="18.705991">
				<time>2022-05-01T10:57:57Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484503" lon="18.706187">
				<time>2022-05-01T10:57:58Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484473" lon="18.706387">
				<time>2022-05-01T10:57:59Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484443" lon="18.706584">
				<time>2022-05-01T10:58:00Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484407" lon="18.706771">
				<time>2022-05-01T10:58:01Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484366" lon="18.706941">
				<time>2022-05-01T10:58:02Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484313" lon="18.707101">
				<time>2022-05-01T10:58:03Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484263" lon="18.707275">
				<time>2022-05-01T10:58:04Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.4842" lon="18.707452">
				<time>2022-05-01T10:58:05Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484139" lon="18.707628">
				<time>2022-05-01T10:58:06Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484096" lon="18.707797">
				<time>2022-05-01T10:58:07Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484063" lon="18.707974">
				<time>2022-05-01T10:58:08Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484032" lon="18.708144">
				<time>2022-05-01T10:58:09Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.484" lon="18.708314">
				<time>2022-05-01T10:58:10Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.483974" lon="18.708476">
				<time>2022-05-01T10:58:11Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.483956" lon="18.708638">
				<time>2022-05-01T10:58:12Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.483955" lon="18.708801">
				<time>2022-05-01T10:58:13Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.483949" lon="18.708958">
				<time>2022-05-01T10:58:14Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.483933" lon="18.709104">
				<time>2022-05-01T10:58:15Z</time>
				<extensions>
					<gpxtpx:TrackPointExtension>
					</gpxtpx:TrackPointExtension>
				</extensions>
			</trkpt>
			<trkpt lat="42.483929" lon="18.709243">
				<time>2022-05-01T10:58:16Z</time>
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
	r := strings.NewReader(xml2)
	points, _ := decodeGpxXml(r)
	outputVars(points, "")
	// Output:
	// tracks['1651402638']={"ll":[[42.486525,18.700998],[42.486510,18.701077],[42.486495,18.701159],[42.486467,18.701269],[42.486466,18.701354],[42.486471,18.701434],[42.486465,18.701532],[42.486464,18.701634],[42.486449,18.701734],[42.486446,18.701846],[42.486438,18.701959],[42.486442,18.702078],[42.486425,18.702218],[42.486402,18.702357],[42.486388,18.702496],[42.486368,18.702630],[42.486339,18.702776],[42.486312,18.702919],[42.486274,18.703041],[42.486218,18.703182],[42.486164,18.703330],[42.486106,18.703476],[42.486031,18.703612],[42.485950,18.703752],[42.485813,18.703870],[42.485669,18.703943],[42.485557,18.704045],[42.485428,18.704121],[42.485314,18.704214],[42.485187,18.704295],[42.485067,18.704390],[42.484947,18.704499],[42.484871,18.704696],[42.484790,18.704864],[42.484702,18.705025],[42.484595,18.705179],[42.484531,18.705357],[42.484557,18.705585],[42.484542,18.705785],[42.484487,18.705991],[42.484503,18.706187],[42.484473,18.706387],[42.484443,18.706584],[42.484407,18.706771],[42.484366,18.706941],[42.484313,18.707101],[42.484263,18.707275],[42.484200,18.707452],[42.484139,18.707628],[42.484096,18.707797],[42.484063,18.707974],[42.484032,18.708144],[42.484000,18.708314],[42.483974,18.708476],[42.483956,18.708638],[42.483955,18.708801],[42.483949,18.708958],[42.483933,18.709104],[42.483929,18.709243]],"dt":[0,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1,1],"dd":[0.00,7.38,7.52,7.62,7.79,7.94,8.25,8.72,8.95,9.41,9.88,10.30,10.67,10.90,11.25,11.65,12.06,12.25,12.50,12.96,13.28,13.56,13.92,14.37,14.66,14.98,15.38,15.67,15.80,15.57,15.67,15.69,15.43,15.48,15.94,16.02,16.36,16.20,16.21,16.15,16.07,16.01,15.82,15.77,15.44,15.45,15.22,15.02,14.79,14.65,14.50,14.24,13.87,13.45,13.34,13.14,12.98,12.77,12.62]}
}

func Example_decodeGpxXml3() {
	r, err := os.Open(".test/3_июня_2021 г.,_17_00.gpx")
	if err != nil {
		fmt.Println(err.Error())
	} else {
		points, err := decodeGpxXml(r)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			t0 := points[0].Time
			for i, p := range points {
				if i > 30 {
					break
				}
				t := p.Time.Sub(t0)
				fmt.Printf("%02.0f:%02.0f:%02.0f %5.2f\n", t.Hours(), t.Minutes(), t.Seconds(), p.Dist)
			}
		}
	}
	// Output:
	// 00:00:00  0.00
	// 00:00:01  2.30
	// 00:00:02  2.69
	// 00:00:03  2.92
	// 00:00:04  3.33
	// 00:00:05  3.60
	// 00:00:06  3.85
	// 00:00:07  6.34
	// 00:00:08  6.01
	// 00:00:09  5.58
	// 00:00:10  6.04
	// 00:00:12  6.32
	// 00:00:13  6.45
	// 00:00:14  6.35
	// 00:00:15  6.51
	// 00:00:16  5.33
	// 00:00:17  5.40
	// 00:00:18  5.47
	// 00:00:19  5.54
	// 00:00:20  5.55
	// 00:00:21  5.69
	// 00:00:22  5.87
	// 00:00:23  5.77
	// 00:00:24  5.64
	// 00:00:25  5.82
	// 00:00:26  5.72
	// 00:00:27  5.72
	// 00:00:28  5.37
	// 00:00:29  5.03
	// 00:00:30  4.41
	// 00:01:31  4.09
}

func Example_decodeGpxXml4() {
	r, err := os.Open("/Users/sfrolkin/Downloads/19_04_2022_11_58.gpx")
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
