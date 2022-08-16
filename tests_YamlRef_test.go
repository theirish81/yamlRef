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

	if data, _ := MergeAndMarshall("test_data/arr_ref.yaml"); string(data) != "rootObject:\n  ref:\n  - foo\n  - bar\n" {
		t.Error("could not import an array")
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

func TestNegative(t *testing.T) {
	if _, err := MergeAndMarshall("test_data/no_ref.yaml"); err == nil {
		t.Error("non existing ref should throw an error ")
	}
	if _, err := MergeAndMarshall("test_data/no_comp.yaml"); err == nil {
		t.Error("non existing ref should throw an error ")
	}
	if _, err := MergeAndMarshall("test_data/no_arr_comp.yaml"); err == nil {
		t.Error("non existing ref should throw an error ")
	}
	if _, err := MergeAndMarshall("test_data/gibberish_ref.yaml"); err == nil {
		t.Error("broken contributing file should return an error")
	}
	if _, err := MergeAndMarshall("test_data/no_ref.yaml"); err == nil {
		t.Error("broken contributing file should return an error")
	}
}
