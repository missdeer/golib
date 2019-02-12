package ebook

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/golang/freetype/truetype"
	"github.com/signintech/gopdf"
	"github.com/signintech/gopdf/fontmaker/core"
)

// Pdf generate PDF file
type pdfBook struct {
	title           string
	height          float64
	pdf             *gopdf.GoPdf
	config          *gopdf.Config
	leftMargin      float64
	topMargin       float64
	paperWidth      float64
	paperHeight     float64
	contentWidth    float64
	contentHeight   float64
	titleFontSize   float64
	contentFontSize float64
	lineSpacing     float64
	output          string
	fontFamily      string
	fontFile        string
	pageType        string
	pagesPerFile    int
	pages           int
	chaptersPerFile int
	chapters        int
	splitIndex      int
}

// Output set the output file path
func (m *pdfBook) Output(o string) {
	m.output = o
}

// Info output self information
func (m *pdfBook) Info() {
	fmt.Println("generating PDF file...")
}

// PagesPerFile how many smaller PDF files are expected to be generated
func (m *pdfBook) PagesPerFile(n int) {
	m.pagesPerFile = n
}

// ChaptersPerFile how many smaller PDF files are expected to be generated
func (m *pdfBook) ChaptersPerFile(n int) {
	m.chaptersPerFile = n
}

// SetLineSpacing set document line spacing
func (m *pdfBook) SetLineSpacing(lineSpacing float64) {
	m.lineSpacing = lineSpacing
}

// SetFontFile set custom font file
func (m *pdfBook) SetFontFile(file string) {
	m.fontFile = file

	// check font files
	fontFd, err := os.OpenFile(m.fontFile, os.O_RDONLY, 0644)
	if err != nil {
		log.Fatalln("can't find font file", m.fontFile, err)
		return
	}

	fontContent, err := ioutil.ReadAll(fontFd)
	fontFd.Close()
	if err != nil {
		log.Fatalln("can't read font file", err)
		return
	}

	font, err := truetype.Parse(fontContent)
	if err != nil {
		log.Fatalln("can't parse TTF font", err)
		return
	}
	m.fontFamily = font.Name(truetype.NameIDFontFamily)
	m.fontFamily = ""

	// calculate Cap Height
	var parser core.TTFParser
	err = parser.Parse(m.fontFile)
	if err != nil {
		log.Print("can't parse TTF font", err)
		return
	}

	//Measure Height
	//get  CapHeight (https://en.wikipedia.org/wiki/Cap_height)
	capHeight := float64(float64(parser.CapHeight()) * 1000.00 / float64(parser.UnitsPerEm()))
	if m.lineSpacing*1000 < capHeight {
		m.lineSpacing = capHeight / 1000
	}
}

// SetMargins dummy funciton for interface
func (m *pdfBook) SetMargins(left float64, top float64) {
	m.leftMargin = left
	m.topMargin = top
	m.contentWidth = m.paperWidth - m.leftMargin*2
	m.contentHeight = m.paperHeight - m.topMargin*2
}

// SetPageType dummy funciton for interface
func (m *pdfBook) SetPageType(pageType string) {
	// https://www.cl.cam.ac.uk/~mgk25/iso-paper-ps.txt
	switch pageType {
	case "a0":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 2384, H: 3370}}
	case "a1":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 1684, H: 2384}}
	case "a2":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 1191, H: 1684}}
	case "a3":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 842, H: 1191}}
	case "a4", "dxg", "10inch":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}
	case "a5":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 420, H: 595}}
	case "a6":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 298, H: 420}}
	case "b0":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 2835, H: 4008}}
	case "b1":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 2004, H: 2835}}
	case "b2":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 1417, H: 2004}}
	case "b3":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 1001, H: 1417}}
	case "b4":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 709, H: 1001}}
	case "b5":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 499, H: 709}}
	case "b6":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 354, H: 499}}
	case "c0":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 2599, H: 3677}}
	case "c1":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 1837, H: 2599}}
	case "c2":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 1298, H: 1837}}
	case "c3":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 918, H: 1298}}
	case "c4":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 649, H: 918}}
	case "c5":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 459, H: 649}}
	case "c6":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 323, H: 459}}
	case "6inch":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 255.12, H: 331.65}} // 90 mm x 117 mm
	case "7inch":
		// FIXME
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 297.64, H: 386.93}}
	case "pc":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}
		m.SetMargins(72, 89.9)
		m.SetFontSize(16, 12)
	case "mobile":
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}
		m.SetFontSize(32, 28)
	default:
		// work as A4 paper size
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}
	}
	m.pageType = pageType
	m.paperWidth = m.config.PageSize.W
	m.paperHeight = m.config.PageSize.H
	m.contentWidth = m.paperWidth - m.leftMargin*2
	m.contentHeight = m.paperHeight - m.topMargin*2
}

// SetFontSize dummy funciton for interface
func (m *pdfBook) SetFontSize(titleFontSize int, contentFontSize int) {
	m.titleFontSize = float64(titleFontSize)
	m.contentFontSize = float64(contentFontSize)
}

// Begin prepare book environment
func (m *pdfBook) Begin() {
	m.beginBook()
	m.newPage()
}

func (m *pdfBook) beginBook() {
	m.pdf = &gopdf.GoPdf{}
	m.pdf.Start(*m.config)
	m.pdf.SetCompressLevel(9)
	m.pdf.SetLeftMargin(m.leftMargin)
	m.pdf.SetTopMargin(m.topMargin)

	if m.fontFile != "" {
		err := m.pdf.AddTTFFont(m.fontFamily, m.fontFile)
		if err != nil {
			log.Print(err.Error())
			return
		}
	}
}

// End generate files that kindlegen needs
func (m *pdfBook) End() {
	m.endBook()
}

func (m *pdfBook) endBook() {
	m.pdf.SetInfo(gopdf.PdfInfo{
		Title:        m.title,
		Author:       `golib/ebook/pdf 用户制作成PDF，并非一定是作品原作者`,
		Creator:      `golib/ebook/pdf，仅限个人研究学习，对其造成的所有后果，软件/库作者不承担任何责任`,
		Producer:     `golib/ebook/pdf，仅限个人研究学习，对其造成的所有后果，软件/库作者不承担任何责任`,
		Subject:      m.title,
		CreationDate: time.Now(),
	})
	if m.pagesPerFile > 0 || m.chaptersPerFile > 0 {
		m.splitIndex++
		m.pdf.WritePdf(fmt.Sprintf("%s_%s(%.4d).pdf", m.title, m.pageType, m.splitIndex))
	} else {
		if m.output == "" {
			m.output = fmt.Sprintf("%s_%s.pdf", m.title, m.pageType)
		}
		m.pdf.WritePdf(m.output)
	}
}

func (m *pdfBook) preprocessContent(content string) string {
	c := strings.Replace(content, `<br/>`, "\n", -1)
	c = strings.Replace(c, `&amp;`, `&`, -1)
	c = strings.Replace(c, `&lt;`, `<`, -1)
	c = strings.Replace(c, `&gt;`, `>`, -1)
	c = strings.Replace(c, `&quot;`, `"`, -1)
	c = strings.Replace(c, `&#39;`, `'`, -1)
	c = strings.Replace(c, `&nbsp;`, ` `, -1)
	c = strings.Replace(c, `</p><p>`, "\n", -1)
	for idx := strings.Index(c, "\n\n"); idx >= 0; idx = strings.Index(c, "\n\n") {
		c = strings.Replace(c, "\n\n", "\n", -1)
	}
	for len(c) > 0 && (c[0] == byte(' ') || c[0] == byte('\n')) {
		c = c[1:]
	}
	for len(c) > 0 && strings.HasPrefix(c, `　`) {
		c = c[len(`　`):]
	}
	return c
}

func (m *pdfBook) newPage() {
	if m.pages > 0 && m.pages == m.pagesPerFile {
		m.endBook()
		m.beginBook()
		m.pages = 0
	}
	m.pdf.AddPage()
	m.pages++
	m.height = 0
	m.pdf.SetFont(m.fontFamily, "", int(m.contentFontSize))
}

func (m *pdfBook) newChapter() {
	if m.chapters > 0 && m.chapters == m.chaptersPerFile {
		m.endBook()
		m.beginBook()
		m.chapters = 0
		m.pages = 0

		m.pdf.AddPage()
		m.pages++
		m.height = 0
	}
	m.chapters++
}

// AppendContent append book content
func (m *pdfBook) AppendContent(articleTitle, articleURL, articleContent string) {
	m.newChapter()
	if m.height+m.titleFontSize*m.lineSpacing > m.contentHeight {
		m.writePageNumber()
		m.newPage()
	}
	m.pdf.SetFont(m.fontFamily, "", int(m.titleFontSize))
	m.writeTextLine(articleTitle, m.titleFontSize)
	m.pdf.SetFont(m.fontFamily, "", int(m.contentFontSize))

	c := m.preprocessContent(articleContent)
	lineBreak := "\n"
	for pos := strings.Index(c, lineBreak); ; pos = strings.Index(c, lineBreak) {
		if pos <= 0 {
			if len(c) > 0 {
				m.writeText(c, m.contentFontSize)
			}
			break
		}
		t := c[:pos]
		m.writeText(t, m.contentFontSize)
		c = c[pos+len(lineBreak):]
	}
	// append a new line at the end of chapter
	if m.height+m.contentFontSize*m.lineSpacing < m.contentHeight {
		m.pdf.Br(m.contentFontSize * m.lineSpacing)
		m.height += m.contentFontSize * m.lineSpacing
	}
}

// SetTitle set book title
func (m *pdfBook) SetTitle(title string) {
	m.title = title
	m.writeCover()
	m.newPage()
}

func (m *pdfBook) writePageNumber() {
	m.pdf.SetFont(m.fontFamily, "", int(m.contentFontSize/2))
	m.pdf.SetY(m.paperHeight - m.contentFontSize/2)
	m.pdf.SetX(m.paperWidth / 2)
	m.pdf.Cell(nil, strconv.Itoa(m.pages-1))
}

func (m *pdfBook) writeCover() {
	titleOnCoverFontSize := 48
	m.pdf.SetFont(m.fontFamily, "B", titleOnCoverFontSize)
	m.pdf.SetY(m.contentHeight/2 - float64(titleOnCoverFontSize))
	m.pdf.SetX(m.leftMargin)
	m.writeText(m.title, float64(titleOnCoverFontSize))
	m.pdf.Br(float64(titleOnCoverFontSize) * m.lineSpacing)
	subtitleOnCoverFontSize := 20
	m.pdf.SetFont(m.fontFamily, "I", subtitleOnCoverFontSize)
	m.writeText(time.Now().Format(time.RFC3339), float64(subtitleOnCoverFontSize))
}

func (m *pdfBook) writeTextLine(t string, fontSize float64) {
	if e := m.pdf.Cell(nil, t); e != nil {
		fmt.Println("cell error:", e, t)
	}
	m.pdf.Br(fontSize * m.lineSpacing)
	m.height += fontSize * m.lineSpacing
}

func (m *pdfBook) writeText(t string, fontSize float64) {
	t = `　　` + t
	for index := 0; ; {
		r, length := utf8.DecodeRuneInString(t[index:])
		if r == utf8.RuneError {
			break
		}
		if length == 1 && !unicode.IsPrint(rune(t[index])) {
			t = t[:index] + t[index+1:]
			continue
		}
		index += length
	}

	count := 0
	index := 0
	for {
		r, length := utf8.DecodeRuneInString(t[index:])
		if r == utf8.RuneError {
			break
		}
		count += length
		if width, _ := m.pdf.MeasureTextWidth(t[:count]); width > m.contentWidth {
			if m.height+m.contentFontSize*m.lineSpacing > m.contentHeight {
				m.writePageNumber()
				m.newPage()
			}
			count -= length
			m.writeTextLine(t[:count], m.contentFontSize)
			t = t[count:]
			index = 0
			count = 0
		} else {
			index += length
		}
	}
	if len(t) > 0 {
		if m.height+m.contentFontSize*m.lineSpacing > m.contentHeight {
			m.writePageNumber()
			m.newPage()
		}
		m.writeTextLine(t, m.contentFontSize)
	}
}
