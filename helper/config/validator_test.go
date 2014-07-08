package config

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/terraform"
)

func TestValidator(t *testing.T) {
	v := &Validator{
		Required: []string{"foo"},
		Optional: []string{"bar"},
	}

	var c *terraform.ResourceConfig

	// Valid
	c = testConfig(t, map[string]interface{}{
		"foo": "bar",
	})
	testValid(v, c)

	// Valid + optional
	c = testConfig(t, map[string]interface{}{
		"foo": "bar",
		"bar": "baz",
	})
	testValid(v, c)

	// Missing required
	c = testConfig(t, map[string]interface{}{
		"bar": "baz",
	})
	testInvalid(v, c)

	// Unknown key
	c = testConfig(t, map[string]interface{}{
		"foo":  "bar",
		"what": "what",
	})
	testInvalid(v, c)
}

func TestValidator_complex(t *testing.T) {
	v := &Validator{
		Required: []string{
			"foo",
			"nested.*",
		},
	}

	var c *terraform.ResourceConfig

	// Valid
	c = testConfig(t, map[string]interface{}{
		"foo": "bar",
		"nested": []map[string]interface{}{
			map[string]interface{}{"foo": "bar"},
		},
	})
	testValid(v, c)

	// Not a nested structure
	c = testConfig(t, map[string]interface{}{
		"foo":    "bar",
		"nested": "baa",
	})
	testInvalid(v, c)
}

func testConfig(
	t *testing.T,
	c map[string]interface{}) *terraform.ResourceConfig {
	r, err := config.NewRawConfig(c)
	if err != nil {
		t.Fatalf("bad: %s", err)
	}

	return terraform.NewResourceConfig(r)
}

func testInvalid(v *Validator, c *terraform.ResourceConfig) {
	ws, es := v.Validate(c)
	if len(ws) > 0 {
		panic(fmt.Sprintf("bad: %#v", ws))
	}
	if len(es) == 0 {
		panic(fmt.Sprintf("bad: %#v", es))
	}
}

func testValid(v *Validator, c *terraform.ResourceConfig) {
	ws, es := v.Validate(c)
	if len(ws) > 0 {
		panic(fmt.Sprintf("bad: %#v", ws))
	}
	if len(es) > 0 {
		panic(fmt.Sprintf("bad: %#v", es))
	}
}
