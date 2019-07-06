package reporter

import (
	"bufio"
	"encoding/xml"
	"io"

	"github.com/pkg/errors"
)

// This model is heavily based on
// https://github.com/jstemmer/go-junit-report

// JUnitTestSuites is a collection of JUnit test suites.
type JUnitTestSuites struct {
	XMLName xml.Name `xml:"testsuites"`
	Suites  []JUnitTestSuite
}

// JUnitTestSuite is a single JUnit test suite which may contain many
// testcases.
type JUnitTestSuite struct {
	XMLName    xml.Name        `xml:"testsuite"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Time       string          `xml:"time,attr"`
	Name       string          `xml:"name,attr"`
	Properties []JUnitProperty `xml:"properties>property,omitempty"`
	TestCases  []JUnitTestCase
}

// JUnitTestCase is a single test case with its result.
type JUnitTestCase struct {
	XMLName     xml.Name          `xml:"testcase"`
	Classname   string            `xml:"classname,attr"`
	Name        string            `xml:"name,attr"`
	Time        string            `xml:"time,attr"`
	SkipMessage *JUnitSkipMessage `xml:"skipped,omitempty"`
	Failure     *JUnitFailure     `xml:"failure,omitempty"`
}

// JUnitSkipMessage contains the reason why a testcase was skipped.
type JUnitSkipMessage struct {
	Message string `xml:"message,attr"`
}

// JUnitProperty represents a key/value pair used to define properties.
type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

// JUnitFailure contains data related to a failed test.
type JUnitFailure struct {
	Message  string `xml:"message,attr"`
	Type     string `xml:"type,attr"`
	Contents string `xml:",chardata"`
}

// NewTestCase returns a pointer to a new TestCase initialized with the supplyed name, classname and time
func NewTestCase(name, classname, time string) *JUnitTestCase {
	return &JUnitTestCase{
		XMLName: xml.Name{
			Local: name,
		},
		Name:      name,
		Classname: classname,
		Time:      time,
	}
}

// FailTest adds a Failure to the TestCase
func (tc *JUnitTestCase) FailTest(message, ftype, contents string) {
	tc.Failure = &JUnitFailure{
		Message:  message,
		Type:     ftype,
		Contents: contents,
	}
}

// Fail fails the test Case and returns it
func (tc *JUnitTestCase) Fail(message, ftype, contents string) *JUnitTestCase {
	tc.FailTest(message, ftype, contents)
	return tc
}

// NewSuite returns a pointer to a new JUnitTestSuite initialized with the supplyed name
func NewSuite(name string) *JUnitTestSuite {
	return &JUnitTestSuite{
		XMLName: xml.Name{
			Local: name,
		},
		Name: name,
	}
}

// AddTestCase adds a test case to the Suite
func (su *JUnitTestSuite) AddTestCase(testCase JUnitTestCase) {
	su.TestCases = append(su.TestCases, testCase)
}

// AddProperty adds a Property to the Suite
func (su *JUnitTestSuite) AddProperty(key, value string) {
	su.Properties = append(su.Properties, JUnitProperty{
		Name:  key,
		Value: value,
	})
}

// WithProperty adds a property to the Suite and returns it
func (su *JUnitTestSuite) WithProperty(key, value string) *JUnitTestSuite {
	su.AddProperty(key, value)
	return su
}

// Reporter is storing JUnit Suites and is able to write them to its writer
type Reporter struct {
	writer io.Writer
	suites JUnitTestSuites
}

// NewReporter returns a new pointer to a new Reporter initialized with the supplyed writer
func NewReporter(output io.Writer) *Reporter {
	return &Reporter{
		writer: output,
		suites: JUnitTestSuites{
			XMLName: xml.Name{
				Local: "Testman",
			},
		},
	}
}

// AddSuite adds a JUnit Suite to the Reporter
func (r *Reporter) AddSuite(suite JUnitTestSuite) {
	r.suites.Suites = append(r.suites.Suites, suite)
}

// GenerateXML writes a JUnit xml representation of the Suites onto the writer
func (r *Reporter) GenerateXML() error {

	// Encode as xml
	xmlBytes, err := xml.MarshalIndent(r.suites, "", "\t")
	if err != nil {
		return errors.Wrap(err, "Cannot generate xml output")
	}

	// Write xml onto the supplied output
	writer := bufio.NewWriter(r.writer)
	if _, err := writer.WriteString(xml.Header); err != nil {
		return errors.Wrap(err, "Failed to write xml header to buffer")
	}
	if _, err := writer.Write(xmlBytes); err != nil {
		return errors.Wrap(err, "Failed to write xml body to buffer")
	}
	if err := writer.WriteByte('\n'); err != nil {
		return errors.Wrap(err, "Failed to write xml final newline to buffer")
	}
	if err := writer.Flush(); err != nil {
		return errors.Wrap(err, "Failed to write buffer to output")
	}

	return nil
}
