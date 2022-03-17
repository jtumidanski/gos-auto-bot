package database

type Model struct {
	Levy         Task `json:"levy"`
	Divination   Task `json:"divination"`
	Council      Task `json:"council"`
	Academy      Task `json:"academy"`
	Rankings     Task `json:"rankings"`
	Harem        Task `json:"harem"`
	Coalition    Task `json:"coalition"`
	Ads          Task `json:"ads"`
	Union        Task `json:"union"`
	DailyCheckIn Task `json:"dailyCheckIn"`
}

type Task struct {
	Execution string `json:"execution"`
	Count     int    `json:"count"`
}
