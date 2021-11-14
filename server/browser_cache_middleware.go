package server

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"hash"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/RobertGrantEllis/t9/logger"
)

var (
	lastModified          time.Time
	formattedLastModified string
)

func init() {
	// Apparently fs.FS does not capture modified dates, so capture one
	// as soon as the executable starts.
	gmt, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
	lastModified = time.Now().In(gmt)
	formattedLastModified = lastModified.Format(time.RFC1123)
}

type browserCacheRecords struct {
	records map[string]*pathRecord
	mutex   *sync.RWMutex
}

func (bcr *browserCacheRecords) getRecord(path string) *pathRecord {
	bcr.mutex.RLock()
	defer bcr.mutex.RUnlock()
	return bcr.records[path]
}

func (bcr *browserCacheRecords) saveRecord(path string, pr *pathRecord) {
	bcr.mutex.Lock()
	defer bcr.mutex.Unlock()
	bcr.records[path] = pr
}

type pathRecord struct {
	etagSum []byte
	size    int

	// TODO: we need separate etags for each header designated in the "vary" header.
	// Example: "Vary: Content-Encoding"-- If "Content-Encoding" header is present
	// in the response, then we must capture the header value and include it in the etag
	// reported to the browser (i.e. with a suffix), and/or calculate different hashes
	// for different responses. IT IS IMPORTANT THAT THIS MIDDLEWARE IS THE OUTERMOST that
	// can mutate responses (i.e. OUTSIDE compresssion middleware).

	formattedETAG string
}

type pathRecordResponseWriter struct {
	http.ResponseWriter
	hasher hash.Hash

	statusCode int
	size       int
	err        error
}

func (prrw *pathRecordResponseWriter) WriteHeader(statusCode int) {
	if prrw.statusCode != 0 {
		return // status code already written
	}
	prrw.statusCode = statusCode
	if statusCode == http.StatusOK {
		prrw.ResponseWriter.Header().Set("Last-Modified", formattedLastModified)
		// TODO: really all other relevant headers (cache-control, expires, etag) should be written here
	}
	prrw.ResponseWriter.WriteHeader(statusCode)
}

func (prrw *pathRecordResponseWriter) Write(buffer []byte) (int, error) {
	prrw.WriteHeader(http.StatusOK) // will do nothing if status code already set

	writtenToResponseWriter, err := prrw.ResponseWriter.Write(buffer)
	if err != nil {
		// error writing to the client means that the hash will be inaccurate so note as much
		prrw.err = err
		return writtenToResponseWriter, err
	}
	if prrw.err != nil || prrw.hasher == nil {
		// error writing to hasher in prior invocation, so just return
		return writtenToResponseWriter, err
	}

	// no error so far
	writtenToHasher, err := prrw.hasher.Write(buffer)
	if err != nil {
		prrw.err = err
	} else if writtenToResponseWriter != writtenToHasher {
		prrw.err = fmt.Errorf(
			"different number of bytes written to response (%d) vs hasher (%d)",
			writtenToResponseWriter,
			writtenToHasher,
		)
	} else {
		prrw.size += writtenToHasher
	}

	return writtenToResponseWriter, nil
}

func enableBrowserCacheMiddleware(inner http.Handler, l logger.Logger) http.Handler {
	bcr := &browserCacheRecords{
		records: map[string]*pathRecord{},
		mutex:   &sync.RWMutex{},
	}

	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		path := req.URL.Path

		// may return nil if there is no record yet
		record := bcr.getRecord(path)

		prrw := &pathRecordResponseWriter{
			ResponseWriter: rw,
			// for now assume we do not need to capture a hash
		}

		if record != nil {
			// no need to calculate ETAG since we already have it
			// TODO: We need different ETAG RECORDS FOR EACH CONTENT_ENCODING!!
			rw.Header().Set("ETAG", record.formattedETAG)
			//TODO: return appropriate headers and responses depending on client headers
			inner.ServeHTTP(prrw, req)
			return
		}

		prrw.hasher = sha1.New()

		inner.ServeHTTP(prrw, req)

		if prrw.err != nil {
			l.Error(prrw.err)
			return
		}

		if prrw.statusCode == http.StatusOK {
			etagSum := prrw.hasher.Sum(nil)
			bcr.saveRecord(path, &pathRecord{
				etagSum: etagSum,
				size:    prrw.size,
				formattedETAG: fmt.Sprintf(
					"%v-%v",
					strconv.Itoa(prrw.size),
					hex.EncodeToString(etagSum),
				),
			})
		}
	})
}
