package attributions

import "time"

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
}

func ReadAttributedFile(filePath string) (AttributedFile, error) {
	return readAttributedFile(filePath)
}

func ReadAllAttributedFiles(dirPath string) ([]AttributedFile, error) {
	return readAllAttributedFiles(dirPath)
}

func AttributeLocalFile(filePath string, attributions ...Attribution) error {
	return attributeLocalFile(filePath, attributions...)
}

func AttributeRawFile(filePath, bytes string, attributions ...Attribution) error {
	return attributeRawFile(filePath, bytes, attributions...)
}
