package shaders

func UpdatePulse(params map[string]any) map[string]any {
	counter := params["Counter"].(int)
	counter++
	if counter > 40 {
		counter = 0
	}
	params["Counter"] = counter
	return params
}
