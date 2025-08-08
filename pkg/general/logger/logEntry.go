package logger

type LogEntry struct {
	Timestamp    string                 `json:"timestamp"`
	Level        string                 `json:"level"`
	Service      string                 `json:"service"`
	Env          string                 `json:"env"`
	Message      string                 `json:"message"`
	RequestBody  string                 `json:"requestBody"`
	ResponseBody string                 `json:"responseBody"`
	RequestID    string                 `json:"requestId,omitempty"`
	TraceID      string                 `json:"traceId,omitempty"`
	SpanID       string                 `json:"spanId,omitempty"`
	Method       string                 `json:"method,omitempty"`
	Path         string                 `json:"path,omitempty"`
	Error        error                  `json:"error,omitempty"`
	Status       int                    `json:"status,omitempty"`
	LatencyMs    int64                  `json:"latencyMs,omitempty"`
	Extra        map[string]interface{} `json:"extra,omitempty"`
}
