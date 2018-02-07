package main

// #include <stdlib.h>
// #include <Python.h>
//
// extern int PyDockerfile_PyArg_ParseTuple_U(PyObject*, PyObject**);
// extern PyObject* PyDockerfile_Py_None;
//
// extern PyObject* PyDockerfile_GoIOError;
// extern PyObject* PyDockerfile_GoParseError;
// extern PyObject* PyDockerfile_NewCommand(
//     PyObject*, PyObject*, PyObject*, PyObject*, PyObject*, PyObject*,
//     PyObject*
// );
import "C"
import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/asottile/dockerfile"
)

func raise(err error) *C.PyObject {
	var tp *C.PyObject
	switch err.(type) {
	case dockerfile.IOError:
		tp = C.PyDockerfile_GoIOError
	case dockerfile.ParseError:
		tp = C.PyDockerfile_GoParseError
	default:
		panic(err)
	}
	cstr := C.CString(err.Error())
	C.PyErr_SetString(tp, cstr)
	C.free(unsafe.Pointer(cstr))
	return nil
}

func stringToPy(s string) *C.PyObject {
	cstr := C.CString(s)
	pystr := C.PyUnicode_FromString(cstr)
	C.free(unsafe.Pointer(cstr))
	return pystr
}

func stringToPyOrNone(s string) *C.PyObject {
	if s == "" {
		C.Py_IncRef(C.PyDockerfile_Py_None)
		return C.PyDockerfile_Py_None
	} else {
		return stringToPy(s)
	}
}

func sliceToTuple(strs []string) *C.PyObject {
	ret := C.PyTuple_New(C.Py_ssize_t(len(strs)))
	for i, str := range strs {
		pystr := stringToPy(str)
		if pystr == nil {
			C.Py_DecRef(ret)
			return nil
		}
		C.PyTuple_SetItem(ret, C.Py_ssize_t(i), pystr)
	}
	return ret
}

func boolToInt(b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}

func cmdsToPy(cmds []dockerfile.Command) *C.PyObject {
	var pyCmd, pySubCmd, pyJson, pyOriginal, pyStartLine, pyValue *C.PyObject
	var pyFlags *C.PyObject
	var ret *C.PyObject
	decrefAll := func() {
		C.Py_DecRef(pyCmd)
		C.Py_DecRef(pySubCmd)
		C.Py_DecRef(pyJson)
		C.Py_DecRef(pyOriginal)
		C.Py_DecRef(pyStartLine)
		C.Py_DecRef(pyFlags)
		C.Py_DecRef(pyValue)
		C.Py_DecRef(ret)
	}

	ret = C.PyTuple_New(C.Py_ssize_t(len(cmds)))
	for i, cmd := range cmds {
		pyCmd = stringToPyOrNone(cmd.Cmd)
		if pyCmd == nil {
			decrefAll()
			return nil
		}

		pySubCmd = stringToPyOrNone(cmd.SubCmd)
		if pySubCmd == nil {
			decrefAll()
			return nil
		}

		pyJson = C.PyBool_FromLong(C.long(boolToInt(cmd.Json)))

		pyOriginal = stringToPyOrNone(cmd.Original)
		if pyOriginal == nil {
			decrefAll()
			return nil
		}

		pyStartLine = C.PyLong_FromLong(C.long(cmd.StartLine))

		pyFlags = sliceToTuple(cmd.Flags)
		if pyFlags == nil {
			decrefAll()
			return nil
		}

		pyValue = sliceToTuple(cmd.Value)
		if pyValue == nil {
			decrefAll()
			return nil
		}

		pyCmd := C.PyDockerfile_NewCommand(
			pyCmd, pySubCmd, pyJson, pyOriginal, pyStartLine, pyFlags, pyValue,
		)
		C.PyTuple_SetItem(ret, C.Py_ssize_t(i), pyCmd)
	}
	return ret
}

func goStringFromArgs(args *C.PyObject) (string, error) {
	var obj *C.PyObject
	if C.PyDockerfile_PyArg_ParseTuple_U(args, &obj) == 0 {
		return "", fmt.Errorf("Failed to parse arguments")
	}
	bytes := C.PyUnicode_AsUTF8String(obj)
	ret := C.GoString(C.PyBytes_AsString(bytes))
	C.Py_DecRef(bytes)
	return ret, nil
}

//export all_cmds
func all_cmds(self *C.PyObject) *C.PyObject {
	return sliceToTuple(dockerfile.AllCmds())
}

//export parse_file
func parse_file(self *C.PyObject, args *C.PyObject) *C.PyObject {
	filename, err := goStringFromArgs(args)
	if err != nil {
		return nil
	}
	cmds, err := dockerfile.ParseFile(filename)
	if err != nil {
		return raise(err)
	}
	return cmdsToPy(cmds)
}

//export parse_string
func parse_string(self *C.PyObject, args *C.PyObject) *C.PyObject {
	s, err := goStringFromArgs(args)
	if err != nil {
		return nil
	}
	cmds, err := dockerfile.ParseReader(bytes.NewBufferString(s))
	if err != nil {
		return raise(err)
	}
	return cmdsToPy(cmds)
}

func main() {}
