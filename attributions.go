package attributions

import (
	"encoding/xml"
	"io/ioutil"
	"os"
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

func (a *Attribution) WriteTo(filePathNoSuffix string) error {
	file, err := xml.MarshalIndent(*a, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePathNoSuffix+".atrib", file, 0644)
}

func ReadFrom(filePathNoSuffix string) (*Attribution, error) {
	var attribution Attribution
	file, err := os.Open(filePath(filePathNoSuffix))
	if err != nil {
		return &attribution, err
	}
	defer file.Close()
	asBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return &attribution, err
	}
	err = xml.Unmarshal(asBytes, &attribution)
	return &attribution, err
}

func filePath(withoutSuffix string) string {
	return withoutSuffix + ".attrib"
}
