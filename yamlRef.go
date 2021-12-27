package yamlRef

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	url2 "net/url"
	"path"
	"strings"
)

func MergeAndMarshall(path string) ([]byte, error) {
	data, err := Merge(path)
	if err != nil {
		return []byte{}, err
	}
	return yaml.Marshal(data)
}

func Merge(filePath string) (interface{}, error) {
	mainBytes, err := ioutil.ReadFile(filePath)
	var data map[interface{}]interface{}
	if err != nil {
		return data, err
	}
	err = yaml.Unmarshal(mainBytes, &data)
	if err != nil {
		return nil, err
	}
	res, err := findAndReplace(data, path.Dir(filePath))
	return res, err
}

func findAndReplace(data interface{}, basePath string) (interface{}, error) {
	switch data := data.(type) {
	case map[interface{}]interface{}:
		for k, v := range data {
			if strVal, ok := v.(string); strings.HasPrefix(strVal, "$ref:") && ok {
				resPath, err := extractPathFromRef(strVal, basePath)
				if err != nil {
					return nil, err
				}
				replaceData, err := Merge(resPath)
				if err != nil {
					return data, err
				}
				data[k] = replaceData
			} else {
				dataL2, err := findAndReplace(v, basePath)
				if err != nil {
					return nil, err
				}
				data[k] = dataL2
			}
		}
	case []interface{}:
		for idx, v := range data {
			if strVal, ok := v.(string); strings.HasPrefix(strVal, "$ref:") && ok {
				resPath, err := extractPathFromRef(strVal, basePath)
				if err != nil {
					return nil, err
				}
				replaceData, err := Merge(resPath)
				if err != nil {
					return data, err
				}
				data[idx] = replaceData
			}
		}
	}
	return data, nil
}

func extractPathFromRef(ref string, basePath string) (string, error) {
	ref = ref[5:]
	url, err := url2.Parse(ref)
	if err != nil {
		return "", err
	}
	if url.Scheme == "file" && !strings.HasPrefix(url.String(), "file:///") {
		url, err = url2.Parse("file://" + path.Join(basePath, url.Host+url.Path))

	}
	return url.Host + url.Path, err
}
