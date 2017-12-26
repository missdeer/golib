package ebook

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/signintech/gopdf"
)

const (
	//595.28, 841.89 = A4
	w               = 595.28
	h               = 841.89
	titleFontSize   = 24
	contentFontSize = 18
)

var (
	// for 10 points top/left/bottom/right margins, A4
	maxW = 575.28
	maxH = 821.89
)

// Pdf generate PDF file
type Pdf struct {
	title  string
	height float64
	pdf    *gopdf.GoPdf
}

// Info output self information
func (m *Pdf) Info() {
	fmt.Println("generating PDF file for Kindle DXG...")
}

// Begin prepare book environment
func (m *Pdf) Begin() {
	m.pdf = &gopdf.GoPdf{}
	m.pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: w, H: h}})
	m.pdf.AddPage()
	err := m.pdf.AddTTFFont(`CustomFont`, "fonts/CustomFont.ttf")
	if err != nil {
		log.Print(err.Error())
		return
	}
	maxW = w - m.pdf.GetX()*2
	maxH = h - m.pdf.GetY()*2
}

// End generate files that kindlegen needs
func (m *Pdf) End() {
	m.pdf.SetInfo(gopdf.PdfInfo{
		Title:        m.title,
		Author:       `类库大魔王制作`,
		Creator:      `类库大魔王制作`,
		Producer:     `GetNovel`,
		Subject:      `不费脑子的适合Kindle DXG看的网络小说`,
		CreationDate: time.Now(),
	})
	m.pdf.WritePdf(m.title + ".pdf")
}

// AppendContent append book content
func (m *Pdf) AppendContent(articleTitle, articleURL, articleContent string) {
	if m.height+titleFontSize+2 > maxH {
		m.pdf.AddPage()
		m.height = 0
	}
	m.pdf.SetFont(`CustomFont`, "", titleFontSize)
	m.pdf.Cell(nil, articleTitle)
	m.pdf.Br(titleFontSize + 2)
	m.height += titleFontSize + 2
	m.pdf.SetFont(`CustomFont`, "", contentFontSize)

	for pos := strings.Index(articleContent, "</p><p>"); ; pos = strings.Index(articleContent, "</p><p>") {
		if pos <= 0 {
			if len(articleContent) > 0 {
				m.writeText(articleContent)
			}
			break
		}
		t := articleContent[:pos]
		m.writeText(t)
		articleContent = articleContent[pos+7:]
	}
}

// SetTitle set book title
func (m *Pdf) SetTitle(title string) {
	m.title = title
}

func (m *Pdf) writeText(t string) {
	t = `　　` + t
	count := 0
	index := 0
	for {
		r, length := utf8.DecodeRuneInString(t[index:])
		if r == utf8.RuneError {
			break
		}
		count += length
		if width, _ := m.pdf.MeasureTextWidth(t[:count]); width > maxW {
			if m.height+contentFontSize+2 > maxH {
				m.pdf.AddPage()
				m.height = 0
			}
			count -= length
			m.pdf.Cell(nil, t[:count])
			m.pdf.Br(contentFontSize + 2)
			m.height += contentFontSize + 2
			t = t[count:]
			index = 0
			count = 0
		} else {
			index += length
		}
	}
	if len(t) > 0 {
		if m.height+contentFontSize+2 > maxH {
			m.pdf.AddPage()
			m.height = 0
		}
		m.pdf.Cell(nil, t)
		m.pdf.Br(contentFontSize + 2)
		m.height += contentFontSize + 2
	}
}
