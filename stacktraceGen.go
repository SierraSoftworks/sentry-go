package sentry

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/pkg/errors"
)

type stackTracer interface {
	StackTrace() errors.StackTrace
}

type stackTraceFrames []*stackTraceFrame

func (c stackTraceFrames) Len() int      { return len(c) }
func (c stackTraceFrames) Swap(i, j int) { c[j], c[i] = c[i], c[j] }
func (c stackTraceFrames) Reverse() {
	for i, j := 0, c.Len()-1; i < j; i, j = i+1, j-1 {
		c.Swap(i, j)
	}
}

func getStacktraceFramesForError(err error) stackTraceFrames {
	if err, ok := err.(stackTracer); ok {
		frames := stackTraceFrames{}
		for _, f := range err.StackTrace() {
			pc := uintptr(f) - 1
			frame := getStacktraceFrame(pc)
			if frame != nil {
				frames = append(frames, frame)
			}
		}

		frames.Reverse()
		return frames
	}

	return stackTraceFrames{}
}

func getStacktraceFrames(skip int) stackTraceFrames {
	pcs := make([]uintptr, 30)
	if c := runtime.Callers(skip+2, pcs); c > 0 {
		frames := stackTraceFrames{}
		for _, pc := range pcs {
			frame := getStacktraceFrame(pc)
			if frame != nil {
				frames = append(frames, frame)
			}
		}

		frames.Reverse()
		return frames
	}

	return stackTraceFrames{}
}

func getStacktraceFrame(pc uintptr) *stackTraceFrame {
	frame := &stackTraceFrame{}

	if fn := runtime.FuncForPC(pc); fn != nil {
		frame.AbsoluteFilename, frame.Line = fn.FileLine(pc)
		frame.Package, frame.Module, frame.Function = formatFuncName(fn.Name())
		frame.Filename = shortFilename(frame.AbsoluteFilename, frame.Package)
	} else {
		frame.AbsoluteFilename = "unknown"
		frame.Filename = "unknown"
	}

	return frame
}

func (f *stackTraceFrame) ClassifyInternal(internalPrefixes []string) {
	if f.Module == "main" {
		f.InApp = true
		return
	}

	for _, prefix := range internalPrefixes {
		if strings.HasPrefix(f.Package, prefix) && !strings.Contains(f.Package, "vendor") {
			f.InApp = true
			return
		}
	}
}

// formatFuncName converts a stack frame function name, which is commonly of the form
// 'github.com/SierraSoftworks/sentry-go/v2.TestStackTraceGenerator.func3', into a well-formed
// package, module, and function name.
// For the above example, this will result in the following values being emitted:
//   - Package: github.com/SierraSoftworks/sentry-go/v2
//   - Module: SierraSoftworks/sentry-go/v2
//   - Function: TestStackTraceGenerator.func3
func formatFuncName(fnName string) (pack, module, name string) {
	name = fnName
	pack = ""
	module = ""

	name = strings.Replace(name, "·", ".", -1)

	packageParts := strings.Split(name, "/")
	codeParts := strings.Split(packageParts[len(packageParts)-1], ".")

	pack = strings.Join(append(packageParts[:len(packageParts)-1], codeParts[0]), "/")
	name = strings.Join(codeParts[1:], ".")

	if len(packageParts) > 2 {
		module = strings.Join(append(packageParts[2:len(packageParts)-1], codeParts[0]), "/")
	} else {
		module = codeParts[0]
	}

	return
}

func shortFilename(absFile, pkg string) string {
	if pkg == "" {
		return absFile
	}

	if idx := strings.Index(absFile, fmt.Sprintf("%s/", pkg)); idx != -1 {
		return absFile[idx:]
	}

	return absFile
}
