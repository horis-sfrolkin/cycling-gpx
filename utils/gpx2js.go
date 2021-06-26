package main

import (
	"encoding/xml"
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

type point struct {
	Lat  float64   `xml:"lat,attr"` //широта в градусах
	Lon  float64   `xml:"lon,attr"` //долгота в градусах
	Time time.Time `xml:"time"`     //Время UTC
	Dist float64   //Расстояние от предыдущей точки
}

// deg2rad преобразует значение угла из градусов в радианы
func deg2rad(deg float64) float64 {
	return deg * math.Pi / 180
}

// distance возвращает расстояние в метрах между двумя географическими точками (lat1, lon1) и (lat2, lon2)
func distance(lat1 float64, lon1 float64, lat2 float64, lon2 float64) float64 {
	const R = 6372795 //средний радиус земли в метрах
	dLat := deg2rad(lat2 - lat1)
	dLon := deg2rad(lon2 - lon1)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) + math.Cos(deg2rad(lat1))*math.Cos(deg2rad(lat2))*math.Sin(dLon/2)*math.Sin(dLon/2)
	dRad := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a)) // дистанция в радианах
	return R * dRad                                      // дистанция в метрах
}

func decodeGpxXml(r io.Reader) ([]point, error) {
	var points []point
	var prev *point
	d := xml.NewDecoder(r)
	for {
		t, err := d.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		switch t := t.(type) {
		case xml.StartElement:
			if t.Name.Local == "trkpt" {
				var p point
				if err := d.DecodeElement(&p, &t); err != nil {
					return nil, err
				}
				if prev != nil {
					p.Dist = distance(prev.Lat, prev.Lon, p.Lat, p.Lon)
					if p.Dist <= 0 {
						continue
					}
				}
				points = append(points, p)
				prev = &p
			}
		}
	}
	return points, nil
}

func outputVars(points []point, output string) error {
	if len(points) < 1 {
		return errors.New("в треке нет данных")
	}
	var w io.WriteCloser
	var err error
	if output == "" {
		w = os.Stdout
	} else {
		var fi os.FileInfo
		if fi, err = os.Stat(output); err == nil {
			if fi.IsDir() {
				t := points[0].Time
				w, err = os.Create(path.Join(output, t.Format("20060102150405")+"_"+strconv.FormatInt(t.Unix(), 10)+".json"))
			} else {
				w, err = os.Create(output)
			}
		}
	}
	if err != nil {
		return err
	}
	outArray := func(prefix string, suffix string, outItem func(point)) {
		fmt.Fprintf(w, "\"%s\":[", prefix)
		for i := 0; i < len(points); i++ {
			if i > 0 {
				fmt.Fprint(w, ",")
			}
			outItem(points[i])
		}
		fmt.Fprintf(w, "]%s", suffix)
	}
	// объект
	fmt.Fprint(w, "{")
	// кординаты
	outArray("ll", ",", func(p point) {
		fmt.Fprintf(w, "[%f,%f]", p.Lat, p.Lon)
	})
	// интервалы времени
	var prev *point
	outArray("dt", ",", func(p point) {
		var dt float64
		if prev != nil {
			dt = p.Time.Sub(prev.Time).Seconds()
		}
		prev = &p
		fmt.Fprintf(w, "%.0f", dt)
	})
	// расстояния
	outArray("dd", "", func(p point) {
		fmt.Fprintf(w, "%.2f", p.Dist)
	})
	fmt.Fprint(w, "}")
	return nil
}

func parseArgs(args []string) (files []string, dest string, err error) {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var input string
	flags.StringVar(&input, "i", "", "Имя входного GPX-файла или файловая маска GPX-файлов")
	flags.StringVar(&dest, "o", ".", "Имя выходного JSON-файла или каталога, куда будет сохранен выходной JSON-файл")
	err = flags.Parse(args)
	if err == nil {
		if input == "" {
			flags.Usage()
			err = errors.New("не указан входной файл")
		} else {
			files, err = filepath.Glob(input)
		}
	}
	return
}

func main() {
	abortIfError := func(exitCode int, err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nОшибка: %s", err.Error())
			os.Exit(exitCode)
		}
	}
	files, dest, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Ошибка в параметрах: %s", err.Error())
		os.Exit(1)
	}
	for _, file := range files {
		fmt.Print(file)
		f, err := os.Open(file)
		abortIfError(2, err)
		points, err := decodeGpxXml(f)
		f.Close()
		abortIfError(3, err)
		err = outputVars(points, dest)
		abortIfError(4, err)
		fmt.Println()
	}
}
