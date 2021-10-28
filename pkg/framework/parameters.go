package framework

import (
	"github.com/pkg/errors"
	"github.com/stackrox/helmtest/internal/rox-imported/stringutils"
	"strings"
)

func expandParameters(params interface{}) ([]map[string]interface{}, error) {
	var unexpandedValMaps []map[string]interface{}
	switch p := params.(type) {
	case map[string]interface{}:
		unexpandedValMaps = []map[string]interface{}{p}
	case []map[string]interface{}:
		unexpandedValMaps = p
	default:
		return nil, errors.Errorf("invalid parameters: expected object or list of objects got %T", p)
	}

	var expandedValMaps []map[string]interface{}
	for _, m := range unexpandedValMaps {
		mExpanded, err := expandValMap(m)
		if err != nil {
			return nil, errors.Wrap(err, "expanding parameter map")
		}
		expandedValMaps = append(expandedValMaps, mExpanded...)
	}
	return expandedValMaps, nil
}

func expandValMap(m map[string]interface{}) ([]map[string]interface{}, error) {
	type paramToExpand struct {
		name string
		values []interface{}
	}
	var paramsToExpand []paramToExpand

	numExpanded := int64(1)
	for k, v := range m {
		if !stringutils.ConsumeSuffix(&k, "*") {
			continue
		}
		vals, ok := v.([]interface{})
		if !ok || len(vals) == 0 {
			return nil, errors.Errorf("parameter %s is supposed to expanded, but its value is not a list (or an empty list)", k)
		}
		paramsToExpand = append(paramsToExpand, paramToExpand{
			name: k,
			values: vals,
		})
		numExpanded *= int64(len(vals))
	}
	if len(paramsToExpand) == 0 {
		return []map[string]interface{}{m}, nil
	}

	baseMap := make(map[string]interface{}, len(m))
	for k, v := range m {
		if !strings.HasSuffix(k, "*") {
			baseMap[k] = v
		}
	}
	for _, p := range paramsToExpand {
		if _, ok := baseMap[p.name]; ok {
			return nil, errors.Errorf("parameter %s is supposed to be expanded, but parameters map already contains an entry for it in unexpanded form", p.name)
		}
	}

	allExpanded := make([]map[string]interface{}, 0, numExpanded)

	for i := int64(0); i < numExpanded; i++ {
		expanded := make(map[string]interface{}, len(m))
		for k, v := range baseMap {
			expanded[k] = v
		}

		idx := i
		for _, p := range paramsToExpand {
			expanded[p.name] = p.values[idx % int64(len(p.values))]
			idx /= int64(len(p.values))
		}

		allExpanded = append(allExpanded, expanded)
	}
	return allExpanded, nil
}