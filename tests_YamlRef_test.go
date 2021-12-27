package yamlRef

import (
	"io/ioutil"
	"testing"
)

func TestLoad2(t *testing.T) {
	data, _ := MergeAndMarshall("test_data/main.yaml")
	res, _ := ioutil.ReadFile("test_data/outcome.yaml")
	if string(data) != string(res) {
		t.Error("Wrong merge and marshall")
	}
}

func TestExtractPathFromRef(t *testing.T) {
	extracted, _ := extractPathFromRef("$ref:file://foo/bar", "")
	if extracted != "foo/bar" {
		t.Error("Wrong path")
	}
	extracted, _ = extractPathFromRef("$ref:file://foo", "")
	if extracted != "foo" {
		t.Error("Wrong path")
	}

	extracted, _ = extractPathFromRef("$ref:file://foo/bar", "/dope")
	if extracted != "/dope/foo/bar" {
		t.Error("Wrong path")
	}
}
