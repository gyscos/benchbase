package main

import "github.com/Gyscos/benchbase"

// Dispatch takes a flat list of benchmarks, and dispatch them according to
// their value of the given spec.
func Dispatch(data []*benchbase.Benchmark, spec string, values []string) [][]*benchbase.Benchmark {
	// Number of categories to compare
	n := len(values)
	result := make([][]*benchbase.Benchmark, n)

	for _, b := range data {
		for i, v := range values {
			if b.Conf[spec] == v {
				result[i] = append(result[i], b)
				break
			}
		}
	}

	return result
}

// Project takes a list of benchmark sorted by spec value, and regroup them by
// common configuration (except for spec)
func Project(data [][]*benchbase.Benchmark, spec string) [][][]*benchbase.Benchmark {
	// Number of categories to compare
	n := len(data)

	// Prepare the map list
	mapList := make([]map[string][]*benchbase.Benchmark, n)
	for i := 0; i < n; i++ {
		mapList[i] = make(map[string][]*benchbase.Benchmark)
	}

	for i, l := range data {
		for _, b := range l {
			key := b.Conf.PartialHash(spec)
			mapList[i][key] = append(mapList[i][key], b)
		}
	}

	// Invert the map and keep populated entried
	m := Invert(mapList)

	// Now just flatten the map
	return Flatten(m)
}

func Invert(mapList []map[string][]*benchbase.Benchmark) map[string][][]*benchbase.Benchmark {
	n := len(mapList)

	result := make(map[string][][]*benchbase.Benchmark)
	population := make(map[string]int)

	for i, m := range mapList {
		// Category #i has its sorted list of benchmarks
		for key, l := range m {
			if result[key] == nil {
				result[key] = make([][]*benchbase.Benchmark, n)
			}
			result[key][i] = l
			population[key]++
		}
	}

	// Only keep populated entries
	for key, p := range population {
		if p < 2 {
			delete(result, key)
		}
	}

	return result
}

func Flatten(m map[string][][]*benchbase.Benchmark) [][][]*benchbase.Benchmark {
	var result [][][]*benchbase.Benchmark

	for _, l := range m {
		result = append(result, l)
	}

	return result
}