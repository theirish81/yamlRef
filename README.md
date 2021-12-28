# YamlRef

An easy Golang library to reference and merge multiple local YAML files into a main one.

## $ref links

To perform the replacement, use a "$ref string" as value in any item of your YAML as in:

* `$ref:file://external.yaml`: replace this $ref string with the value of the `external.yaml` file, located in the same
directory as the main file
* `$ref:file:///home/theirish81/external.yaml` replace this $ref with the value of the `external.yaml` file, located in 
an absolute path
* `$ref:file://external.yaml?comp=externalBot` replace this $ref with a specific object `externalBot` (root level only)
described in the `external.yaml` file, located in the same directory as the main file

## Example

**main.yaml**:

```yaml
rootObjct:
  foo: bar
  extRef: "$ref:file://external1.yaml"
  myArray:
    - foo
    - bar
    - "$ref:file://external1.yaml?comp=externalBot"
```
**external1.yaml**:

```yaml
externalObject:
  externalFoo: externalBar
  externalArray:
    - a
    - b
    - c
externalBot:
  bot: true
```
Invoking `MergeAndMarshall("main.yaml")` would produce the following output (in []byte):

**outcome**:

```yaml
rootObjct:
  extRef:
    externalBot:
      bot: true
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
    - bot: true
```

You can also invoke the `Merge(path string)` function to obtain the unmarshalled data structure.