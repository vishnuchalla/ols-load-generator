package main

// Type to store the test config.
type TestConfig struct {
	Rps        int      `json:"rps"`
	Host       string   `json:"host"`
	HitSize    int      `json:"hitsize"`
	AuthToken  string   `json:"-"`
	Uuid       string   `json:"uuid"`
	ESHost     string   `json:"eshost"`
	ESIndex    string   `json:"esindex"`
	MetricStep int      `json:"metricstep"`
	Profiles   []string `json:"profiles"`
}
