package cmd

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/iAsuma/lsq-cli/internal/consts"
	"github.com/iAsuma/lsq-cli/utility/mlog"
	"github.com/iAsuma/lsq-cli/utility/utils"
)

const (
	defaultApiPath        = `api/v1`
	defaultControllerPath = `internal/controller`
	defaultModelPath      = `internal/model`
	defaultLogicPath      = `internal/logic`
	defaultLastLogicPath  = `internal/service/logic`
)

const (
	defaultApiSuffix        = `_api`
	defaultControllerSuffix = ``
	defaultModelSuffix      = `_model`
	defaultLogicSuffix      = `_logic`
)

const (
	cGenMcUsage = `lsq gen mc [OPTION]`
	cGenMcBrief = `automatically generate GoFrameV2 files for api/controller/logic/model`
	cGenMcEg    = ``
	cGenMcAd    = ``

	cGenMcArgName       = `required, the name for the business file`
	cGenMcArgApi        = `generate api entrance struct file`
	cGenMcArgController = `generate controller file`
	cGenMcArgLogic      = `generate logic file`
	cGenMcArgModel      = `generate model file`
	cGenMcArgOverwrite  = `overwrite the file, danger!`
	cGenMcArgLast       = `generate business file for goframe 2.0 (< 2.1) `

	tplControllerCamelCase = `{TplControllerCamelCase}`
	tplLogicCamelCase      = `{TplLogicCamelCase}`
)

func init() {
	gtag.Sets(g.MapStrStr{
		`cGenMcUsage`:         cGenMcUsage,
		`cGenMcBrief`:         cGenMcBrief,
		`cGenMcEg`:            cGenMcEg,
		`cGenMcAd`:            cGenMcAd,
		`cGenMcArgName`:       cGenMcArgName,
		`cGenMcArgApi`:        cGenMcArgApi,
		`cGenMcArgController`: cGenMcArgController,
		`cGenMcArgLogic`:      cGenMcArgLogic,
		`cGenMcArgModel`:      cGenMcArgModel,
		`cGenMcArgOverwrite`:  cGenMcArgOverwrite,
		`cGenMcArgLast`:       cGenMcArgLast,
	})
}

type (
	cGenMcInput struct {
		g.Meta     `name:"mc" usage:"{cGenMcUsage}" brief:"{cGenMcBrief}" eg:"{cGenMcEg}" ad:"{cGenMcAd}"`
		Name       string `name:"name"       short:"n" brief:"{cGenMcArgName}"`
		Api        bool   `name:"api"        short:"a" brief:"{cGenMcArgApi}" orphan:"true"`
		Controller bool   `name:"controller" short:"c" brief:"{cGenMcArgController}" orphan:"true"`
		Logic      bool   `name:"logic"      short:"l" brief:"{cGenMcArgLogic}" orphan:"true"`
		Model      bool   `name:"model"      short:"m" brief:"{cGenMcArgModel}" orphan:"true"`
		Overwrite  bool   `name:"overwrite"  short:"o" brief:"{cGenMcArgOverwrite}" orphan:"true"`
		Last       bool   `name:"last"       short:"lt" brief:"{cGenMcArgLast}"  orphan:"true"`
	}
	cGenMcOutput struct{}
)

func (c cGen) Mc(ctx context.Context, in cGenMcInput) (out *cGenMcOutput, err error) {
	if in.Name == "" {
		mlog.Fatalf("name is empty")
		return
	}

	var (
		nameCamelCase = gstr.CaseCamel(in.Name)
	)

	genAll := true
	genApi := false
	genController := false
	genLogic := false
	genModel := false

	if in.Api {
		genAll = false
		genApi = true
	}

	if in.Controller {
		genAll = false
		genController = true
	}

	if in.Logic {
		genAll = false
		genLogic = true
	}

	if in.Model {
		genAll = false
		genModel = true
	}

	if genAll || genApi {
		generateApiFile(in)
	}

	if genAll || genController {
		generateControllerFile(in, nameCamelCase)
	}

	if genAll || genLogic {
		generateLogicFile(in, nameCamelCase)
	}

	if genAll || genModel {
		generateModelFile(in)
	}

	return
}

func generateApiFile(in cGenMcInput) {
	fileName := fmt.Sprintf("%s%s.go", in.Name, defaultApiSuffix)
	path := gfile.Join(defaultApiPath, fileName)

	if in.Overwrite || !gfile.Exists(path) {
		if err := gfile.PutContents(path, getApiFileContent()); err != nil {
			mlog.Fatalf("write api file get something wrong", path, err)
		} else {
			utils.GoFmt(path)
			mlog.Print("generated:", path)
		}
	}
}

func generateControllerFile(in cGenMcInput, nameCamelCase string) {
	fileName := fmt.Sprintf("%s%s.go", in.Name, defaultControllerSuffix)
	path := gfile.Join(defaultControllerPath, fileName)

	if in.Overwrite || !gfile.Exists(path) {
		if err := gfile.PutContents(path, getControllerContent(nameCamelCase)); err != nil {
			mlog.Fatalf("write controller file get something wrong", path, err)
		} else {
			utils.GoFmt(path)
			mlog.Print("generated:", path)
		}
	}
}

func generateLogicFile(in cGenMcInput, nameCamelCase string) {
	fileName := fmt.Sprintf("%s%s.go", in.Name, defaultLogicSuffix)
	var path string
	if in.Last {
		path = gfile.Join(defaultLastLogicPath, fileName)
	} else {
		path = gfile.Join(defaultLogicPath, fileName)
	}

	if in.Overwrite || !gfile.Exists(path) {
		if err := gfile.PutContents(path, getLogicContent(nameCamelCase)); err != nil {
			mlog.Fatalf("write logic file get something wrong", path, err)
		} else {
			utils.GoFmt(path)
			mlog.Print("generated:", path)
		}
	}
}

func generateModelFile(in cGenMcInput) {
	fileName := fmt.Sprintf("%s%s.go", in.Name, defaultModelSuffix)
	path := gfile.Join(defaultModelPath, fileName)

	if in.Overwrite || !gfile.Exists(path) {
		if err := gfile.PutContents(path, getModelFileContent()); err != nil {
			mlog.Fatalf("write api file get something wrong", path, err)
		} else {
			utils.GoFmt(path)
			mlog.Print("generated:", path)
		}
	}
}

func getApiFileContent() string {
	return consts.TemplateApiContent
}

func getControllerContent(nameCamelCase string) string {
	controllerContent := gstr.ReplaceByMap(consts.TemplateControllerContent, g.MapStrStr{
		tplControllerCamelCase: nameCamelCase,
	})

	return controllerContent
}

func getLogicContent(nameCamelCase string) string {
	logicContent := gstr.ReplaceByMap(consts.TemplateLogicContent, g.MapStrStr{
		tplLogicCamelCase: nameCamelCase,
	})

	return logicContent
}

func getModelFileContent() string {
	return consts.TemplateModelContent
}
