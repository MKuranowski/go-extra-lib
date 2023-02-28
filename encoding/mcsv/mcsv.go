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

	"golang.org/x/exp/maps"
)

// Reader reads records from a CSV file.
type Reader struct {
	// Reader.Reader is the [csv.Reader] actually used for parsing the CSV file.
	//
	// Almost all options of the [csv.Reader] are available and can be set
	// before the first call to Read / ReadAll.
	//
	// The two unavailable options are `ReuseRecord` and `FieldsPerRecord`,
	// those are controlled internally by the mcsv.Reader, and their values
	// must not be changed.
	csv.Reader

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

	// lastRecord returned by Read() if ReuseRecord is enabled
	lastRecord map[string]string
}

// NewReader returns a Reader pulling CSV records from r.
//
// The first row is assumed to be the header row.
func NewReader(r io.Reader) *Reader {
	n := &Reader{Reader: *csv.NewReader(r)}
	n.Reader.ReuseRecord = true
	return n
}

// NewReaderWithHeader returns a Reader pulling CSV records from r.
//
// Assumes that r does not contain a header row; instead header is
// used as the column names. All rows in the CSV file must have len(header) fields.
func NewReaderWithHeader(r io.Reader, header []string) *Reader {
	n := &Reader{Reader: *csv.NewReader(r), Header: header}
	n.Reader.ReuseRecord = true
	n.Reader.FieldsPerRecord = len(header)
	return n
}

func (r *Reader) ensureHeader() (err error) {
	if r.Header != nil {
		return nil
	}

	header, err := r.Reader.Read()
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
	recordList, err := r.Reader.Read()
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
