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
	url, err := url2.Parse("file://" + filePath)
	if err != nil {
		return nil, err
	}
	return merge(url)
}

func merge(url *url2.URL) (interface{}, error) {
	mainBytes, err := ioutil.ReadFile(url.Host + url.Path)
	var data map[interface{}]interface{}
	if err != nil {
		return data, err
	}
	err = yaml.Unmarshal(mainBytes, &data)
	if err != nil {
		return nil, err
	}
	res, err := findAndReplace(data, path.Dir(url.Host+url.Path))
	if comp, ok := url.Query()["comp"]; ok {
		casted := res.(map[interface{}]interface{})
		obj := casted[comp[0]]
		return obj, err
	}
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
				replaceData, err := merge(resPath)
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
				replaceData, err := merge(resPath)
				if err != nil {
					return data, err
				}
				data[idx] = replaceData
			}
		}
	}
	return data, nil
}

func extractPathFromRef(ref string, basePath string) (*url2.URL, error) {
	ref = ref[5:]
	url, err := url2.Parse(ref)
	if err != nil {
		return nil, err
	}
	if url.Scheme == "file" && !strings.HasPrefix(url.String(), "file:///") {
		url, err = url2.Parse("file://" + path.Join(basePath, url.Host+url.Path) + "?" + url.RawQuery)

	}
	return url, err
}
