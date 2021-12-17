package attributions

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type localAttributedFile struct {
	attributions   []Attribution
	dataPath       string
	sha256Checksum string
}

func (f *localAttributedFile) Attributions() []Attribution {
	return f.attributions
}

func (f *localAttributedFile) Read() ([]byte, error) {
	b, err := ioutil.ReadFile(f.dataPath)
	if err != nil {
		return b, err
	}
	s := bytesSHA256(b)
	if s != f.sha256Checksum {
		return b, fmt.Errorf("Checksum mismatch - attribution may be stale %s versus %s.", s, f.sha256Checksum)
	}
	return b, nil
}

func (f *localAttributedFile) ReadString() (string, error) {
	b, err := f.Read()
	return string(b), err
}

type rawAttributedFile struct {
	attributions []Attribution
	data         string
}

func (f *rawAttributedFile) Attributions() []Attribution {
	return f.attributions
}

func (f *rawAttributedFile) Read() ([]byte, error) {
	return []byte(f.data), nil
}

func (f *rawAttributedFile) ReadString() (string, error) {
	return f.data, nil
}

const attributedFileExtension = ".attrib"
const rawAttributedFileExtension = ".raw" + attributedFileExtension
const localAttributedFileExtension = ".local" + attributedFileExtension

func readAttributedFile(filePath string) (AttributedFile, error) {
	if strings.HasSuffix(filePath, rawAttributedFileExtension) {
		return readRawAttributedFile(filePath)
	} else if strings.HasSuffix(filePath, localAttributedFileExtension) {
		return readLocalAttributedFile(filePath)
	}
	return nil, fmt.Errorf("Unrecognized attributed file format %s.", filePath)
}

func readAllAttributedFiles(dirPath string) ([]AttributedFile, error) {
	var files []string
	visit := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("Recursive search error at path %s: %v", path, err)
		}
		if strings.HasSuffix(path, attributedFileExtension) {
			files = append(files, path)
		}
		return nil
	}
	err := filepath.Walk(dirPath, visit)
	empty := []AttributedFile{}
	if err != nil {
		return empty, err
	}
	result := make([]AttributedFile, len(files))
	for i, file := range files {
		af, err := ReadAttributedFile(file)
		if err != nil {
			return empty, fmt.Errorf("Error reading raw attributed file: %v", err)
		}
		result[i] = af
	}
	return result, err
}

func readLocalAttributedFile(filePath string) (AttributedFile, error) {
	result := &localAttributedFile{}
	file, err := os.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()
	asBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return result, fmt.Errorf("Error reading file %s: %v", filePath, err)
	}
	err = xml.Unmarshal(asBytes, result)
	if err != nil {
		return result, fmt.Errorf("Error parsing file %s: %v", filePath, err)
	}
	return result, nil
}

func readRawAttributedFile(filePath string) (AttributedFile, error) {
	result := &rawAttributedFile{}
	file, err := os.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("Error opening file %s: %v", filePath, err)
	}
	defer file.Close()
	asBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return result, fmt.Errorf("Error reading file %s: %v", filePath, err)
	}
	err = xml.Unmarshal(asBytes, result)
	if err != nil {
		return result, fmt.Errorf("Error parsing file %s: %v", filePath, err)
	}
	return result, nil
}

func attributeLocalFile(filePath string, attributions ...Attribution) error {
	sha256, err := fileSHA256(filePath)
	if err != nil {
		return fmt.Errorf("Error computing sha256 for file %s: %v", filePath, err)
	}
	laf := &localAttributedFile{
		attributions:   attributions,
		dataPath:       filePath,
		sha256Checksum: sha256,
	}
	aPath := filePath + localAttributedFileExtension
	file, err := xml.MarshalIndent(laf, "", " ")
	if err != nil {
		return fmt.Errorf("Error serializing the local attribution for content %s: %v", filePath, err)
	}
	return ioutil.WriteFile(aPath, file, 0777)
}

func attributeRawFile(filePath, bytes string, attributions ...Attribution) error {
	raf := &rawAttributedFile{
		attributions: attributions,
		data:         bytes,
	}
	file, err := xml.MarshalIndent(raf, "", " ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filePath+rawAttributedFileExtension, file, 0777)
}
