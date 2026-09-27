package uws1

import (
	"bytes"
	"encoding/json"
)

func cloneExecutionRecord(record ExecutionRecord) ExecutionRecord {
	cloned := record
	cloned.Result = cloneExecutionResult(record.Result)
	if len(record.Outputs) > 0 {
		cloned.Outputs = make(map[string]any, len(record.Outputs))
		for key, value := range record.Outputs {
			cloned.Outputs[key] = value
		}
	}
	return cloned
}

func cloneExecutionResult(result any) any {
	switch typed := result.(type) {
	case json.RawMessage:
		return json.RawMessage(bytes.Clone(typed))
	case map[string]any:
		cloned := make(map[string]any, len(typed))
		for key, value := range typed {
			cloned[key] = cloneExecutionResult(value)
		}
		return cloned
	case []any:
		cloned := make([]any, len(typed))
		for index, value := range typed {
			cloned[index] = cloneExecutionResult(value)
		}
		return cloned
	default:
		return result
	}
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
