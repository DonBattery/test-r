package model

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"runtime"
	"time"
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

// JUnitReportXML writes a JUnit xml representation of the given report to w
func JUnitReportXML(report *Report, noXMLHeader bool, w io.Writer) error {
	suites := JUnitTestSuites{}

	// convert Report to JUnit test suites
	for _, app := range report.Apps {
		ts := JUnitTestSuite{
			Tests:      len(app.Builds),
			Failures:   0,
			Time:       formatTime(app.Duration),
			Name:       app.Name,
			Properties: []JUnitProperty{},
			TestCases:  []JUnitTestCase{},
		}

		ts.Properties = append(ts.Properties, JUnitProperty{"go.version", runtime.Version()})

		// individual builds
		for _, build := range app.Builds {
			testCase := JUnitTestCase{
				Classname: app.Name,
				Name:      build.Name,
				Time:      formatTime(build.Duration),
				Failure:   nil,
			}

			if build.Result == FAIL {
				ts.Failures++
				testCase.Failure = &JUnitFailure{
					Message:  "Failed",
					Type:     "Build",
					Contents: "This build failed misserably",
				}
			}

			ts.TestCases = append(ts.TestCases, testCase)
		}

		suites.Suites = append(suites.Suites, ts)
	}

	// to xml
	bytes, err := xml.MarshalIndent(suites, "", "\t")
	if err != nil {
		return err
	}

	writer := bufio.NewWriter(w)

	if !noXMLHeader {
		writer.WriteString(xml.Header)
	}

	writer.Write(bytes)
	writer.WriteByte('\n')
	writer.Flush()

	return nil
}

func formatTime(d time.Duration) string {
	return fmt.Sprintf("%.3f", d.Seconds())
}
