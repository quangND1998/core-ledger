package wv

import "github.com/elliotchance/orderedmap/v3"

func OrderedMapToMap(om *orderedmap.OrderedMap[string, any]) map[string]interface{} {
	result := make(map[string]interface{})
	for el := om.Front(); el != nil; el = el.Next() {
		// Nếu value lại là OrderedMap thì đệ quy
		switch v := el.Value.(type) {
		case *orderedmap.OrderedMap[string, any]:
			result[el.Key] = OrderedMapToMap(v)
		default:
			result[el.Key] = v
		}
	}
	return result
}
