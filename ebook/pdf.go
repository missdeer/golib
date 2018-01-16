package ebook

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/golang/freetype/truetype"
	"github.com/missdeer/gopdf"
	pdf "github.com/unidoc/unidoc/pdf/model"
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
	implicitMerge   bool
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
	defer fontFd.Close()

	fontContent, err := ioutil.ReadAll(fontFd)
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
	// if m.pagesPerFile == 0 && m.chaptersPerFile == 0 {
	// 	m.implicitMerge = true
	// 	m.chaptersPerFile = 10
	// }
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

func (m *pdfBook) mergeFile(inputPath string, pdfWriter *pdf.PdfWriter) error {
	f, err := os.Open(inputPath)
	if err != nil {
		return err
	}

	defer f.Close()

	pdfReader, err := pdf.NewPdfReader(f)
	if err != nil {
		return err
	}

	numPages, err := pdfReader.GetNumPages()
	if err != nil {
		return err
	}

	for i := 0; i < numPages; i++ {
		pageNum := i + 1

		page, err := pdfReader.GetPage(pageNum)
		if err != nil {
			return err
		}

		err = pdfWriter.AddPage(page)
		if err != nil {
			return err
		}
	}
	return nil
}

// End generate files that kindlegen needs
func (m *pdfBook) End() {
	m.endBook()
	if m.implicitMerge {
		pdfWriter := pdf.NewPdfWriter()

		var inputPaths []string
		for i := 1; ; i++ {
			inputPath := fmt.Sprintf("%s_%s(%.4d).pdf", m.title, m.pageType, i)
			if _, err := os.Stat(inputPath); os.IsNotExist(err) {
				break
			}
			inputPaths = append(inputPaths, inputPath)
		}

		for _, inputPath := range inputPaths {
			if err := m.mergeFile(inputPath, &pdfWriter); err != nil {
				log.Println("merge", inputPath, "failed:", err)
			}
			os.Remove(inputPath)
		}

		if m.output == "" {
			m.output = fmt.Sprintf("%s_%s.pdf", m.title, m.pageType)
		}
		fWrite, err := os.Create(m.output)
		if err != nil {
			log.Println("creating final PDF file failed", err)
			return
		}

		err = pdfWriter.Write(fWrite)
		if err != nil {
			log.Println("writing PDF failed", err)
		}
		fWrite.Close()
	}
}

func (m *pdfBook) endBook() {
	m.pdf.SetInfo(gopdf.PdfInfo{
		Title:        m.title,
		Author:       `GetNovel用户制作成PDF，并非小说原作者`,
		Creator:      `GetNovel，仅限个人研究学习，对其造成的所有后果，软件作者不负任何责任`,
		Producer:     `GetNovel，仅限个人研究学习，对其造成的所有后果，软件作者不负任何责任`,
		Subject:      m.title + `：不费脑子的适合电子书设备（如Kindle DXG）看的网络小说`,
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
		m.newPage()
	}
	m.pdf.SetFont(m.fontFamily, "", int(m.titleFontSize))
	m.writeTitleLine(articleTitle)
	m.pdf.SetFont(m.fontFamily, "", int(m.contentFontSize))

	c := m.preprocessContent(articleContent)
	lineBreak := "\n"
	for pos := strings.Index(c, lineBreak); ; pos = strings.Index(c, lineBreak) {
		if pos <= 0 {
			if len(c) > 0 {
				m.writeText(c)
			}
			break
		}
		t := c[:pos]
		m.writeText(t)
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
}

func (m *pdfBook) writeTitleLine(t string) {
	m.pdf.Cell(nil, t)
	m.pdf.Br(m.titleFontSize * m.lineSpacing)
	m.height += m.titleFontSize * m.lineSpacing
}

func (m *pdfBook) writeContentLine(t string) {
	m.pdf.Cell(nil, t)
	m.pdf.Br(m.contentFontSize * m.lineSpacing)
	m.height += m.contentFontSize * m.lineSpacing
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
		if width, _ := m.pdf.MeasureTextWidth(t[:count]); width > m.contentWidth {
			if m.height+m.contentFontSize*m.lineSpacing > m.contentHeight {
				m.newPage()
			}
			count -= length
			m.writeContentLine(t[:count])
			t = t[count:]
			index = 0
			count = 0
		} else {
			index += length
		}
	}
	if len(t) > 0 {
		if m.height+m.contentFontSize*m.lineSpacing > m.contentHeight {
			m.newPage()
		}
		m.writeContentLine(t)
	}
}
