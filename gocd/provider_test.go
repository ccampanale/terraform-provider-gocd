package gocd

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/sergi/go-diff/diffmatchpatch"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"
	"github.com/beamly/go-gocd/gocd"
)

var (
	testGocdProviders map[string]terraform.ResourceProvider
	testGocdProvider  *schema.Provider
)

type TestStepJSONComparison struct {
	Index        int
	ID           string
	Config       string
	ExpectedJSON string
}

func init() {

	testGocdProvider = Provider().(*schema.Provider)
	testGocdProviders = map[string]terraform.ResourceProvider{
		"gocd": testGocdProvider,
	}

}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

	t.Run("Impl", testProviderImpl)
}

func testProviderImpl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testGetClient() *gocd.Client {
	//GOCD_URL=http://127.0.0.1:8153/go/;GOCD_SKIP_SSL_CHECK=1;TF_ACC=1

	var url, u, p string
	var nossl bool

	url = os.Getenv("GOCD_URL")
	u = os.Getenv("GOCD_USERNAME")
	p = os.Getenv("GOCD_PASSWORD")
	nossl = os.Getenv("GOCD_SKIP_SSL_CHECK") == "1"

	cfg := &gocd.Configuration{
		Server:       url,
		Username:     u,
		Password:     p,
		SkipSslCheck: nossl,
	}

	return cfg.Client()
}

func testAccPreCheck(t *testing.T) {

	if u := os.Getenv("GOCD_URL"); u == "" {
		t.Fatal("GOCD_URL must be set for acceptance tests.")
	}

	if s := os.Getenv("GOCD_SKIP_SSL_CHECK"); s == "" {
		t.Fatal("GOCD_SKIP_SSL_CHECK must be set for acceptance tests.")
	}

	err := testGocdProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}
}

// Loads a test file resource from the 'test' directory.
func testFile(name string) string {
	f, err := ioutil.ReadFile(fmt.Sprintf("test/%s", name))
	if err != nil {
		panic(err)
	}

	return string(f)
}

func testTaskDataSourceStateValue(id string, name string, value string, index int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		root := s.RootModule()
		rs, ok := root.Resources[id]
		if !ok {
			return fmt.Errorf("In '%d'.\nNot found: %s", index, id)
		}
		if rs.Primary.ID == "" {
			return fmt.Errorf("In '%d'.\nNo ID is set", index)
		}

		if v := rs.Primary.Attributes[name]; v != value {
			dmp := diffmatchpatch.New()
			rawDiffs := dmp.DiffMain(v, value, true)
			rawDiff := dmp.DiffPrettyText(rawDiffs)

			err := fmt.Errorf("In '%d'.\nValue mismatch for 'json' is:\n%s", index, rawDiff)
			return err
		}

		return nil
	}
}

func testStepComparisonCheck(t *TestStepJSONComparison) []resource.TestStep {
	return []resource.TestStep{
		{
			Config: t.Config,
			Check: func(s *terraform.State) error {
				root := s.RootModule()
				rs, ok := root.Resources[t.ID]
				if !ok {
					return fmt.Errorf("In '%d'.\nNot found: %s", t.Index, t.ID)
				}
				if rs.Primary.ID == "" {
					return fmt.Errorf("In '%d'.\nNo ID is set", t.Index)
				}

				if v := rs.Primary.Attributes["json"]; v != t.ExpectedJSON {
					dmp := diffmatchpatch.New()
					rawDiffs := dmp.DiffMain(v, t.ExpectedJSON, true)
					rawDiff := dmp.DiffPrettyText(rawDiffs)

					err := fmt.Errorf("In '%d'.\nValue mismatch for 'json' is:\n'%s'", t.Index, rawDiff)
					return err
				}

				return nil
			},
		},
	}
}

func randomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyz"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}
