package drum

type JobSetting func(*RunStatus)

func SetJobRetryTime(t int) JobSetting {
	return func(rs *RunStatus) {
		rs.MaxTry = t
	}
}

func SetJobTryStep(s int) JobSetting {
	return func(rs *RunStatus) {
		rs.TryStep = s
	}
}

func SetJobMaxStep(s int) JobSetting {
	return func(rs *RunStatus) {
		rs.MaxStep = s
	}
}
