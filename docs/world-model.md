# World Model

The world model describes everything that is available for you to reason about in your test expectations. It contains,
among other things, the parsed rendering output, the message generated from `NOTES.txt`, error messages (if any)
as well as the input values. The "world" is given as a JSON object, and is used as the input for every `jq` filter
that is evaluated as an expectation.

This JSON object contains the following properties:
- `helm`: a JSON object representing the values passed to the Helm rendering engine, such as `Values`, `Release`,
  `Chart` etc.
- `notesRaw`: a string containing the message that would be displayed after running `helm install` or `helm upgrade`,
  rendered from the `NOTES.txt` template.
- `notes`: same as `notesRaw`, but in normalized form - all sequences of one or more whitespaces (spaces, tabs, line
  breaks) having been replaced by a single space (`' '`) character. This allows checking for the occurrence of phrases
  without paying attention to line wrapping, and is almost always to preferable over `notesRaw`.
- `errorRaw`: a string containing the Helm-generated error message occurred from rendering templates (if any). Note that
  schema validation errors do not end up here.
- `error`: same as `errorRaw`, but in normalized form (see `notes` above).
- `objects`: an array of all the objects in the rendering output, parsed directly from YAML.

Additionally, for every object kind that occurred in the rendering output, the world object contains a field named after
the plural form in lowercase. This field references a name-indexed object of all Kubernetes object of the given
resource, provided the name is unique across all namespaces. That is, to locate the deployment "server", you can thus
either write
```
.objects[] | select(.metadata.kind == "Deployment" and .metadata.name == "server")
```
or simply `.deployments.server`. If there are deployments named `server` in multiple distinct
namespaces (which shouldn't usually be the case with Helm charts), however, `.deployments.server` will be undefined.
