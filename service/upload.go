package service

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"zone.com/common"
	"zone.com/util"
)

func (s *service) upload(c echo.Context) error {
	//-----------
	// Read file
	//-----------

	// Source
	clientfd, err := c.FormFile("file")
	if err != nil {
		return err
	}
	src, err := clientfd.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	// Destination
	destLocalPath := "./temp/"
	checkAndCreateDir(destLocalPath)

	fileName := "temp.jpg"
	if c.FormValue("name") != "" {
		fileName = c.FormValue("name")
	}
	localpath := fmt.Sprintf("%s%s", destLocalPath, fileName)
	dst, err := os.OpenFile(localpath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer dst.Close()

	// Copy
	// 利用io.TeeReader在读取文件内容时计算hash值
	fhash := sha1.New()
	io.Copy(dst, io.TeeReader(src, fhash))
	hstr := hex.EncodeToString(fhash.Sum(nil))

	// copy to web dir
	targetFile := fmt.Sprintf("%s%s", util.FileDir, fileName)
	checkAndCreateDir(util.FileDir)
	if _, err := copy(localpath, targetFile); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, common.NewHttpMsgData(echo.Map{
		"fileName": fileName,
		"hash":     hstr,
	}))
}

func checkAndCreateDir(path string) {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		os.MkdirAll(path, 0755)
	}
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
