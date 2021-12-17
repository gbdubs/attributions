package attributions

import (
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"sync"
)

var memoizedSHA256sLock = sync.RWMutex{}
var memoizedSHA256s = make(map[string]string)

func fileSHA256(filePath string) (string, error) {
	memoizedSHA256sLock.RLock()
	r := memoizedSHA256s[filePath]
	memoizedSHA256sLock.RUnlock()
	if r != "" {
		return r, nil
	}
	f, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	r = fmt.Sprintf("%x", h.Sum(nil))
	memoizedSHA256sLock.Lock()
	memoizedSHA256s[filePath] = r
	memoizedSHA256sLock.Unlock()
	return r, nil
}

func bytesSHA256(b []byte) string {
	h := sha256.New()
	h.Write(b)
	return fmt.Sprintf("%x", h.Sum(nil))
}
