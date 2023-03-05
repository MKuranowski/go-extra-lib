// Copyright (c) 2023 Mikołaj Kuranowski
// SPDX-License-Identifier: MIT

package mcsv_test

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/MKuranowski/go-extra-lib/encoding/mcsv"
	"github.com/MKuranowski/go-extra-lib/io2"
	"github.com/MKuranowski/go-extra-lib/iter"
	"github.com/MKuranowski/go-extra-lib/testing2/assert"
	"github.com/MKuranowski/go-extra-lib/testing2/check"
)

const (
	dataNewLines = "field1,field2,field3\r\n" +
		"\"hello\",\"is it \"\"me\"\"\",\"you're\n" +
		"looking for\"\r\n" +
		"this is going to be,\"another\n" +
		"broken row\",\"very confusing\"\r\n"
)

type readerTest struct {
	name   string
	input  string
	header []string // If nil, use mcsv.NewReader; otherwise use mcsv.NewReaderWithHeader
	result []map[string]string

	comma       rune // if non zero, set Reader.Comma
	preserveBOM bool
}

var readerTests = []readerTest{
	{
		name: "EuCities",
		input: `"City","Country"
"Berlin","Germany"
"Madrid","Spain"
"Rome","Italy"
"Bucharest","Romania"
"Paris","France"
`,
		header: nil,
		result: []map[string]string{
			{"City": "Berlin", "Country": "Germany"},
			{"City": "Madrid", "Country": "Spain"},
			{"City": "Rome", "Country": "Italy"},
			{"City": "Bucharest", "Country": "Romania"},
			{"City": "Paris", "Country": "France"},
		},
	},

	{
		name: "MathConstants",
		input: `pi,3.1416
sqrt2,1.4142
phi,1.618
e,2.7183
`,
		header: []string{"name", "value"},
		result: []map[string]string{
			{"name": "pi", "value": "3.1416"},
			{"name": "sqrt2", "value": "1.4142"},
			{"name": "phi", "value": "1.618"},
			{"name": "e", "value": "2.7183"},
		},
	},

	{
		name: "MetroSystems",
		input: `"City"	"Stations"	"System Length"
"New York"	"424"	"380"
"Shanghai"	"345"	"676"
"Seoul"	"331"	"353"
"Beijing"	"326"	"690"
"Paris"	"302"	"214"
"London"	"270"	"402"
`,
		result: []map[string]string{
			{"City": "New York", "Stations": "424", "System Length": "380"},
			{"City": "Shanghai", "Stations": "345", "System Length": "676"},
			{"City": "Seoul", "Stations": "331", "System Length": "353"},
			{"City": "Beijing", "Stations": "326", "System Length": "690"},
			{"City": "Paris", "Stations": "302", "System Length": "214"},
			{"City": "London", "Stations": "270", "System Length": "402"},
		},
		comma: '\t',
	},

	{
		name:  "Newlines",
		input: dataNewLines,
		result: []map[string]string{
			{"field1": "hello", "field2": "is it \"me\"", "field3": "you're\nlooking for"},
			{"field1": "this is going to be", "field2": "another\nbroken row", "field3": "very confusing"},
		},
	},

	{
		name:   "SkipsBOM",
		input:  "\uFEFFa,b,c\n1,2,3\n",
		result: []map[string]string{{"a": "1", "b": "2", "c": "3"}},
	},

	{
		name:   "SkipsBOMWithHeader",
		input:  "\uFEFF1,2,3\n",
		header: []string{"a", "b", "c"},
		result: []map[string]string{{"a": "1", "b": "2", "c": "3"}},
	},

	{
		name:        "PreserveBOM",
		input:       "\uFEFFa,b,c\n1,2,3\n",
		result:      []map[string]string{{"\uFEFFa": "1", "b": "2", "c": "3"}},
		preserveBOM: true,
	},

	{
		name:        "PreserveBOMWithHeader",
		input:       "\uFEFF1,2,3\n",
		header:      []string{"a", "b", "c"},
		result:      []map[string]string{{"a": "\uFEFF1", "b": "2", "c": "3"}},
		preserveBOM: true,
	},
}

func runReadTest(t *testing.T, r *mcsv.Reader, expected []map[string]string) {
	it := iter.ZipLongest(nil, iter.OverIOReader[map[string]string](r), iter.OverSlice(expected))
	i := 0
	for it.Next() {
		elem := it.Get()
		check.DeepEqMsg(t, elem[0], elem[1], fmt.Sprintf("row %d", i))
		i++
	}
	assert.NoErr(t, it.Err())
}

func runReadAllTest(t *testing.T, r *mcsv.Reader, expected []map[string]string) {
	got, err := r.ReadAll()
	assert.NoErr(t, err)

	it := iter.ZipLongest(nil, iter.OverSlice(got), iter.OverSlice(expected))
	i := 0
	for it.Next() {
		elem := it.Get()
		check.DeepEqMsg(t, elem[0], elem[1], fmt.Sprintf("row %d", i))
		i++
	}
	assert.NoErr(t, it.Err())
}

func getReaderForTest(test readerTest) (r *mcsv.Reader) {
	inputReader := strings.NewReader(test.input)
	if test.header != nil {
		r = mcsv.NewReaderWithHeader(inputReader, test.header)
	} else {
		r = mcsv.NewReader(inputReader)
	}

	if test.comma != 0 {
		r.Comma = test.comma
	}
	r.PreserveBOM = test.preserveBOM

	return
}

func ExampleReader() {
	in := `number,value
pi,3.1416
sqrt2,1.4142
phi,1.618
e,2.7183
`

	r := mcsv.NewReader(strings.NewReader(in))

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Number", record["number"], "has value", record["value"])
	}

	// Output:
	// Number pi has value 3.1416
	// Number sqrt2 has value 1.4142
	// Number phi has value 1.618
	// Number e has value 2.7183
}

func ExampleReader_missing_header() {
	in := `pi,3.1416
sqrt2,1.4142
phi,1.618
e,2.7183
`

	r := mcsv.NewReaderWithHeader(strings.NewReader(in), []string{"number", "value"})

	for {
		record, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Number", record["number"], "has value", record["value"])
	}

	// Output:
	// Number pi has value 3.1416
	// Number sqrt2 has value 1.4142
	// Number phi has value 1.618
	// Number e has value 2.7183
}

func TestReaderRead(t *testing.T) {
	for _, test := range readerTests {
		r := getReaderForTest(test)
		t.Run(test.name, func(t *testing.T) { runReadTest(t, r, test.result) })
	}
}

func TestReaderReadAll(t *testing.T) {
	for _, test := range readerTests {
		r := getReaderForTest(test)
		t.Run(test.name, func(t *testing.T) { runReadAllTest(t, r, test.result) })
	}
}

func TestReaderReadReuseRecord(t *testing.T) {
	for _, test := range readerTests {
		r := getReaderForTest(test)
		r.ReuseRecord = true
		t.Run(test.name, func(t *testing.T) { runReadTest(t, r, test.result) })
	}
}

func TestReaderReadAllReuseRecord(t *testing.T) {
	for _, test := range readerTests {
		r := getReaderForTest(test)
		r.ReuseRecord = true
		t.Run(test.name, func(t *testing.T) { runReadAllTest(t, r, test.result) })
	}
}

const readerBenchmarkData = `f1,f2,f3,f4
a,b,c,d
w,x,y,z
i,j,k,l
1,2,3,4
`

func runReadBenchmark(b *testing.B, header []string, initFunc func(*mcsv.Reader)) {
	b.ReportAllocs()

	var r *mcsv.Reader
	in := io2.Repeated(readerBenchmarkData, b.N)
	if header != nil {
		r = mcsv.NewReaderWithHeader(in, header)
	} else {
		r = mcsv.NewReader(in)
	}
	if initFunc != nil {
		initFunc(r)
	}

	b.ResetTimer()
	for {
		_, err := r.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReadWithHeader(b *testing.B) {
	runReadBenchmark(b, []string{"f1", "f2", "f3", "f4"}, nil)
}

func BenchmarkReadWithoutHeader(b *testing.B) {
	runReadBenchmark(b, []string{"f1", "f2", "f3", "f4"}, nil)
}

func BenchmarkReadWithHeaderReuseRecord(b *testing.B) {
	runReadBenchmark(b, []string{"f1", "f2", "f3", "f4"}, func(r *mcsv.Reader) { r.ReuseRecord = true })
}

func BenchmarkReadWithoutHeaderReuseRecord(b *testing.B) {
	runReadBenchmark(b, []string{"f1", "f2", "f3", "f4"}, func(r *mcsv.Reader) { r.ReuseRecord = true })
}

type writerTest struct {
	name   string
	header []string
	input  []map[string]string
	output string

	writeHeader bool
	comma       rune // if non zero, set Writer.Comma
}

var writerTests = []writerTest{
	{
		name:   "EuCities",
		header: []string{"City", "Country"},

		input: []map[string]string{
			{"City": "Berlin", "Country": "Germany"},
			{"City": "Madrid", "Country": "Spain"},
			{"City": "Rome", "Country": "Italy"},
			{"City": "Bucharest", "Country": "Romania"},
			{"City": "Paris", "Country": "France"},
		},

		output: `City,Country
Berlin,Germany
Madrid,Spain
Rome,Italy
Bucharest,Romania
Paris,France
`,

		writeHeader: true,
	},

	{
		name:   "MathConstants",
		header: []string{"name", "value"},

		input: []map[string]string{
			{"name": "pi", "value": "3.1416"},
			{"name": "sqrt2", "value": "1.4142"},
			{"name": "phi", "value": "1.618"},
			{"name": "e", "value": "2.7183"},
		},

		output: `pi,3.1416
sqrt2,1.4142
phi,1.618
e,2.7183
`,

		writeHeader: false,
	},

	{
		name:   "MetroSystems",
		header: []string{"City", "Stations", "System Length"},

		input: []map[string]string{
			{"City": "New York", "Stations": "424", "System Length": "380"},
			{"City": "Shanghai", "Stations": "345", "System Length": "676"},
			{"City": "Seoul", "Stations": "331", "System Length": "353"},
			{"City": "Beijing", "Stations": "326", "System Length": "690"},
			{"City": "Paris", "Stations": "302", "System Length": "214"},
			{"City": "London", "Stations": "270", "System Length": "402"},
		},

		output: `City	Stations	System Length
New York	424	380
Shanghai	345	676
Seoul	331	353
Beijing	326	690
Paris	302	214
London	270	402
`,

		writeHeader: true,
		comma:       '\t',
	},

	{
		name:   "Newlines",
		header: []string{"field1", "field2", "field3"},

		input: []map[string]string{
			{"field1": "hello", "field2": "is it \"me\"", "field3": "you're\nlooking for"},
			{"field1": "this is going to be", "field2": "another\nbroken row", "field3": "very confusing"},
		},

		output: `field1,field2,field3
hello,"is it ""me""","you're
looking for"
this is going to be,"another
broken row",very confusing
`,

		writeHeader: true,
	},

	{
		name:   "MismatchedFields",
		header: []string{"country", "population"},

		input: []map[string]string{
			{"country": "France", "population": "68000000", "capitol": "Paris"},
			{"country": "Germany", "population": "84000000"},
			{"country": "Spain", "capitol": "Madrid"},
		},

		output: `country,population
France,68000000
Germany,84000000
Spain,
`,

		writeHeader: true,
	},
}

func runWriterTest(t *testing.T, test writerTest, writeAll bool) {
	// Create the writer
	out := &strings.Builder{}
	w := mcsv.NewWriter(out, test.header)

	if test.comma != 0 {
		w.Comma = test.comma
	}

	if test.writeHeader {
		err := w.WriteHeader()
		assert.NoErr(t, err)
	}

	if writeAll {
		err := w.WriteAll(test.input)
		assert.NoErr(t, err)
	} else {
		for _, record := range test.input {
			err := w.Write(record)
			assert.NoErr(t, err)
		}

		w.Flush()
		assert.NoErr(t, w.Error())
	}

	check.Eq(t, out.String(), test.output)
}

func ExampleWriter() {
	data := []map[string]string{
		{"City": "Berlin", "Country": "Germany", "Population": "3 677 472"},
		{"City": "Madrid", "Country": "Spain", "Population": "3 223 334"},
		{"City": "Rome", "Country": "Italy"},
		{"City": "Bucharest", "Country": "Romania", "Area": "240 km²"},
		{"City": "Paris", "Country": "France", "Population": "2 165 423", "Area": "105 km²"},
	}

	w := mcsv.NewWriter(os.Stdout, []string{"City", "Country", "Population"})

	// Write the header row
	err := w.WriteHeader()
	if err != nil {
		log.Fatalln(err)
	}

	// Write every record
	for _, record := range data {
		err = w.Write(record)
		if err != nil {
			log.Fatalln(err)
		}
	}

	// Flush the output
	w.Flush()
	err = w.Error()
	if err != nil {
		log.Fatalln(err)
	}

	// Output:
	// City,Country,Population
	// Berlin,Germany,3 677 472
	// Madrid,Spain,3 223 334
	// Rome,Italy,
	// Bucharest,Romania,
	// Paris,France,2 165 423
}

func TestWrite(t *testing.T) {
	for _, test := range writerTests {
		t.Run(test.name, func(t *testing.T) { runWriterTest(t, test, false) })
	}
}

func TestWriteAll(t *testing.T) {
	for _, test := range writerTests {
		t.Run(test.name, func(t *testing.T) { runWriterTest(t, test, true) })
	}
}

var writeBenchmarkRecords = []map[string]string{
	{"name": "pi", "value": "3.1416"},
	{"name": "sqrt2", "value": "1.4142"},
	{"name": "phi", "value": "1.618"},
	{"name": "e", "value": "2.7183"},
	{"name": "i", "value": "0+1i"},
}

func BenchmarkWrite(b *testing.B) {
	// NOTE: There are no performance-related settings (like what Reader has),
	//       and so there is only this one write benchmark.

	w := mcsv.NewWriter(io.Discard, []string{"name", "value"})
	for i := 0; i < b.N; i++ {
		for _, record := range writeBenchmarkRecords {
			w.Write(record)
		}
	}
}
