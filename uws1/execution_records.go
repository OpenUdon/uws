package uws1

func cloneExecutionRecord(record ExecutionRecord) ExecutionRecord {
	cloned := record
	if len(record.Outputs) > 0 {
		cloned.Outputs = make(map[string]any, len(record.Outputs))
		for key, value := range record.Outputs {
			cloned.Outputs[key] = value
		}
	}
	return cloned
}

func cloneExecutionRecords(records map[string]ExecutionRecord) map[string]ExecutionRecord {
	if records == nil {
		return nil
	}
	out := make(map[string]ExecutionRecord, len(records))
	for key, record := range records {
		out[key] = cloneExecutionRecord(record)
	}
	return out
}

func cloneCurrentExecution(current *CurrentExecutionContext) *CurrentExecutionContext {
	if current == nil {
		return nil
	}
	out := *current
	if len(current.Outputs) > 0 {
		out.Outputs = make(map[string]any, len(current.Outputs))
		for key, value := range current.Outputs {
			out.Outputs[key] = value
		}
	}
	return &out
}
