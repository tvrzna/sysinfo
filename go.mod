module github.com/tvrzna/sysinfo

go 1.21.0

require github.com/tvrzna/go-utils/args v0.0.2

require github.com/tvrzna/pkgtray/checker v0.1.1

replace github.com/tvrzna/sysinfo/metric => ./metric

replace github.com/tvrzna/sysinfo/metric/smartctl => ./smartctl
