package main

import (
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
	"strings"

	drawille "github.com/Kerrigan29a/drawille-go"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func parseArgs(path, projection *string, scale *float64) {
	flag.Usage = func() {
		fmt.Println("")
		fmt.Printf("Usage:  %s [OPTIONS] -w map_csv\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Println("")
		fmt.Printf("Version: %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Date: %s\n", date)
		fmt.Println("")
	}
	projections := []string{"flat", "miller37", "miller43", "miller50"}
	flag.StringVar(path, "w", "", "World map `path`")
	flag.Float64Var(scale, "s", 1.0, "World map `scale`")
	flag.StringVar(projection, "p", "flat", "Projection `method`: "+strings.Join(projections, ", "))
	flag.Parse()

	if *path == "" {
		fmt.Fprintf(os.Stderr, "Must supply a world map file\n\n")
		flag.Usage()
		os.Exit(1)
	}
	if !stringInSlice(*projection, projections) {
		fmt.Fprintf(os.Stderr, "Projection must be one of: %s\n\n", strings.Join(projections, ", "))
		flag.Usage()
		os.Exit(1)
	}
}

func main() {
	var path string
	var projection string
	var scale float64
	parseArgs(&path, &projection, &scale)
	check(parseMap(path, projection, scale))
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func parseMap(path, projection string, scale float64) error {
	/* Open canvas */
	g := drawille.NewCanvas()
	g.Inverse = true

	/* Open file */
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	/* Open CSV reade */
	r := csv.NewReader(f)
	r.Comma = ' '
	r.Comment = '#'
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(record) != 2 {
			return errors.New("Expected latitude, longitude format")
		}

		latitude, longitude, err := parsePoint(record[1], record[0])
		if err != nil {
			return err
		}
		x, y := 0, 0
		switch projection {
		case "flat":
			y, x = calculateFlatSquareProjection(latitude, longitude, scale)
		case "miller37":
			y, x = calculateMiller37Projection(latitude, longitude, scale)
		case "miller43":
			y, x = calculateMiller43Projection(latitude, longitude, scale)
		case "miller50":
			y, x = calculateMiller50Projection(latitude, longitude, scale)
		default:
			return fmt.Errorf("unknown projection: %s", projection)
		}

		//fmt.Printf("%g N, %g E -> y:%d x:%d\n", latitude, longitude, y, x)
		g.Set(x, y)
	}
	fmt.Print(g.String())

	return nil
}

func parsePoint(latitude, longitude string) (float64, float64, error) {
	lat, err := strconv.ParseFloat(latitude, 64)
	if err != nil {
		return 0, 0, err
	}
	if lat < -90 || lat > 90 {
		return 0, 0, fmt.Errorf("wrong value for latitude: %g", lat)
	}
	long, err := strconv.ParseFloat(longitude, 64)
	if err != nil {
		return 0, 0, err
	}
	if long < -180 || long > 180 {
		return 0, 0, fmt.Errorf("wrong value for longitude: %g", long)
	}
	return lat, long, nil
}

func calculateFlatSquareProjection(latitude, longitude, scale float64) (int, int) {
	return calculateEquirectangularProjection(latitude, longitude, scale, 0, 0)
}

/* Provides minimal overall scale distortion */
func calculateMiller37Projection(latitude, longitude, scale float64) (int, int) {
	return calculateEquirectangularProjection(latitude, longitude, scale, 37.0+30.0/60.0, 0)
}

/* Provides minimal scale distortion over continents */
func calculateMiller43Projection(latitude, longitude, scale float64) (int, int) {
	return calculateEquirectangularProjection(latitude, longitude, scale, 43.0, 0)
}

/* Miller 1949 */
func calculateMiller50Projection(latitude, longitude, scale float64) (int, int) {
	return calculateEquirectangularProjection(latitude, longitude, scale, 50.0+28.0/60.0, 0)
}

func calculateEquirectangularProjection(latitude, longitude, scale, initialLatitude, initialLongitude float64) (int, int) {
	xf := (longitude - initialLongitude) * math.Cos(initialLatitude)
	yf := latitude - initialLatitude
	y := int(math.Round(yf * scale))
	x := int(math.Round(xf * scale))
	return y, x
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
