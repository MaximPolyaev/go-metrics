package server

const (
	updateAction = "/update/"

	updateGaugeAction   = updateAction + metricsGaugeType + "/"
	updateCounterAction = updateAction + metricsCounterType + "/"
)
