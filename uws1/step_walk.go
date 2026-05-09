package uws1

import "fmt"

type stepTreeWalkHandlers struct {
	step     func(path string, step *Step) error
	nilStep  func(path string) error
	caseNode func(path string, c *Case) error
	nilCase  func(path string) error
}

func walkStepTree(path string, steps []*Step, h stepTreeWalkHandlers) error {
	for i, step := range steps {
		stepPath := fmt.Sprintf("%s[%d]", path, i)
		if step == nil {
			if h.nilStep != nil {
				if err := h.nilStep(stepPath); err != nil {
					return err
				}
			}
			continue
		}
		if h.step != nil {
			if err := h.step(stepPath, step); err != nil {
				return err
			}
		}
		if err := walkStepTree(stepPath+".steps", step.Steps, h); err != nil {
			return err
		}
		if err := walkCaseTree(stepPath+".cases", step.Cases, h); err != nil {
			return err
		}
		if err := walkStepTree(stepPath+".default", step.Default, h); err != nil {
			return err
		}
	}
	return nil
}

func walkCaseTree(path string, cases []*Case, h stepTreeWalkHandlers) error {
	for i, c := range cases {
		casePath := fmt.Sprintf("%s[%d]", path, i)
		if c == nil {
			if h.nilCase != nil {
				if err := h.nilCase(casePath); err != nil {
					return err
				}
			}
			continue
		}
		if h.caseNode != nil {
			if err := h.caseNode(casePath, c); err != nil {
				return err
			}
		}
		if err := walkStepTree(casePath+".steps", c.Steps, h); err != nil {
			return err
		}
	}
	return nil
}
