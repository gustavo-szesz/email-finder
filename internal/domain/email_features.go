package domain

type DisposableResult struct {
	IsDisposable bool   `json:"is_disposable"`
	Provider     string `json:"provider,omitempty"`
}

type TypoResult struct {
	Suspicious bool   `json:"suspicious"`
	Target     string `json:"target,omitempty"`
}

type CatchAllResult struct {
	IsCatchAll bool   `json:"is_catch_all"`
	TestEmail  string `json:"test_email"`
}
