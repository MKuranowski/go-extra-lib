// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

// mcsv is a wrapper around the built-in [encoding/csv] package,
// parsing [RFC 4180] CSV files into map[string]string rows.
//
// [RFC 4180]: https://rfc-editor.org/rfc/rfc4180.html
package mcsv

import (
	"encoding/csv"
	"errors"
	"io"
	"unicode/utf8"

	"golang.org/x/exp/maps"
)

// Reader reads records from a CSV io.Reader.
type Reader struct {
	// Reader.Reader is the [csv.Reader] actually used for parsing the CSV file.
	//
	// Almost all options of the [csv.Reader] are available and can be set
	// before the first call to Read / ReadAll.
	//
	// The two unavailable options are `ReuseRecord` and `FieldsPerRecord`,
	// those are controlled internally by the mcsv.Reader, and their values
	// must not be changed.
	*csv.Reader

	// Header is a slice of column names to be used as keys in returned records.
	//
	// Set by the first call to Read or NewReaderWithHeader, must not be modified.
	Header []string

	// ReuseRecord shadows the csv.Reader.ReuseRecord setting.
	//
	// If ReuseRecord is set to true, the same map will be returned by all calls to Read().
	//
	// Ignored by the ReadAll() function.
	//
	// The csv.Reader.ReuseRecord setting is controlled by mcsv.Reader and must not be changed.
	ReuseRecord bool

	// PreserveBOM ensures an initial byte-order-mark in the first ever read from the file
	// is not removed. The default behavior removes the BOM.
	PreserveBOM bool

	// lastRecord returned by Read() if ReuseRecord is enabled
	lastRecord map[string]string

	// didRead is flag used to control the behavior of BOM removal.
	removedBOM bool
}

// NewReader returns a Reader pulling CSV records from r.
//
// The first row is assumed to be the header row.
func NewReader(r io.Reader) *Reader {
	n := &Reader{Reader: csv.NewReader(r)}
	n.Reader.ReuseRecord = true
	return n
}

// NewReaderWithHeader returns a Reader pulling CSV records from r.
//
// Assumes that r does not contain a header row; instead header is
// used as the column names. All rows in the CSV file must have len(header) fields.
func NewReaderWithHeader(r io.Reader, header []string) *Reader {
	n := &Reader{Reader: csv.NewReader(r), Header: header}
	n.Reader.ReuseRecord = true
	n.Reader.FieldsPerRecord = len(header)
	return n
}

// readRow returns the result of calling r.Reader.Read,
// with additionally handling the byte-order-mark.
func (r *Reader) readRow() (row []string, err error) {
	row, err = r.Reader.Read()

	// Remove the byte-order-mark
	if err == nil && !r.PreserveBOM && !r.removedBOM && len(row) > 0 {
		r.removedBOM = true
		first, size := utf8.DecodeRuneInString(row[0])
		if first == '\uFEFF' {
			row[0] = row[0][size:]
		}
	}

	return
}

func (r *Reader) ensureHeader() (err error) {
	if r.Header != nil {
		return nil
	}

	header, err := r.readRow()
	if err != nil {
		return
	} else {
		// r.Reader.FieldsPerRecord set by csv.Reader.Read()
		r.Header = make([]string, len(header))
		copy(r.Header, header)
	}
	return nil
}

// Read reads a record from the CSV file and returns it.
// Always returns either a non-nil record or a non-nil err, but never both.
// The exception from [csv.Reader.Read] does not apply.
// If there are no more records to read, returns (nil, io.EOF).
//
// If ReuseRecord is set, this function may return the same map as a previous call to Read(),
// just with the values changed.
func (r *Reader) Read() (record map[string]string, err error) {
	// Ensure the header was parsed
	err = r.ensureHeader()
	if err != nil {
		return
	}

	// retrieve the next record
	recordList, err := r.readRow()
	if err == nil {
		// prepare the record map
		if r.ReuseRecord && r.lastRecord != nil {
			record = r.lastRecord
		} else {
			record = make(map[string]string, len(r.Header))
		}

		// put data from the record into a map
		for i, colName := range r.Header {
			record[colName] = recordList[i]
		}

		// cache the record if ReuseRecord is true
		if r.ReuseRecord {
			r.lastRecord = record
		}
	}

	return
}

// ReadAll repeatedly calls Read to read all remaining records from a file.
// If an error occurs, returns all records read up to the error and the error itself.
// Successful ReadAll call returns (records, nil), as io.EOF is not deemed an error in this context.
//
// ReuseRecord setting is automatically set to false.
func (r *Reader) ReadAll() (records []map[string]string, err error) {
	err = r.ensureHeader()
	if err != nil {
		return
	}

	// Ensure mcsv.Reader.Read() returns newly-allocated maps.
	r.ReuseRecord = false

	for {
		var record map[string]string
		record, err = r.Read()

		if errors.Is(err, io.EOF) {
			err = nil
			break
		} else if err != nil {
			break
		}

		records = append(records, maps.Clone(record))
	}

	return
}

// Writer writes CSV records into a io.Writer.
//
// Writes to the underlying file are buffered and client must call Flush()
// to ensure data was actually written to the io.Writer.
// Any encountered errors may be checked with the Error() method.
type Writer struct {
	// Writer.Writer is the [csv.Writer] actually used for encoding and writing the CSV data.
	//
	// All options of the [csv.Writer] are available and can be set
	// before the first call to Read / ReadAll.
	*csv.Writer

	// Header is a slice of column names to be used as keys in returned records.
	//
	// Set by NewWriter, must not be modified after construction.
	Header []string

	// row is a pre-allocated slice used for calling csv.Writer.Write.
	row []string
}

// NewWriter returns a *Writer for encoding records as CSV and
// writing them to an underlying io.Writer.
//
// The header row is not written to the file by default - use the WriteHeader().
func NewWriter(w io.Writer, header []string) *Writer {
	return &Writer{
		Writer: csv.NewWriter(w),
		Header: header,
		row:    make([]string, len(header)),
	}
}

// WriteHeader writes the header row.
//
// All writes to the underlying io.Writer are buffered and some data
// may not be actually written unless Flush() is called.
func (w *Writer) WriteHeader() error {
	return w.Writer.Write(w.Header)
}

// Write writes a record to the CSV file.
//
// Any missing fields are replaced with an empty string,
// and extra fields are ignored.
//
// All writes to the underlying io.Writer are buffered and some data
// may not be actually written unless Flush() is called.
func (w *Writer) Write(record map[string]string) error {
	for i, column := range w.Header {
		w.row[i] = record[column]
	}
	return w.Writer.Write(w.row)
}

// WriteAll calls Write for every provided record, and then calls Flush().
func (w *Writer) WriteAll(records []map[string]string) error {
	for _, record := range records {
		err := w.Write(record)
		if err != nil {
			return err
		}
	}
	w.Flush()
	return w.Error()
}
