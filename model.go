package main

type SearchParam struct {
	Newest   int
	Keyword  string
	MinPrice string
	Matchid  int
	SortBy   string
}

type workerDoCrawlReturn struct {
	Result  [][]string
	MatchID int
	Dur     string
	Err     error
}

type BasicResp struct {
	Msg     string
	Data    interface{}
	Latency string
}
