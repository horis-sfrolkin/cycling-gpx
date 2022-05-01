package main

import (
	"bufio"
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
	"strings"
	"time"
)

type point struct {
	Lat  float64   `xml:"lat,attr"` //широта в градусах
	Lon  float64   `xml:"lon,attr"` //долгота в градусах
	Time time.Time `xml:"time"`     //Время UTC
	Dist float64   //Расстояние от предыдущей точки в метрах
}

const maxSpeed = 20   //максимальная адекватная скорость в м/с
const smoothCount = 5 //количество соседних точек назад и вперед для сглаживания
const smoothTime = 5  //интервал в секундах назад и вперед для сглаживания

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
				var dist float64
				if err := d.DecodeElement(&p, &t); err != nil {
					return nil, err
				}
				if prev != nil {
					dt := p.Time.Sub(prev.Time).Seconds()
					if dt > 0 {
						dist = distance(prev.Lat, prev.Lon, p.Lat, p.Lon)
						if dist <= 0 {
							continue
						}
						if dt > 0 && dist/dt > maxSpeed {
							continue
						}
					} else if dt == 0 { // время новой точки не изменилось - удаляем предыдущую, чтобы заменить ее
						points = points[:len(points)-1]
					} else { // время новой точки меньше предыдущей - игнорируем
						continue
					}
				}
				points = append(points, p)
				prev = &p
			}
		}
	}
	lastPointIndex := len(points) - 1
	// без сглаживания скоростей
	// for i := 1; i <= lastPointIndex; i++ {
	// 	points[i].Dist = distance(points[i-1].Lat, points[i-1].Lon, points[i].Lat, points[i].Lon)
	// }
	// // сглаживание скоростей по отдаленным точкам
	for i := 1; i <= lastPointIndex; i++ {
		l := i - smoothCount
		if l < 0 {
			l = 0
		}
		minTime := points[i].Time.Add(-time.Second * smoothTime)
		for l < (i-1) && points[l].Time.Before(minTime) {
			l++
		}
		maxTime := points[i].Time.Add(time.Second * smoothTime)
		m := i + smoothCount
		if m > lastPointIndex {
			m = lastPointIndex
		}
		for m > i && points[m].Time.After(maxTime) {
			m--
		}
		dist := distance(points[l].Lat, points[l].Lon, points[m].Lat, points[m].Lon)
		points[i].Dist = dist / float64(m-l)
	}
	return points, nil
}

func outputVars(points []point, output string) (string, error) {
	if len(points) < 1 {
		return "", errors.New("в треке нет данных")
	}
	var w io.WriteCloser
	var err error
	begTime := points[0].Time
	begTimeUnix := strconv.FormatInt(begTime.Unix(), 10)
	if output == "" {
		w = os.Stdout
	} else {
		output = path.Join(output, begTimeUnix+".js")
		w, err = os.Create(output)
		if err != nil {
			return "", err
		}
		defer w.Close()
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
	fmt.Fprintf(w, "tracks['%s']={", begTimeUnix)
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
	return output, nil
}

func parseArgs(args []string) (files []string, dest string, html string, err error) {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	var input string
	flags.StringVar(&input, "i", "", "Имя входного GPX-файла или файловая маска GPX-файлов")
	flags.StringVar(&dest, "o", ".", "Имя каталога, куда будет сохранен выходной JSON-файл")
	flags.StringVar(&html, "s", "", "Путь html-файла, в который надо вписать ссылки на json-данные поездок")
	err = flags.Parse(args)
	if err == nil {
		if input == "" {
			err = errors.New("не указан входной файл")
		} else {
			files, err = filepath.Glob(input)
		}
		if err == nil {
			var fi os.FileInfo
			fi, err = os.Stat(dest)
			if (err != nil && os.IsNotExist(err)) || (err == nil && !fi.IsDir()) {
				err = fmt.Errorf("выходной каталог '%s' не является каталогом", dest)
			}
		}
	}
	if err != nil {
		flags.Usage()
	}
	return
}

func readHtml(html string) (prevLines []string, postLines []string, err error) {
	var mode int = 0 // 0-prev 1-in 2-post
	var r *os.File
	if r, err = os.Open(html); err != nil {
		return
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		var line = scanner.Text()
		switch mode {
		case 0:
			prevLines = append(prevLines, line)
			if strings.Contains(line, "<!-- begin of routers -->") {
				mode = 1
			}
		case 1:
			if strings.Contains(line, "<!-- end of routers -->") {
				postLines = append(postLines, line)
				mode = 2
			}
		case 2:
			postLines = append(postLines, line)
		}
	}
	if err = scanner.Err(); err != nil {
		return
	}
	if mode != 2 {
		err = fmt.Errorf("в файле '%s' неверная разметка местоположения ссылок на json-данные", html)
	}
	return
}

func writeHtml(w *os.File, htmlDir string, dest string, prevLines []string, postLines []string) error {
	files, err := filepath.Glob(path.Join(dest, "*.js"))
	if err != nil {
		return err
	}
	if len(files) < 1 {
		return fmt.Errorf("в каталоге '%s' нет json-данных поездок", dest)
	}
	for _, line := range prevLines {
		fmt.Fprintln(w, line)
	}
	for _, file := range files {
		if file, err = filepath.Rel(htmlDir, file); err != nil {
			return err
		}
		if _, err = fmt.Fprintln(w, `    <script src="`+file+`"></script>`); err != nil {
			return err
		}
	}
	for _, line := range postLines {
		fmt.Fprintln(w, line)
	}
	w.Close()
	return nil
}

func main() {
	abortIfError := func(exitCode int, err error) {
		if err != nil {
			fmt.Fprintf(os.Stderr, "\nОшибка: %s", err.Error())
			os.Exit(exitCode)
		}
	}
	files, dest, html, err := parseArgs(os.Args[1:])
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
		output, err := outputVars(points, dest)
		abortIfError(4, err)
		if output != "" {
			fmt.Print(" -> " + output)
		}
		fmt.Println()
	}
	if html != "" {
		prevLines, postLines, err := readHtml(html)
		abortIfError(5, err)
		if w, err := os.Create(html); err != nil {
			abortIfError(6, err)
		} else {
			defer w.Close()
			err := writeHtml(w, filepath.Dir(html), dest, prevLines, postLines)
			abortIfError(7, err)
		}
	}
}
