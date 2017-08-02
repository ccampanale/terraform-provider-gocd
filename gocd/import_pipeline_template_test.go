package gocd

import (
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

func TestResourcePipelineTemplate_Import(t *testing.T) {

}

func testAccCheckPipelineTemplateDestroy(s *terraform.State) error {
	//pt := testGocdProvider.Meta().(*gocd.Client).PipelineTemplates
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "gocd_pipeline_template" {
			continue
		}
	}

	return nil
}
