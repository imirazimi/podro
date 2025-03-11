package podroapp

type GetProvidersWeaklyReportReqeust struct {
}

type Report struct {
	Provider string  `json:"provider"`
	Average  float64 `json:"average"`
}
type GetProvidersWeaklyReportResponse struct {
	WeaklyReports []Report `json:"weakly_reports"`
}

type UpdateOrdersStatusReqeust struct {
}

type UpdateOrdersStatusResponse struct {
}

type CallProviderAPIRequest struct {
}

type CallProviderAPIResponse struct {
	Orders []struct {
		ID     uint   `json:"id"`
		Status string `json:"status"`
	} `json:"data"`
}
