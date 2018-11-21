package handler

import (
	"errors"
	"fmt"
	"github.com/syyongx/llog/types"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"
)

// Stores logs to files that are rotated every day and a limited number of files are kept.
//
// This rotation is only intended to be used as a workaround. Using logrotate to
// handle the rotation is strongly encouraged when you can use it.
type RotatingFile struct {
	*File

	filename       string
	maxFiles       int
	mustRotate     bool
	nextRotation   int
	filenameFormat string
	dateFormat     string
	perm           os.FileMode
	sync.Mutex
}

// level: The minimum logging level at which this handler will be triggered
// bubble: Whether the messages that are handled can bubble up the stack or not
// filePerm: Optional file permissions (default (0644) are only for owner read/write)
func NewRotatingFile(filename string, maxFiles, level int, bubble bool, filePerm os.FileMode) *RotatingFile {
	rf := &RotatingFile{
		filename:       filename,
		maxFiles:       maxFiles,
		filenameFormat: "{filename}-{date}",
		dateFormat:     "2006-01-02",
		nextRotation:   time.Now().AddDate(0, 0, 1).Day(), // tomorrow
		perm:           filePerm,
	}
	path := rf.timedFilename()
	rf.File = NewFile(path, level, bubble, filePerm)
	rf.File.Writer = rf.Write
	return rf
}

// Set filename format.
func (rf *RotatingFile) SetFilenameFormat(filenameFormat, dateFormat string) error {
	match, _ := regexp.MatchString("^2006(([/_.-]?01)([/_.-]?02)?)?$", dateFormat)
	if !match {
		return errors.New("invalid date format")
	}
	if n := strings.Index(filenameFormat, "{date}"); n < 0 {
		return errors.New("invalid filename format, format should contain at least {date}")
	}

	rf.filenameFormat = filenameFormat
	rf.dateFormat = dateFormat
	rf.Close()

	return nil
}

// Write to file.
func (rf *RotatingFile) Write(record *types.Record) {
	if rf.nextRotation < record.Datetime.Day() {
		rf.mustRotate = true
		rf.Close()
	}

	rf.File.Write(record)
}

// Closes the handler.
func (rf *RotatingFile) Close() {
	rf.File.Close()

	if rf.mustRotate {
		rf.rotate()
	}
}

// Rotates the files.
func (rf *RotatingFile) rotate() error {
	// update path
	rf.Path = rf.timedFilename()
	rf.Fd = nil
	// tomorrow
	rf.nextRotation = time.Now().AddDate(0, 0, 1).Day()
	// skip GC of old logs if files are unlimited
	if rf.maxFiles == 0 {
		return nil
	}
	files, err := filepath.Glob(rf.globPattern())
	if err != nil {
		return err
	}
	if len(files) > rf.maxFiles {
		// no files to remove
		return nil
	}
	// Sorting the files by name to remove the older ones
	sort.Strings(files)
	for _, file := range files[:rf.maxFiles] {
		err := os.Remove(file)
		if err != nil {
			return err
		}
	}

	rf.mustRotate = false
	return nil
}

// Get timed filename
func (rf *RotatingFile) timedFilename() string {
	dir := filepath.Dir(rf.filename)
	basename := filepath.Base(rf.filename)
	ext := filepath.Ext(rf.filename)
	if ext != "" {
		basename = basename[:strings.Index(basename, ext)]
	}

	date := time.Unix(time.Now().Unix(), 0).Format(rf.dateFormat)
	timedFilename := strings.NewReplacer("{filename}", basename, "{date}", date).Replace(dir + "/" + rf.filenameFormat)
	timedFilename += ext

	return timedFilename
}

// Get blob pattern
func (rf *RotatingFile) globPattern() string {
	dir := filepath.Dir(rf.filename)
	basename := filepath.Base(rf.filename)
	ext := filepath.Ext(rf.filename)

	glob := strings.NewReplacer("{filename}", basename, "{date}", "*").Replace(dir + "/" + rf.filenameFormat)
	glob += ext
	fmt.Println(glob)

	return glob
}
