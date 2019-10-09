package metricSource

type setHelper map[string]bool

func newSetHelperFromTriggerPatternMetrics(metrics TriggerPatternMetrics) setHelper {
	result := make(setHelper, len(metrics))
	for metricName := range metrics {
		result[metricName] = true
	}
	return result
}

func (h setHelper) difference(other setHelper) setHelper {
	result := make(setHelper, len(h))
	for metricName := range h {
		result[metricName] = true
	}
	for metricName := range other {
		if _, ok := result[metricName]; ok {
			delete(result, metricName)
		}
	}
	return result
}

func (h setHelper) intersection(other setHelper) setHelper {
	result := make(setHelper, len(h))
	for metricName := range h {
		if _, ok := other[metricName]; ok {
			result[metricName] = true
		}
	}
	return result
}

func isOneMetricMap(metrics map[string]MetricData) (bool, string) {
	if len(metrics) == 1 {
		for metricName := range metrics {
			return true, metricName
		}
	}
	return false, ""
}

// TriggerPatternMetrics ...
// TODO(litleleprikon): fill
type TriggerPatternMetrics map[string]MetricData

// newTriggerPatternMetricsWithCapacity ...
// TODO(litleleprikon): fill
func newTriggerPatternMetricsWithCapacity(capacity int) TriggerPatternMetrics {
	return make(TriggerPatternMetrics, capacity)
}

// NewTriggerPatternMetrics ..
// TODO(litleleprikon): fill
func NewTriggerPatternMetrics(source FetchedPatternMetrics) TriggerPatternMetrics {
	result := newTriggerPatternMetricsWithCapacity(len(source))
	for _, m := range source {
		result[m.Name] = m
	}
	return result
}

func (m TriggerPatternMetrics) filterByNames(names setHelper) TriggerPatternMetrics {
	result := newTriggerPatternMetricsWithCapacity(len(names))
	for name := range names {
		result[name] = m[name]
	}
	return result
}

// TriggerMetrics ...
// TODO(litleleprikon): fill
type TriggerMetrics map[string]TriggerPatternMetrics

// NewTriggerMetricsWithCapacity is a constructor function that creates TriggerMetrics with initialized empty fields
func NewTriggerMetricsWithCapacity(capacity int) TriggerMetrics {
	return make(TriggerMetrics, capacity)
}

func (m TriggerMetrics) FilterAloneMetrics() (TriggerMetrics, MetricsToCheck) {
	result := NewTriggerMetricsWithCapacity(len(m))
	aloneMetrics := make(MetricsToCheck)

	for targetName, patternMetrics := range m {
		if oneMetricMap, metricName := isOneMetricMap(patternMetrics); oneMetricMap {
			aloneMetrics[targetName] = patternMetrics[metricName]
			continue
		}
		result[targetName] = m[targetName]
	}
	return result, aloneMetrics
}

// CrossIntersection ...
// TODO(litleleprikon): fill
func (m TriggerMetrics) CrossIntersection() TriggerMetrics {
	result := NewTriggerMetricsWithCapacity(len(m))

	if len(m) == 0 {
		return result
	}

	multiMetricsTarget, _ := m.multiMetricsTarget()
	commonMetrics := newSetHelperFromTriggerPatternMetrics(m[multiMetricsTarget])

	for _, patternMetrics := range m {
		currentMetrics := newSetHelperFromTriggerPatternMetrics(patternMetrics)
		commonMetrics = commonMetrics.intersection(currentMetrics)
	}

	for targetName, patternMetrics := range m {
		result[targetName] = patternMetrics.filterByNames(commonMetrics)
	}
	return result
}

// multiMetricsTarget is a function that finds any first target with
// amount of metrics greater than one and returns set with names of this metrics.
func (m TriggerMetrics) multiMetricsTarget() (string, setHelper) {
	commonMetrics := make(setHelper)
	for targetName, metrics := range m {
		if len(metrics) > 1 {
			for metricName := range metrics {
				commonMetrics[metricName] = true
			}
			return targetName, commonMetrics
		}
	}
	return "", nil
}

// ConvertForCheck is a function that converts TriggerMetrics with structure
// map[TargetName]map[MetricName]MetricData to ConvertedTriggerMetrics
// with structure map[MetricName]map[TargetName]MetricData and fill with
// duplicated metrics targets that have only one metric. Second return value is
// a map with names of targets that had only one metric as key and original metric name as value.
func (m TriggerMetrics) ConvertForCheck() TriggerMetricsToCheck {
	result := make(TriggerMetricsToCheck)
	_, commonMetrics := m.multiMetricsTarget()

	hasAtLeastOneMultiMetricsTarget := commonMetrics != nil

	if !hasAtLeastOneMultiMetricsTarget && len(m) <= 1 {
		return result
	}

	for targetName, targetMetrics := range m {
		oneMetricTarget, oneMetricName := isOneMetricMap(targetMetrics)

		for metricName := range commonMetrics {
			if _, ok := result[metricName]; !ok {
				result[metricName] = make(MetricsToCheck, len(m))
			}

			if oneMetricTarget {
				result[metricName][targetName] = m[targetName][oneMetricName]
				continue
			}

			result[metricName][targetName] = m[targetName][metricName]
		}
	}
	return result
}

// MetricsToCheck ...
// TODO(litleleprikon): fill
type MetricsToCheck map[string]MetricData

// MetricName ...
// TODO(litleleprikon): fill
func (m MetricsToCheck) MetricName() string {
	for _, metric := range m {
		return metric.Name
	}
	return ""
}

// Merge ...
// TODO(litleleprikon): fill
func (m MetricsToCheck) Merge(other MetricsToCheck) MetricsToCheck {
	result := make(MetricsToCheck, len(m)+len(other))
	for k, v := range m {
		result[k] = v
	}
	for k, v := range other {
		result[k] = v
	}
	return result
}

// IsUbiquitousMetrics ...
// TODO(litleleprikon): fill
func (m MetricsToCheck) IsUbiquitousMetrics() bool {
	var commonMetric string
	firstMetric := true
	for _, metric := range m {
		if firstMetric {
			commonMetric = metric.Name
		}
		if metric.Name != commonMetric {
			return false
		}
	}
	return true
}

// TriggerMetricsToCheck ...
// TODO(litleleprikon): fill
type TriggerMetricsToCheck map[string]MetricsToCheck
