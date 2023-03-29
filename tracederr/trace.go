package tracederr

import (
	"fmt"
	"runtime"
	"strings"
	"sync"
)

var pool sync.Pool

// default value for how many stacks on each error.
var defaultDeep = 100

// skip = 1, it is means to skip the first caller, usually main.go.
var defaultSkip = 1

type trace struct {
	fileName    string
	lineNum     int
	packageName string
	funcName    string
}

func (t *trace) String() string {
	str := fmt.Sprintf("%s:%d", t.fileName, t.lineNum)
	if t.packageName == "" && t.funcName == "" {
		return str
	}

	return str + fmt.Sprintf(" [%s](%s)", t.packageName, t.funcName)
}

func (t *trace) File() string {
	return t.fileName
}

func (t *trace) Line() int {
	return t.lineNum
}

func (t *trace) Func() string {
	return t.funcName
}

// StackTrace is a function to backtrack traces
// which let user to set how deep of trace and how many skip to reach all of those trace
func StackTrace(skip, deep int) []trace {
	if deep <= 0 {
		deep = defaultDeep
	}

	if skip < 0 {
		skip = defaultSkip
	}
	skip += 2

	dataLen := deep + 10
	var data []uintptr
	if tmp1 := pool.Get(); tmp1 != nil {
		switch tmp2 := tmp1.(type) {
		case []uintptr:
			if len(tmp2) >= dataLen {
				data = tmp2
			}
		}
	}

	if data == nil {
		data = make([]uintptr, dataLen)
	}

	pc := data[:dataLen]
	pc = pc[:runtime.Callers(skip, pc)]
	if len(pc) == 0 {
		return nil
	}

	traces := make([]trace, 0, len(pc))

	frames := runtime.CallersFrames(pc)
	for {
		var pkgName, fnName string
		frame, more := frames.Next()

		// as written on docs, frame.Function is package path-qualified function name
		pkgName = frame.Function
		// If Func is not nil then Function == Func.Name().
		fnName = frame.Function

		// replace them with function which split pkg and fn
		if frame.Func != nil {
			pkgName, fnName = pkgNameAndFnName(frame.Func)
		}

		// we don't need stack from runtime and other repos
		if frame.Line != 0 && frame.File != "" &&
			!inPkg(pkgName, "runtime") &&
			!inPkg(pkgName, "net/http") {
			traces = append(traces, trace{
				fileName:    frame.File,
				lineNum:     frame.Line,
				packageName: pkgName,
				funcName:    fnName,
			})
		}
		if len(traces) == deep {
			break
		}
		if !more {
			break
		}
	}

	pool.Put(&data)
	return traces
}

func pkgNameAndFnName(fn *runtime.Func) (string, string) {
	name := fn.Name()
	pkg := ""

	// The name includes the path name to the package, which is unnecessary
	// since the file name is already included.  Plus, it has center dots.
	// That is, we see
	//  runtime/debug.*T·ptrmethod
	// and want
	//  *T.ptrmethod
	// Since the package path might contains dots (e.g. code.google.com/...),
	// we first remove the path prefix if there is one.
	if lastslash := strings.LastIndex(name, "/"); lastslash >= 0 {
		pkg += name[:lastslash] + "/"
		name = name[lastslash+1:]
	}
	if period := strings.Index(name, "."); period >= 0 {
		pkg += name[:period]
		name = name[period+1:]
	}

	name = strings.Replace(name, "·", ".", -1)
	return pkg, name
}

func inPkg(what, pkg string) bool {
	if !strings.HasPrefix(what, pkg) {
		return false
	}

	if len(what) == len(pkg) {
		return true
	}

	if len(what) > len(pkg) {
		return what[len(pkg)] == '.' || what[len(pkg)] == '/'
	}

	return false
}
