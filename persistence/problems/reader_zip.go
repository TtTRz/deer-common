package problems

import (
    "archive/zip"
    "bufio"
    "bytes"
    "crypto/rsa"
    "github.com/LanceLRQ/deer-common/persistence"
    commonStructs "github.com/LanceLRQ/deer-common/structs"
    "github.com/LanceLRQ/deer-common/utils"
    "github.com/pkg/errors"
    "golang.org/x/crypto/openpgp"
    "io"
    "io/ioutil"
    "log"
    "path"
    "reflect"
)


// 校验判题结果数据包
func validateProblemPackageZip(zipArchive *zip.ReadCloser) (bool, error) {
    // 搜索config文件
    configFile, _, err := FindInZip(zipArchive, "problem.json")
    if err != nil {
        return false, err
    }

    configBytes, err := ioutil.ReadAll(*configFile)
    if err != nil {
        return false, err
    }

    hash, err := persistence.SHA256Streams([]io.Reader{ bytes.NewReader(configBytes) })
    if err != nil {
        return false, err
    }

    // signFile
    signFile, _, err := FindInZip(zipArchive, ".sign")
    if err != nil {
        return false, err
    }

    signature, err := ioutil.ReadAll(*signFile)
    if err != nil {
        return false, err
    }

    // gpgFile
    withGPG := true
    gpgFile, _, err := FindInZip(zipArchive, ".gpg")
    if err != nil {
        withGPG = false
    }

    // 进行签名校验
    if withGPG && gpgFile != nil {
        // Read GPG Keys
        elist, err := openpgp.ReadArmoredKeyRing(bufio.NewReader(*gpgFile))
        if err != nil {
            return false, err
        }
        if len(elist) < 1 {
            return false, errors.Errorf("GPG key error")
        }
        publicKey := elist[0].PrimaryKey.PublicKey.(*rsa.PublicKey)
        err = persistence.RSA2048Verify(hash, signature, publicKey)
        if err != nil {
            return false, err
        }
    } else {
        return reflect.DeepEqual(hash, signature), nil
    }
    return true, nil
}

func doProblemPackageValidationZip (zipArchive *zip.ReadCloser, validate bool) error {
    ok, err := validateProblemPackageZip(zipArchive)
    var errmsg error
    if !ok || err != nil {
        if err != nil {
            errmsg = errors.Errorf("validate package hash error: %s", err.Error())
        }
        errmsg = errors.Errorf("validate package hash error")
    }
    // 如果出错并且现在必须要验证错误，则抛出
    if errmsg != nil {
        if validate {
            return errmsg
        } else {
            log.Println("Warning! Package signature validation failed.")
        }
    }
    return nil
}

// 读取题目信息(ZIP)
// workDir只在需要解压的时候才会用到
func ReadProblemInfoZip(problemFile string, unpack, validate bool, workDir string) (*commonStructs.JudgeConfiguration, string, error) {
    // 打开文件
    zipReader, err := zip.OpenReader(problemFile)
    if err != nil {
       return nil, "", errors.Errorf("open file (%s) error: %s", problemFile, err.Error())
    }
    defer zipReader.Close()

    config := commonStructs.JudgeConfiguration{}
    // 搜索config文件
    configFile, _, err := FindInZip(zipReader, "problem.json")
    if err != nil {
        return nil, "", err
    }
    // 读取Config信息
    configByte, err := ioutil.ReadAll(*configFile)
    if err != nil {
        return nil, "", err
    }
    utils.JSONBytesObject(configByte, &config)
    // 校验签名
    err = doProblemPackageValidationZip(zipReader, validate)
    if err != nil {
        return nil, "", err
    }

    if unpack {
        // 解压
        err = UnZip(zipReader, workDir)
        if err != nil {
            return nil, "", err
        }
    }

    return &config, path.Join(workDir, "problem.json"), nil
}

// 读取题目携带的GPG信息(ZIP)
func ReadProblemGPGInfoZip(problemFile string) (string, error) {
    zipReader, err := zip.OpenReader(problemFile)
    if err != nil {
        return "", errors.Errorf("open file (%s) error: %s", problemFile, err.Error())
    }
    defer zipReader.Close()

    file, _, err := FindInZip(zipReader, ".gpg")
    if err != nil {
        if IsFileNotFoundError(err) {
            return "", errors.Errorf("no GPG public key")
        }
        return "", err
    }

    elist, err := openpgp.ReadArmoredKeyRing(bufio.NewReader(*file))
    if err != nil {
        return "", err
    }
    if len(elist) < 1 {
        return "", errors.Errorf("GPG key error")
    }
    rel := ""
    for _, identify := range elist[0].Identities {
        rel += identify.Name + "\n"
    }
    return rel, nil
}