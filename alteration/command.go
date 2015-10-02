package alteration

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/sarulabs/statix/helpers"
	"github.com/sarulabs/statix/resource"
)

// TmpInputFile defines a temporary input file that can be used in ExecCommand.
// The content of the file will be the result of the Resource dump.
// Suffix will be added at the end of the name of the temporary file.
type TmpInputFile struct {
	Resource resource.Resource
	Suffix   string
}

// TmpOutputFile defines a temporary output file that can be used in ExecCommand.
// Suffix will be added at the end of the name of the temporary file.
type TmpOutputFile struct {
	Suffix string
}

// ExecCommand executes a command that returns a resource.
// Arguments should be strings, except for two.
// - one may be a TmpInputFile
// - one may be a TmpOutputFile
// That is because you will often need a temporary input file
// and a temporary output file to  generate a resource.
// The returned Resource content is the command standard output
// except when a TmpOutputFile is used.
func ExecCommand(name string, args ...interface{}) (resource.Resource, error) {
	parsedArgs, err := parseArgs(args...)
	if err != nil {
		return &resource.Empty{}, err
	}

	var inputFile, outputFile *os.File

	// Create temporary input file if necessary
	if parsedArgs.tmpInputFile != nil {
		inputFile, err = helpers.TempFile("", "statix_filter_", parsedArgs.tmpInputFile.Suffix)
		if err != nil {
			return &resource.Empty{}, err
		}
		defer func() {
			inputFile.Close()
			os.Remove(inputFile.Name())
		}()

		content, err := parsedArgs.tmpInputFile.Resource.Dump()
		if err != nil {
			return &resource.Empty{}, err
		}

		_, err = inputFile.Write(content)
		if err != nil {
			return &resource.Empty{}, err
		}
	}

	// Create temporary output file if necessary
	if parsedArgs.tmpOutputFile != nil {
		outputFile, err = helpers.TempFile("", "statix_filter_", parsedArgs.tmpOutputFile.Suffix)
		if err != nil {
			return &resource.Empty{}, err
		}
		defer func() {
			outputFile.Close()
			os.Remove(outputFile.Name())
		}()
	}

	// update args with the name of temporary files
	for i, arg := range parsedArgs.args {
		if arg == statixTmpInputFile {
			parsedArgs.args[i] = inputFile.Name()
		}
		if arg == statixTmpOutputFile {
			parsedArgs.args[i] = outputFile.Name()
		}
	}

	// execute command
	bufOut := bytes.NewBuffer(nil)
	bufErr := bytes.NewBuffer(nil)
	command := exec.Command(name, parsedArgs.args...)
	command.Stdout = bufOut
	command.Stderr = bufErr
	err = command.Run()
	if err != nil {
		return &resource.Empty{}, fmt.Errorf("command error on (%s, %s) :\n%s", name, parsedArgs.args, bufErr.String())
	}

	if parsedArgs.tmpOutputFile == nil {
		return resource.NewBytes(bufOut.Bytes()), nil
	}

	c, err := ioutil.ReadFile(outputFile.Name())
	if err != nil {
		return &resource.Empty{}, err
	}

	return resource.NewBytes(c), nil
}

const statixTmpInputFile = "{{STATIX_TMP_INPUT_FILE}}"
const statixTmpOutputFile = "{{STATIX_TMP_OUTPUT_FILE}}"

type parsedArgs struct {
	tmpInputFile  *TmpInputFile
	tmpOutputFile *TmpOutputFile
	args          []string
}

func parseArgs(args ...interface{}) (parsedArgs, error) {
	pcas := parsedArgs{}
	for _, a := range args {
		switch arg := a.(type) {
		case TmpInputFile:
			if pcas.tmpInputFile != nil {
				return parsedArgs{}, errors.New("only one TmpInputFile allowed")
			}
			pcas.tmpInputFile = &arg
			pcas.args = append(pcas.args, statixTmpInputFile)

		case TmpOutputFile:
			if pcas.tmpOutputFile != nil {
				return parsedArgs{}, errors.New("only one tmpOutputFile allowed")
			}
			pcas.tmpOutputFile = &arg
			pcas.args = append(pcas.args, statixTmpOutputFile)

		case string:
			pcas.args = append(pcas.args, arg)

		default:
			return parsedArgs{}, errors.New("argument type not supported")
		}
	}
	return pcas, nil
}
