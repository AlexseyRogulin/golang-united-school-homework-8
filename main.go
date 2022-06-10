package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

const (
	argOp       = "operation"
	argId       = "id"
	argItem     = "item"
	argFilename = "fileName"

	opList     = "list"
	opAdd      = "add"
	opFindById = "findById"
	opRemove   = "remove"
)

var (
	operationFl = flag.String("operation", "", "operation you want to execute")
	idFl        = flag.String("id", "", "user id")
	itemFl      = flag.String("item", "", "item")
	filenameFl  = flag.String("fileName", "", "file with users data")

	errOperationFlagNotSpecified = fmt.Errorf("-operation flag has to be specified")
	errFilenameFlagNotSpecified  = fmt.Errorf("-fileName flag has to be specified")
	errItemFlagNotSpecified      = fmt.Errorf("-item flag has to be specified")
	errIdFlagNotSpecidied        = fmt.Errorf("-id flag has to be specified")
)

func unknownOperationFlag(operation string) error {
	return fmt.Errorf("Operation %v not allowed!", operation)
}

func parseArgs() Arguments {
	flag.Parse()
	return Arguments{argOp: *operationFl, argId: fmt.Sprint(idFl), argItem: *itemFl, argFilename: *filenameFl}
}

func Perform(args Arguments, writer io.Writer) error {
	if args[argFilename] == "" {
		return errFilenameFlagNotSpecified
	}
	if args[argOp] == "" {
		return errOperationFlagNotSpecified
	}
	switch args[argOp] {
	case opList:
		return operationList(args, writer)
	case opAdd:
		return operationAdd(args, writer)
	case opFindById:
		return operationFindById(args, writer)
	case opRemove:
		return operationRemove(args, writer)
	default:
		return unknownOperationFlag(args[argOp])
	}
}

func operationList(args Arguments, writer io.Writer) error {
	data, err := getFileRawData(args[argFilename])
	if err != nil {
		return fmt.Errorf("restore data from %q: %w", args[argFilename], err)
	}
	writer.Write(data)
	return nil
}

func operationAdd(args Arguments, writer io.Writer) error {
	if args[argItem] == "" {
		return errItemFlagNotSpecified
	}
	var data userDataList
	if err := data.Restore(args[argFilename]); err != nil {
		return fmt.Errorf("restore data from file: %w", err)
	}
	if err := data.AddString(args[argItem]); err != nil {
		writer.Write([]byte(err.Error()))
	}
	if err := data.SaveTo(args[argFilename]); err != nil {
		return fmt.Errorf("save changes: %w", err)
	}
	return nil
}

func operationFindById(args Arguments, writer io.Writer) error {
	if args[argId] == "" {
		return errIdFlagNotSpecidied
	}
	var data userDataList
	if err := data.Restore(args[argFilename]); err != nil {
		return fmt.Errorf("restore data from file: %w", err)
	}
	res := data.FindById(args[argId])
	writer.Write(res)
	return nil
}

func operationRemove(args Arguments, writer io.Writer) error {
	if args[argId] == "" {
		return errIdFlagNotSpecidied
	}
	var data userDataList
	if err := data.Restore(args[argFilename]); err != nil {
		return fmt.Errorf("restore data from file: %w", err)
	}
	if err := data.RemoveById(args[argId]); err != nil {
		writer.Write([]byte(err.Error()))
	}
	if err := data.SaveTo(args[argFilename]); err != nil {
		return fmt.Errorf("save changes: %w", err)
	}
	return nil
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
