package yamlRef

import (
	"errors"
	"gopkg.in/yaml.v2"
	url2 "net/url"
	"os"
	"path"
	"strings"
)

// MergeAndMarshall will take the path to a YAML file, look for any $ref and perform the merges accordingly.
// Once done, it will marshall the data structure back to proper YAML
func MergeAndMarshall(path string) ([]byte, error) {
	data, err := Merge(path)
	if err != nil {
		return []byte{}, err
	}
	return yaml.Marshal(data)
}

// Merge will take the path to a YAML file, look for any $ref and perform the merges accordingly.
// Once done, it will return the raw data structure as an interface{}
func Merge(filePath string) (interface{}, error) {
	// converting a path to a file:// URL
	url, err := url2.Parse("file://" + filePath)
	if err != nil {
		return nil, err
	}
	return merge(url)
}

// merge will take the path to a YAML file in the form of a URL, look for any $ref and perform the
// merges accordingly. Once done, it will return the raw data structure as an interface{}
func merge(url *url2.URL) (interface{}, error) {
	mainBytes, err := os.ReadFile(url.Host + url.Path)
	var data map[interface{}]interface{}
	if err != nil {
		return data, err
	}
	err = yaml.Unmarshal(mainBytes, &data)
	if err != nil {
		return nil, err
	}
	res, err := findAndReplace(data, path.Dir(url.Host+url.Path))
	if err != nil {
		return nil, err
	}
	// if the URL contains a "comp" query param, we want to reference a specific object in the imported YAML
	if comp, ok := url.Query()["comp"]; ok {
		// we need to assume that the loaded data structure is in fact a map
		if casted, ok := res.(map[interface{}]interface{}); ok && len(comp) > 0 {
			// extracting the specific object, if it's in fact present
			if obj, ok := casted[comp[0]]; ok {
				return obj, err
			} else {
				// if it's not present, then we'll need to throw an error
				return nil, errors.New("comp not found in referenced object")
			}

		} else {
			return nil, errors.New("referenced YAML file does not contain a map or comp is invalid")
		}

	}
	// if the URL DOES NOT contain a "comp" query param, then we're good, and we can return the whole data structure
	return res, nil
}

// findAndReplace will recursively look for $refs and replace them with the loaded data structure.
// The updated data structure is returned, in the form of an interface{}
func findAndReplace(data interface{}, basePath string) (interface{}, error) {
	switch data := data.(type) {
	// If it's a map, we need to drill down
	case map[interface{}]interface{}:
		for k, v := range data {
			// If the value is a string and starts with $ref...
			if strVal, ok := v.(string); strings.HasPrefix(strVal, "$ref:") && ok {
				// ... we convert the $ref into a proper URL
				resPath, err := refToUrl(strVal, basePath)
				if err != nil {
					return nil, err
				}
				// we load the referenced data structure, and we make sure to process refs we may find in there
				replaceData, err := merge(resPath)
				if err != nil {
					return data, err
				}
				data[k] = replaceData
			} else {
				// If the value is not a string or does not start with $ref, we recursively dig into it
				dataL2, err := findAndReplace(v, basePath)
				if err != nil {
					return nil, err
				}
				data[k] = dataL2
			}
		}
	// In case we find an array, we go through each item looking for $refs
	case []interface{}:
		for idx, v := range data {
			// If the item ia string and starts with $ref...
			if strVal, ok := v.(string); strings.HasPrefix(strVal, "$ref:") && ok {
				// ... we convert the $ref into a proper URL
				resPath, err := refToUrl(strVal, basePath)
				if err != nil {
					return nil, err
				}
				// we load the referenced data structure, and we make sure to process refs we may find in there
				replaceData, err := merge(resPath)
				if err != nil {
					return data, err
				}
				// finally, we replace the item at the given index
				data[idx] = replaceData
			}
		}
	}
	return data, nil
}

// refToUrl will receive a $ref string and, given a basePath, will return a URL for the referenced YAML file
func refToUrl(ref string, basePath string) (*url2.URL, error) {
	// remove the $ref: prefix
	ref = ref[5:]
	// turn it into a URL
	url, err := url2.Parse(ref)
	if err != nil {
		return nil, err
	}
	// if it's a file:// URL, and it's not absolute...
	if url.Scheme == "file" && !strings.HasPrefix(url.String(), "file:///") {
		//... we combine the provided path with the base path
		url, err = url2.Parse("file://" + path.Join(basePath, url.Host+url.Path) + "?" + url.RawQuery)

	}
	return url, err
}
