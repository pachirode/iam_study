package pumps

import "github.com/pachirode/iam_study/internal/pump/analytics"

func (p *CommonPumpConfig) SetFilters(filters analytics.AnalyticsFilters) {
	p.filters = filters
}

func (p *CommonPumpConfig) GetFilters() analytics.AnalyticsFilters {
	return p.filters
}

func (p *CommonPumpConfig) SetTimeout(timeout int) {
	p.timeout = timeout
}

func (p *CommonPumpConfig) GetTimeout() int {
	return p.timeout
}

func (p *CommonPumpConfig) SetOmitDetailedRecording(omitDetailedRecording bool) {
	p.OmitDetailedRecording = omitDetailedRecording
}

func (p *CommonPumpConfig) GetOmitDetailedRecording() bool {
	return p.OmitDetailedRecording
}
