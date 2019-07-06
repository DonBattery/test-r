package cmd

import (
	"fmt"
	"os"

	"github.com/donbattery/test-r/reporter"
	"github.com/spf13/cobra"
)

// xmlCmd represents the xml command
var xmlCmd = &cobra.Command{
	Use:   "xml",
	Short: "Generates JUnit xml",
	Long:  `it just does`,
	Run: func(cmd *cobra.Command, args []string) {
		generateXML()
	},
}

func init() {
	RootCmd.AddCommand(xmlCmd)
}

func generateXML() {
	junit := reporter.NewReporter(os.Stdout)
	suite1 := reporter.NewSuite("App 1")
	suite1.AddProperty("Stack", "Linux")
	suite1.AddTestCase(*reporter.NewTestCase("Builds1", "On valamilyen stack", "123s"))
	suite2 := reporter.NewSuite("App 2 hah")
	failCase := reporter.NewTestCase("Builds 1", "On valamilyen stack", "15435s")
	failCase.Fail("Ez elfailelt", "Fatal", "hu de elfaileltem 911")
	suite2.AddTestCase(*failCase)
	suite2.AddTestCase(*reporter.NewTestCase("Builds 2", "On masmilyen stack", "1s"))
	suite2.AddTestCase(*reporter.NewTestCase("Builds 3", "On ugyanolyan stack", "435s"))
	// suite3 := reporter.NewSuite("App 3 heh").WithProperty("Prop1", "Val1").WithProperty("prop2", "val2").WithProperty("Prop3", "val3")
	// suite3.AddTestCase(reporter.NewTestCase(""))
	junit.AddSuite(*suite1)
	junit.AddSuite(*suite2)
	if err := junit.GenerateXML(); err != nil {
		fmt.Printf("FAtal ERROROROrorororor\n")
	}
}
