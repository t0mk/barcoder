package main

import (
	"fmt"
	"image/gif"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	yaml "gopkg.in/yaml.v2"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/code128"
	"github.com/jung-kurt/gofpdf"
	"github.com/urfave/cli"
)

func readfile(p string) []byte {
	filename, _ := filepath.Abs("./" + p)
	content, err := ioutil.ReadFile(filename)

	if err != nil {
		panic(err)
	}
	return content
}

func stripchars(str, chr string) string {
	return strings.Map(func(r rune) rune {
		if strings.IndexRune(chr, r) < 0 {
			return r
		}
		return -1
	}, str)
}

type Payment struct {
	Name   string
	IBAN   string
	Amount float64
	Ref    string
	Date   string `yaml:"omitempty"`
}

//func barcod(iban, eur, refnum, ddate string) string {
func barcod(p Payment) string {
	log.Printf("%+v", p)
	// This is something in Finland... I don't even remember where I
	// learned this from
	codeIBAN := "4" + p.IBAN[2:]
	codeEUR := fmt.Sprintf("%08s", stripchars(strconv.FormatFloat(p.Amount, 'f', 2, 64), "."))
	codeRefNum := fmt.Sprintf("%023s", p.Ref)
	codeDate := stripchars(p.Date, "-")[2:]
	return codeIBAN + codeEUR + codeRefNum + codeDate

}

func codeToFile(c, fn string) error {
	e, err := code128.Encode(c)
	if err != nil {
		return err
	}

	cd, err := barcode.Scale(e, 800, 200)
	if err != nil {
		return err
	}

	file, _ := os.Create(fn)
	defer file.Close()

	gif.Encode(file, cd, &gif.Options{NumColors: 256})
	return nil

}

func barcodesPrint(c *cli.Context) error {
	templ := c.String("templ")
	if templ == "" {
		panic(fmt.Errorf("You must supply template file"))
	}
	outfile := c.String("outfile")
	if outfile == "" {
		panic(fmt.Errorf("You must supply filename for result pdf"))
	}
	yamlFile := readfile(templ)

	var ps []Payment

	err := yaml.Unmarshal(yamlFile, &ps)
	if err != nil {
		panic(err)
	}

	pdf, err := getPdf(ps, c.String("date"))
	if err != nil {
		panic(err)
	}
	err = pdf.OutputFileAndClose(outfile)
	if err != nil {
		panic(err)
	}

	return nil
}

func main() {
	now := time.Now()
	app := &cli.App{
		Action: barcodesPrint,
		Name:   "barcodesPrint",
		Usage:  "generates PDF invoice",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "templ"},
			&cli.StringFlag{Name: "outfile"},
			&cli.StringFlag{Name: "date", Value: now.Add(168 * time.Hour).String()[:10]},
		},
	}
	app.Run(os.Args)
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	rand.Seed(time.Now().UTC().UnixNano())

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func getPdf(ps []Payment, date string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "", 12)
		pdf.CellFormat(0, 10, "generated with https://github.com/t0mk/barcoder", "", 0, "C", false, 0, "")
	})

	pdf.AddPage()
	y := 10.0
	tr := pdf.UnicodeTranslatorFromDescriptor("") // "" defaults to "cp1252"
	for i, p := range ps {
		if y > 230.0 {
			pdf.AddPage()
			y = 10.0
		}

		p.Date = date
		pdf.SetXY(20, y)
		pdf.SetFont("Helvetica", "B", 12)
		pdf.Cell(10, 30, tr(p.Name))
		pdf.SetFont("Helvetica", "", 12)

		pdf.SetXY(20, y+5)
		pdf.Cell(10, 30, "Summa:")

		pdf.SetXY(40, y+5)
		pdf.Cell(10, 30, strconv.FormatFloat(p.Amount, 'f', 2, 64)+" EUR")

		pdf.SetXY(20, y+10)
		pdf.Cell(10, 30, "IBAN:")

		pdf.SetXY(40, y+10)
		pdf.Cell(10, 30, p.IBAN)

		pdf.SetXY(90, y+5)
		pdf.Cell(10, 30, "Viite:")

		pdf.SetXY(120, y+5)
		pdf.Cell(10, 30, p.Ref)

		pdf.SetXY(90, y+10)
		pdf.Cell(10, 30, tr("Päivämäärä: "))

		pdf.SetXY(120, y+10)
		pdf.Cell(10, 30, p.Date)

		var opt gofpdf.ImageOptions
		bcname := fmt.Sprintf("%d%s.gif", i, randString(8))
		codeToFile(barcod(p), bcname)

		opt.ImageType = "gif"
		pdf.ImageOptions(bcname, 8, y+30, 160, 20, false, opt, 0, "")

		y += 50

	}

	pdf.AliasNbPages("{nb}") // replace {nb}

	return pdf, nil
}
