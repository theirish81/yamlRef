package yamlRef

import (
	"os"
	"testing"
)

func TestMergeAndMarshall(t *testing.T) {
	data, _ := MergeAndMarshall("test_data/main.yaml")
	res, _ := os.ReadFile("test_data/outcome.yaml")
	if string(data) != string(res) {
		t.Error("Wrong merge and marshall")
	}
}

func TestExtractPathFromRef(t *testing.T) {
	extracted, _ := refToUrl("$ref:file://foo/bar", "")
	if extracted.Host+extracted.Path != "foo/bar" {
		t.Error("Wrong path")
	}
	extracted, _ = refToUrl("$ref:file://foo", "")
	if extracted.Host+extracted.Path != "foo" {
		t.Error("Wrong path")
	}

	extracted, _ = refToUrl("$ref:file://foo/bar", "/dope")
	if extracted.Host+extracted.Path != "/dope/foo/bar" {
		t.Error("Wrong path")
	}

	extracted, _ = refToUrl("$ref:file://foo/bar?comp=foobar", "/dope")
	if extracted.Query()["comp"][0] != "foobar" {
		t.Error("Wrong query param")
	}

	extracted, _ = refToUrl("$ref:file:///foo/bar", "/dope")
	if extracted.Host+extracted.Path != "/foo/bar" {
		t.Error("Wrong path")
	}
}
