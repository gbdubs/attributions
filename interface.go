package attributions

import (
	"time"
)

type Attribution struct {
	OriginUrl           string    `xml:"origin_url"`
	CollectedAt         time.Time `xml:"collected_at"`
	OriginalTitle       string    `xml:"original_title"`
	Author              string    `xml:"author"`
	AuthorUrl           string    `xml:"author_url"`
	License             string    `xml:"license"`
	LicenseUrl          string    `xml:"license_url"`
	CreatedAt           time.Time `xml:"created_at"`
	Context             []string  `xml:"context"`
	ScrapingMethodology string    `xml:"scraping_methodology"`
}

type AttributedFile interface {
	Attributions() []Attribution
	Read() ([]byte, error)
	ReadString() (string, error)
	SHA256() string
	del() error
}

type AttributedFilePointer struct {
	FilePath string
}

func ReadAttributedFile(filePath string) (AttributedFile, error) {
	return readAttributedFile(filePath)
}

func (a AttributedFilePointer) ReadAttributedFile() (AttributedFile, error) {
	return readAttributedFile(a.FilePath)
}

func ReadAllAttributedFilePointers(dirPath string) ([]AttributedFilePointer, error) {
	return readAllAttributedFilePointers(dirPath)
}

func ReadAllAttributedFiles(dirPath string) ([]AttributedFile, error) {
	return readAllAttributedFiles(dirPath)
}

func AttributeLocalFile(filePath string, attributions ...Attribution) (AttributedFilePointer, error) {
	return attributeLocalFile(filePath, attributions...)
}

func AttributeRawFile(filePath, bytes string, attributions ...Attribution) (AttributedFilePointer, error) {
	return attributeRawFile(filePath, bytes, attributions...)
}

func (a AttributedFilePointer) Delete() error {
	return a.del()
}

func (a AttributedFilePointer) CopyTo(newLocalPath string) (AttributedFilePointer, error) {
	return a.copyTo(newLocalPath)
}

func (a AttributedFilePointer) CopyToDir(newDirPath string) (AttributedFilePointer, error) {
	return a.copyTo(newDirPath + "/" + a.base())
}

func (a AttributedFilePointer) Base() string {
	return a.base()
}
