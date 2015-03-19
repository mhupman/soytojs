package main

import (
	"fmt"
	"io"
	"os"

	_ "yext/pages/publishing/templating" // Need this to get access to the soyhtml.Directives in urlscheme.go

	"github.com/robfig/soy"
	"github.com/robfig/soy/ast"
	"github.com/robfig/soy/soyjs"
	"github.com/robfig/soy/template"
)

var (
	templateDirectory = os.Args[1]
)

func init() {
	InjectSoyJSFuncs()
}

func main() {
	fmt.Println(templateDirectory)

	fmt.Println("Loading soy templates from: ", templateDirectory)

	registry, err := soy.NewBundle().
		AddTemplateDir(templateDirectory).
		Compile()

	if err != nil {
		panic(err)
	}

	f, err := os.Create("out.js")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = SoyToJS(registry, f)

	if err != nil {
		panic(err)
	}
}

func SoyToJS(r *template.Registry, out io.Writer) error {
	g := soyjs.NewGenerator(r)

	for _, soyfile := range r.SoyFiles {
		fmt.Println(soyfile.Name)
		if err := g.WriteFile(out, soyfile.Name); err != nil {
			return err
		}
	}

	return nil
}

func InjectSoyJSFuncs() {
	soyjs.Funcs["collapseDays"] = externalJSFunc("collapseDays", 1)
	soyjs.Funcs["timef"] = externalJSFunc("timef", 2)
	soyjs.Funcs["timestampf"] = externalJSFunc("timestampf", 2)
	soyjs.Funcs["ltimef"] = externalJSFunc("ltimef", 3)
	soyjs.Funcs["ltimestampf"] = externalJSFunc("ltimestampf", 3)
	soyjs.Funcs["lnumberf"] = externalJSFunc("lnumberf", 2)
	soyjs.Funcs["lpercentf"] = externalJSFunc("lpercentf", 2)
	soyjs.Funcs["lcurrencyf"] = externalJSFunc("lcurrencyf", 3)
	soyjs.Funcs["llanguageName"] = externalJSFunc("llanguageName", 2)
	soyjs.Funcs["lcountryName"] = externalJSFunc("lcountryName", 2)
	soyjs.Funcs["llocaleName"] = externalJSFunc("llocaleName", 2)
	soyjs.Funcs["lregionName"] = externalJSFunc("lregionName", 3)

	soyjs.Funcs["substring"] = soyjs.Func{FuncSubstring, []int{2, 3}}
	soyjs.Funcs["sameDay"] = externalJSFunc("sameDay", 2)
	soyjs.Funcs["sameYear"] = externalJSFunc("sameYear", 2)
	soyjs.Funcs["listItems"] = externalJSFunc("listItems", 1)
	soyjs.Funcs["strlen"] = soyjs.Func{FuncStrlen, []int{1}}
	soyjs.Funcs["stripOutPhoneDigits"] = externalJSFunc("stripOutPhoneDigits", 1)
	soyjs.Funcs["gmap"] = serverSideOnlyWarning("gmap")
	soyjs.Funcs["fullState"] = externalJSFunc("fullState", 1)
	soyjs.Funcs["prettyPrintPhone"] = externalJSFunc("prettyPrintPhone", 1)
	soyjs.Funcs["augmentList"] = externalJSFunc("augmentList", 2)
	soyjs.Funcs["sortList"] = externalJSFunc("sortList", 2)
	soyjs.Funcs["sortListByKeys"] = externalJSFunc("sortListByKeys", 2)
	soyjs.Funcs["groupListByKey"] = externalJSFunc("groupListByKey", 2)
	soyjs.Funcs["replace"] = soyjs.Func{FuncReplace, []int{3, 4}}
	soyjs.Funcs["slice"] = soyjs.Func{FuncSlice, []int{2, 3}}
}

const jsFuncNamespace = "yext.pages.soy."

func externalJSFunc(name string, maxArgs int) soyjs.Func {
	apply := func(js soyjs.JSWriter, args []ast.Node) {
		js.Write(jsFuncNamespace, name, "(", args[0])

		for i := 1; i < len(args); i++ {
			js.Write(", ", args[i])
		}

		js.Write(")")
	}

	return soyjs.Func{apply, []int{maxArgs}}
}

func serverSideOnlyWarning(name string) soyjs.Func {
	apply := func(js soyjs.JSWriter, args []ast.Node) {
		js.Write("alert(", name, "+ \" is not implemented for JS templates\")")
	}

	return soyjs.Func{apply, []int{0}}
}

func FuncSubstring(js soyjs.JSWriter, args []ast.Node) {
	switch len(args) {
	case 1:
		js.Write(args[0], ".sub(", args[1], ")")
	default:
		js.Write(args[0], ".substring(", args[1], ", ", args[2], ")")
	}
}

func FuncStrlen(js soyjs.JSWriter, args []ast.Node) {
	js.Write(args[0], ".length")
}

func FuncReplace(js soyjs.JSWriter, args []ast.Node) {
	js.Write(args[0], ".replace(", args[1], ",", args[2], ")")
}

func FuncSlice(js soyjs.JSWriter, args []ast.Node) {
	js.Write(args[0], ".slice(", args[1], ")")

	if len(args) > 2 {
		js.Write(", ", args[2])
	}

	js.Write(")")
}
