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
	Attribs        []Attribution
	DataPath       string
	SHA256Checksum string
}

func (f *localAttributedFile) Attributions() []Attribution {
	return f.Attribs
}

func (f *localAttributedFile) Read() ([]byte, error) {
	b, err := ioutil.ReadFile(f.DataPath)
	if err != nil {
		return b, err
	}
	s := bytesSHA256(b)
	if s != f.SHA256Checksum {
		return b, fmt.Errorf("Checksum mismatch - file %s - stale=%s current=%s.", f.DataPath, f.SHA256Checksum, s)
	}
	return b, nil
}

func (afp AttributedFilePointer) GetDataPath() (string, error) {
	af, err := afp.ReadAttributedFile()
	if err != nil {
		return "", fmt.Errorf("reading afp: %w", err)
	}
	laf, ok := af.(*localAttributedFile)
	if !ok {
		return "", fmt.Errorf("can't call GetDataPath on type %T", afp)
	}
	return laf.DataPath, nil
}

func (f *localAttributedFile) ReadString() (string, error) {
	b, err := f.Read()
	return string(b), err
}

func (f *localAttributedFile) SHA256() string {
	return f.SHA256Checksum
}

func (f *localAttributedFile) del() error {
	return os.Remove(f.DataPath)
}

type rawAttributedFile struct {
	Attribs []Attribution
	Data    string
}

func (f *rawAttributedFile) Attributions() []Attribution {
	return f.Attribs
}

func (f *rawAttributedFile) Read() ([]byte, error) {
	return []byte(f.Data), nil
}

func (f *rawAttributedFile) ReadString() (string, error) {
	return f.Data, nil
}

func (f *rawAttributedFile) SHA256() string {
	return bytesSHA256([]byte(f.Data))
}

func (f *rawAttributedFile) del() error {
	return nil
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
	panic(fmt.Errorf("Unrecognized attributed file format %s.", filePath))
}

func readAllAttributedFilePointers(dirPath string) ([]AttributedFilePointer, error) {
	var files []AttributedFilePointer
	visit := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("During AFP search at path %s: %v", path, err)
		}
		if strings.HasSuffix(path, attributedFileExtension) {
			files = append(files, AttributedFilePointer{FilePath: path})
		}
		return nil
	}
	err := filepath.Walk(dirPath, visit)
	if err != nil {
		return []AttributedFilePointer{}, err
	}
	return files, nil
}

func readAllAttributedFiles(dirPath string) ([]AttributedFile, error) {
	pointers, err := readAllAttributedFilePointers(dirPath)
	if err != nil {
		return []AttributedFile{}, err
	}
	result := make([]AttributedFile, len(pointers))
	for i, ptr := range pointers {
		af, err := ptr.ReadAttributedFile()
		if err != nil {
			return []AttributedFile{}, fmt.Errorf("While reading raw attributed file: %v", err)
		}
		result[i] = af
	}
	return result, err
}

func readLocalAttributedFile(filePath string) (AttributedFile, error) {
	result := &localAttributedFile{}
	file, err := os.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("While opening local attributed file %s: %v", filePath, err)
	}
	defer file.Close()
	asBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return result, fmt.Errorf("While reading local attributed file %s: %v", filePath, err)
	}
	err = xml.Unmarshal(asBytes, result)
	if err != nil {
		return result, fmt.Errorf("While parsing local attributed file %s: %v", filePath, err)
	}
	return result, nil
}

func readRawAttributedFile(filePath string) (AttributedFile, error) {
	result := &rawAttributedFile{}
	file, err := os.Open(filePath)
	if err != nil {
		return result, fmt.Errorf("While opening raw attributed file %s: %v", filePath, err)
	}
	defer file.Close()
	asBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return result, fmt.Errorf("While reading raw attributed file %s: %v", filePath, err)
	}
	err = xml.Unmarshal(asBytes, result)
	if err != nil {
		return result, fmt.Errorf("While parsing raw attributed file %s: %v", filePath, err)
	}
	return result, nil
}

func attributeLocalFile(filePath string, attributions ...Attribution) (AttributedFilePointer, error) {
	sha256, err := fileSHA256(filePath)
	if err != nil {
		return AttributedFilePointer{}, fmt.Errorf("While computing sha256 for file %s: %v", filePath, err)
	}
	laf := &localAttributedFile{
		Attribs:        attributions,
		DataPath:       filePath,
		SHA256Checksum: sha256,
	}
	aPath := filePath + localAttributedFileExtension
	file, err := xml.MarshalIndent(laf, "", " ")
	if err != nil {
		return AttributedFilePointer{}, fmt.Errorf("While serializing the local attribution for content %s: %v", filePath, err)
	}
	err = ioutil.WriteFile(aPath, file, 0777)
	if err != nil {
		return AttributedFilePointer{}, fmt.Errorf("While writing local attributed file to %v: %v", aPath, err)
	}
	return AttributedFilePointer{FilePath: aPath}, nil
}

func attributeRawFile(filePath, bytes string, attributions ...Attribution) (AttributedFilePointer, error) {
	raf := &rawAttributedFile{
		Attribs: attributions,
		Data:    bytes,
	}
	file, err := xml.MarshalIndent(raf, "", " ")
	if err != nil {
		return AttributedFilePointer{}, fmt.Errorf("While marshalling raw attributed file: %v", err)
	}
	path := filePath + rawAttributedFileExtension
	err = ioutil.WriteFile(path, file, 0777)
	if err != nil {
		return AttributedFilePointer{}, fmt.Errorf("While writing raw attributed file to %v: %v", path, err)
	}
	return AttributedFilePointer{FilePath: path}, nil
}

func (a AttributedFilePointer) del() error {
	af, err := a.ReadAttributedFile()
	if err != nil {
		return fmt.Errorf("While reading attributed file at %s for deletion: %v", a.FilePath, err)
	}
	err = af.del()
	if err != nil {
		return fmt.Errorf("While deleting attributed file at %s: %v", a.FilePath, err)
	}
	return os.Remove(a.FilePath)
}

func (a AttributedFilePointer) copyTo(newPath string) (AttributedFilePointer, error) {
	af, err := a.ReadAttributedFile()
	if err != nil {
		return a, fmt.Errorf("While reading attributed file at %s: %v", a.FilePath, err)
	}
	bytes, err := af.Read()
	if err != nil {
		return a, fmt.Errorf("While reading bytes from file referenced by %s: %v", a.FilePath, err)
	}
	err = os.MkdirAll(filepath.Dir(newPath), 0777)
	if err != nil {
		return a, fmt.Errorf("While creating directories to %s: %v", newPath, err)
	}
	err = ioutil.WriteFile(newPath, bytes, 0777)
	if err != nil {
		return a, fmt.Errorf("While writing copy of file to %s: %v", newPath, err)
	}
	afp, err := AttributeLocalFile(newPath, af.Attributions()...)
	if err != nil {
		return a, fmt.Errorf("While attributing newly copied file at %s: %v", newPath, err)
	}
	return afp, nil
}

func (a AttributedFilePointer) base() string {
	s := filepath.Base(a.FilePath)
	s = strings.ReplaceAll(s, rawAttributedFileExtension, "")
	s = strings.ReplaceAll(s, localAttributedFileExtension, "")
	return s
}
