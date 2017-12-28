package ebook

import (
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/signintech/gopdf"
)

// Pdf generate PDF file
type pdfBook struct {
	title           string
	height          float64
	pdf             *gopdf.GoPdf
	config          *gopdf.Config
	leftMargin      float64
	topMargin       float64
	w               float64
	h               float64
	maxW            float64
	maxH            float64
	titleFontSize   float64
	contentFontSize float64
	lineSpacing     float64
	fontFamily      string
	fontFile        string
}

// Info output self information
func (m *pdfBook) Info() {
	fmt.Println("generating PDF file...")
}

// SetLineSpacing dummy funciton for interface
func (m *pdfBook) SetLineSpacing(lineSpacing float64) {
	m.lineSpacing = lineSpacing
}

// SetFontFamily set custom font family
func (m *pdfBook) SetFontFamily(family string) {
	m.fontFamily = family
}

// SetFontFile set custom font file
func (m *pdfBook) SetFontFile(file string) {
	m.fontFile = file
}

// SetMargins dummy funciton for interface
func (m *pdfBook) SetMargins(left float64, top float64) {
	m.leftMargin = left
	m.topMargin = top
	m.maxW = m.w - m.leftMargin*2
	m.maxH = m.h - m.topMargin*2
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
	default:
		// work as A4 paper size
		m.config = &gopdf.Config{PageSize: gopdf.Rect{W: 595.28, H: 841.89}}
	}
	m.w = m.config.PageSize.W
	m.h = m.config.PageSize.H
	m.maxW = m.w - m.leftMargin*2
	m.maxH = m.h - m.topMargin*2
}

// SetFontSize dummy funciton for interface
func (m *pdfBook) SetFontSize(titleFontSize int, contentFontSize int) {
	m.titleFontSize = float64(titleFontSize)
	m.contentFontSize = float64(contentFontSize)
}

// Begin prepare book environment
func (m *pdfBook) Begin() {
	m.pdf = &gopdf.GoPdf{}
	m.pdf.Start(*m.config)
	m.pdf.SetLeftMargin(m.leftMargin)
	m.pdf.SetTopMargin(m.topMargin)
	m.pdf.AddPage()
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
	m.pdf.SetInfo(gopdf.PdfInfo{
		Title:        m.title,
		Author:       `GetNovel用户制作成PDF，并非小说原作者`,
		Creator:      `GetNovel，仅限个人研究学习，对其造成的所有后果，软件作者不负任何责任`,
		Producer:     `GetNovel，仅限个人研究学习，对其造成的所有后果，软件作者不负任何责任`,
		Subject:      m.title + `：不费脑子的适合电子书设备（如Kindle DXG）看的网络小说`,
		CreationDate: time.Now(),
	})
	m.pdf.WritePdf(m.title + ".pdf")
}

// AppendContent append book content
func (m *pdfBook) AppendContent(articleTitle, articleURL, articleContent string) {
	if m.height+m.titleFontSize+2 > m.maxH {
		m.pdf.AddPage()
		m.height = 0
	}
	m.pdf.SetFont(m.fontFamily, "", int(m.titleFontSize))
	m.pdf.Cell(nil, articleTitle)
	m.pdf.Br(m.titleFontSize * m.lineSpacing)
	m.height += m.titleFontSize * m.lineSpacing
	m.pdf.SetFont(m.fontFamily, "", int(m.contentFontSize))

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
func (m *pdfBook) SetTitle(title string) {
	m.title = title
}

func (m *pdfBook) writeText(t string) {
	t = `　　` + t
	count := 0
	index := 0
	for {
		r, length := utf8.DecodeRuneInString(t[index:])
		if r == utf8.RuneError {
			break
		}
		count += length
		if width, _ := m.pdf.MeasureTextWidth(t[:count]); width > m.maxW {
			if m.height+m.contentFontSize+2 > m.maxH {
				m.pdf.AddPage()
				m.height = 0
			}
			count -= length
			m.pdf.Cell(nil, t[:count])
			m.pdf.Br(m.contentFontSize * m.lineSpacing)
			m.height += m.contentFontSize * m.lineSpacing
			t = t[count:]
			index = 0
			count = 0
		} else {
			index += length
		}
	}
	if len(t) > 0 {
		if m.height+m.contentFontSize+2 > m.maxH {
			m.pdf.AddPage()
			m.height = 0
		}
		m.pdf.Cell(nil, t)
		m.pdf.Br(m.contentFontSize * m.lineSpacing)
		m.height += m.contentFontSize * m.lineSpacing
	}
}
