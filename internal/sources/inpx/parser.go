package inpx

import (
	"archive/zip"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/HoskeOwl/PoorBookExtractor/internal/entities"
	"github.com/HoskeOwl/PoorBookExtractor/internal/logs"
	"github.com/HoskeOwl/PoorBookExtractor/internal/sources/inp"
	"go.uber.org/zap"
)

/*
inpx is a Zip archive, inside of which are text files 'inp'
nd of record - <CR><LF> - <0x0D,0x0A>
*/

type InpxParser struct {
	filename string
	path     string
	progress *atomic.Int32
}

func NewInpxParser(filename string) *InpxParser {
	return &InpxParser{progress: &atomic.Int32{}}
}

func (imp *InpxParser) readZipFile(zf *zip.File) ([]byte, error) {
	file, err := zf.Open()
	if err != nil {
		return nil, err
	}
	defer file.Close()
	buf := bytes.NewBuffer(make([]byte, 0, zf.UncompressedSize64))
	io.Copy(buf, file)
	return buf.Bytes(), nil
}

func (imp *InpxParser) readZip() (*zip.Reader, func() error, error) {
	file, err := os.Open(filepath.Join(imp.path, imp.filename))
	if err != nil {
		return nil, nil, err
	}
	stat, err := file.Stat()
	if err != nil {
		return nil, nil, err
	}
	zipReader, err := zip.NewReader(file, stat.Size())
	if err != nil {
		return nil, nil, err
	}
	return zipReader, file.Close, nil
}

func (imp *InpxParser) Parse(ctx context.Context) ([]entities.Book, error) {
	log := logs.GetFromContext(ctx).With(zap.String("action", "parse_inpx"))
	zipReader, closer, err := imp.readZip()
	if err != nil {
		return nil, err
	}
	defer closer()
	books := make([]entities.Book, 0)
	// Read all the files from zip archive
	for _, zipFile := range zipReader.File {
		if !strings.HasSuffix(zipFile.Name, ".inp") {
			log.Debug("skipping file", zap.String("filename", zipFile.Name))
			continue
		}
		log.Debug("Reading file:", zap.String("filename", zipFile.Name))
		unzippedFileBytes, err := imp.readZipFile(zipFile)
		if err != nil {
			return nil, err
		}
		books, err = inp.ParseBooks(ctx, unzippedFileBytes)
		if err != nil {
			return nil, err
		}
		for _, book := range books {
			metadata := entities.BookMetadata{
				ArchiveName: zipFile.Name,
				Filepath:    filepath.Join(imp.path),
			}
			book.Metadata = metadata
		}
	}
	return books, nil
}

func getProgress(total int, current int) int32 {
	return int32(float64(current) / float64(total) * 100)
}

func (imp *InpxParser) ParseSkipErrors(ctx context.Context) ([]entities.Book, error) {
	imp.progress.Store(0)
	defer imp.progress.Store(100)

	log := logs.GetFromContext(ctx).With(zap.String("action", "parse_inpx_skip_errors"))
	zipReader, closer, err := imp.readZip()
	if err != nil {
		return nil, err
	}
	defer closer()
	// Read all the files from zip archive
	books := make([]entities.Book, 0)
	total := len(zipReader.File)
	for idx, zipFile := range zipReader.File {
		if !strings.HasSuffix(zipFile.Name, ".inp") {
			log.Debug("skipping file", zap.String("filename", zipFile.Name))
			continue
		}
		log.Debug("reading file", zap.String("filename", zipFile.Name))
		unzippedFileBytes, err := imp.readZipFile(zipFile)
		if err != nil {
			log.Error("error ")
			continue
		}
		metadata := entities.BookMetadata{
			ArchiveName: zipFile.Name[:len(zipFile.Name)-3] + "zip",
			Filepath:    filepath.Join(imp.path),
		}
		books, err = inp.ParseBooksWithMetadata(ctx, unzippedFileBytes, metadata)
		if err != nil {
			return nil, err
		}
		imp.progress.Store(getProgress(total, idx+1))
	}
	return books, nil
}

func (imp *InpxParser) GetProgress() int {
	return int(imp.progress.Load())
}

func (imp *InpxParser) ParseBooks(ctx context.Context, path string) ([]entities.Book, error) {

	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	dir := filepath.Dir(path)
	filename := filepath.Base(path)

	imp.path = dir
	imp.filename = filename

	books, err := imp.ParseSkipErrors(ctx)
	if err != nil {
		return nil, err
	}
	return books, nil
}
