// based on http://golang.org/src/io/io.go?s=12247:12307#L340
package testi

import (
	"encoding/hex"
	"hash"
	"io"
	// this is broken:
	//"golang.org/x/crypto/sha3"
	"github.com/tonnerre/golang-go.crypto/sha3"
)

func HashCopy(dst io.Writer, src io.Reader, digest hash.Hash) (written int64, err error) {
	buf := make([]byte, 32*1024)
	for {
		nr, er := src.Read(buf)
		if nr > 0 {
			nw, ew := dst.Write(buf[0:nr])
			if nw > 0 {
				written += int64(nw)
				nw_digest, ew_digest := digest.Write(buf[0:nw])
				if nw != nw_digest {
					err = io.ErrShortWrite
					break
				}
				if ew_digest != nil {
					err = ew_digest
					break
				}
			}
			if ew != nil {
				err = ew
				break
			}
			if nr != nw {
				err = io.ErrShortWrite
				break
			}
		}
		if er == io.EOF {
			break
		}
		if er != nil {
			err = er
			break
		}
	}
	return written, err
}

func Sha3HashCopy(dst io.Writer, src io.Reader) (written int64, sha3sum string, err error) {
	h := sha3.NewKeccak256()
	written, err = HashCopy(dst, src, h)
	if err != nil {
		return 0, "", err
	}
	sha3sum = hex.EncodeToString(h.Sum(nil))
	return written, sha3sum, err

}
