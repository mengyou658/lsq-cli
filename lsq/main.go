package main

import (
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/iAsuma/lsq-cli/internal/cmd"
	_ "github.com/iAsuma/lsq-cli/internal/packed"
	"github.com/iAsuma/lsq-cli/utility/mlog"
)

func main() {
	defer func() {
		if exception := recover(); exception != nil {
			if err, ok := exception.(error); ok {
				mlog.Print(err.Error())
			} else {
				panic(exception)
			}
		}
	}()

	var (
		ctx = gctx.New()
	)

	command, err := gcmd.NewFromObject(cmd.LSQ)
	if err != nil {
		panic(err)
	}

	err = command.AddObject(
		cmd.Gen,
		cmd.Version,
		cmd.Dec,
	)
	if err != nil {
		panic(err)
	}
	command.Run(ctx)
}
