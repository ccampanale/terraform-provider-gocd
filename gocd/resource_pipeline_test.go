package gocd

import (
	r "github.com/hashicorp/terraform/helper/resource"
	"regexp"
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
)

func testResourcePipeline(t *testing.T) {
	t.Run("Basic", testResourcePipelineBasic)
	t.Run("ImportBasic", testResourcePipelineImportBasic)
	t.Run("FullStack1", testResourcePipelineFullStack1)
	t.Run("FullStack2", testResourcePipelineFullStack2)
	t.Run("DisableAutoUpdate", testResourcePipelineDisableAutoUpdate)
	t.Run("LinkedDependencies", testResourcePipelineLinkedDependencies)
	t.Run("HandleMissingPipelines", testResourcePipelineMissing)
}

func testResourcePipelineMissing(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline.0.rsc.tf"),
			},
			{
				PreConfig: func() {
					c := testGetClient()
					_, _, err := c.PipelineConfigs.Delete(context.Background(), "pipeline0-terraform")
					assert.NoError(t, err)
				},
				Config: testFile("resource_pipeline.4.rsc.tf"),
			},
		},
	})
}

func testResourcePipelineLinkedDependencies(t *testing.T) {
	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_linked_dep.0.rsc.tf"),
			},
			{
				Config: testFile("resource_pipeline_linked_dep.1.rsc.tf"),
			},
		},
	})
}

func testResourcePipelineDisableAutoUpdate(t *testing.T) {
	// Managing auto_update on a per material basis is not possible through the current GoCD API as of 01/10/2017.
	// For details see, https://github.com/beamly/terraform-provider-gocd/issues/32
	errRE, err := regexp.Compile("The `auto_update` attribute has been disabled until a way to manage updates atomically has been devised")
	if err != nil {
		t.Error(err)
	}

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline_auto_update.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr(
						"gocd_pipeline.pipeline1",
						"name",
						"pipeline1",
					),
				),
				ExpectError: errRE,
			},
		},
	})
}

func testResourcePipelineBasic(t *testing.T) {

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline.0.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr(
						"gocd_pipeline.test-pipeline",
						"name",
						"pipeline0-terraform",
					),
				),
			},
			{
				Config: testFile("resource_pipeline.1.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr(
						"gocd_pipeline.test-pipeline",
						"name",
						"pipeline0-terraform",
					),
				),
			},
		},
	})
}

func testResourcePipelineFullStack1(t *testing.T) {

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline.3.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr(
						"gocd_pipeline.test-pipeline3",
						"name",
						"test-pipeline3",
					),
				),
			},
		},
	})
}
func testResourcePipelineFullStack2(t *testing.T) {

	r.Test(t, r.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testGocdProviders,
		CheckDestroy: testGocdPipelineDestroy,
		Steps: []r.TestStep{
			{
				Config: testFile("resource_pipeline.4.rsc.tf"),
				Check: r.ComposeTestCheckFunc(
					r.TestCheckResourceAttr(
						"gocd_pipeline.test-pipeline",
						"name",
						"test-pipeline",
					),
				),
			},
		},
	})
}
