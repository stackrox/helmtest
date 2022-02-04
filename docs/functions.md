Functions
===============

As `helmtest` predicates are `jq` filters,
[all the functions available in `jq`](https://stedolan.github.io/jq/manual/#Builtinoperatorsandfunctions) are available
in `helmtest` as well.

Additionally, `helmtest` adds a number of additional functions to facilitate writing test cases for Kubernetes
manifests:

| Function name(s)     | Description |
| -------------------- | ----------- |
| `fromyaml`, `toyaml` | The equivalent of `fromjson` and `tojson` for YAML. |
| `assertNotExist`     | This function will fail if ever executed. Semantically equivalent to just writing `false`, but prints the offending object. |
| `assertThat(f)`      | Asserts that a filter `f` holds for the input object. If `. \| f` evaluates to `false`, this will print the value of `.` as well as the original string representation of `f`. Hence, while `... \| .name == "foo"` and `\| assertThat(.name == "foo")` are semantically equivalent, the latter is preferable as it is much easier to debug. |
| `assumeThat(f)`      | Assumes that a filter `f` holds for the input object. If it doesn't, the evaluation is aborted for the given input object, and no failure is triggered. |
| `print`              | Prints input directly with `fmt.Println` and returns it. To print all objects in a test as `yaml`, write `.objects[] \| toyaml \| print`. |
