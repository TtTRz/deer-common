/* Python Compiler Provider
 * (C) 2019 LanceLRQ
 *
 * This code is licenced under the GPLv3.
 */
package provider

import (
    "fmt"
    "strings"
)

type Py2CompileProvider struct {
    CodeCompileProvider
}

type Py3CompileProvider struct {
    CodeCompileProvider
}

func NewPy2CompileProvider () *Py2CompileProvider {
    return &Py2CompileProvider{
        CodeCompileProvider{
            isReady: false,
            realTime: true,
            Name: "python2",
        },
    }
}


func (prov *Py2CompileProvider) Init(code string, workDir string) error {
    prov.realTime = true
    prov.codeContent = code
    prov.workDir = workDir
    prov.Name = "python2"

    err := prov.checkWorkDir()
    if err != nil {
        return err
    }

    err = prov.initFiles(".py", "")
    if err != nil {
        return nil
    }
    prov.isReady = true
    return nil
}

func (prov *Py2CompileProvider) Compile() (result bool, errmsg string) {
    return true, ""
}

func (prov *Py2CompileProvider) GetRunArgs() (args []string) {
    argsRaw := fmt.Sprintf(CompileCommands.Python2, prov.codeFilePath)
    args = strings.Split(argsRaw, " ")
    return
}

func (prov *Py2CompileProvider) IsCompileError(remsg string) bool {
    return strings.Contains(remsg, "SyntaxError") ||
        strings.Contains(remsg, "IndentationError") ||
        strings.Contains(remsg, "ImportError")
}


func NewPy3CompileProvider () *Py3CompileProvider {
    return &Py3CompileProvider{
        CodeCompileProvider{
            isReady: false,
            realTime: true,
            Name: "python3",
        },
    }
}


func (prov *Py3CompileProvider) Init(code string, workDir string) error {
    prov.codeContent = code
    prov.workDir = workDir

    err := prov.checkWorkDir()
    if err != nil {
        return err
    }

    err = prov.initFiles(".py", "")
    if err != nil {
        return nil
    }
    prov.isReady = true
    return nil
}

func (prov *Py3CompileProvider) Compile() (result bool, errmsg string) {
    return true, ""
}

func (prov *Py3CompileProvider) GetRunArgs() (args []string) {
    argsRaw := fmt.Sprintf(CompileCommands.Python3, prov.codeFilePath)
    args = strings.Split(argsRaw, " ")
    return
}

func (prov *Py3CompileProvider) IsCompileError(remsg string) bool {
    return strings.Contains(remsg, "SyntaxError") ||
        strings.Contains(remsg, "IndentationError") ||
        strings.Contains(remsg, "ImportError")
}
