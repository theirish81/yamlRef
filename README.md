# YamlRef

An easy Golang library to reference and merge multiple YAML files into a main one.

Example (main.yaml):
```yaml
rootObjct:
  foo: bar
  extRef: "$ref:file://external1.yaml"
  myArray:
    - foo
    - bar
    - "$ref:file://external1.yaml"
```
And (external1.yaml):
```yaml
externalObject:
  externalFoo: externalBar
  externalArray:
    - a
    - b
    - c
```
Invoking `MergeAndMarshall("main.yaml")` would produce the following output (in []byte):

```yaml
rootObjct:
  extRef:
    externalObject:
      externalArray:
      - a
      - b
      - c
      externalFoo: externalBar
  foo: bar
  myArray:
  - foo
  - bar
  - externalObject:
      externalArray:
      - a
      - b
      - c
      externalFoo: externalBar
```

You can also invoke the `Merge(path string)` function to obtain the unmarshalled data structure.