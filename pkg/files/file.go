package files

import (
	"archive/zip"
	"bufio"
	"bypctl/pkg/cmd"
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	cZip "github.com/klauspost/compress/zip"
	"github.com/mholt/archiver/v4"
	"github.com/spf13/afero"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type File struct {
	Fs afero.Fs
}

func NewFile() File {
	return File{
		Fs: afero.NewOsFs(),
	}
}

func (f File) OpenFile(dst string) (fs.File, error) {
	return f.Fs.Open(dst)
}

func (f File) GetContent(dst string) ([]byte, error) {
	afs := &afero.Afero{Fs: f.Fs}
	cByte, err := afs.ReadFile(dst)
	if err != nil {
		return nil, err
	}
	return cByte, nil
}

func (f File) SetKeyValue(dst, key, value string) error {
	afs := &afero.Afero{Fs: f.Fs}
	cByte, err := afs.ReadFile(dst)
	if err != nil {
		return err
	}
	lines := strings.Split(string(cByte), "\n")
	for i, line := range lines {
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 && parts[0] == key {
			lines[i] = fmt.Sprintf("%s=%s", key, value)
		}
	}

	newText := strings.Join(lines, "\n")

	if err := f.SaveFile(dst, newText, 0644); err != nil {
		return err
	}
	return nil
}

func (f File) CreateDir(dst string, mode fs.FileMode) error {
	return f.Fs.MkdirAll(dst, mode)
}

func (f File) CreateFile(dst string) error {
	if _, err := f.Fs.Create(dst); err != nil {
		return err
	}
	return nil
}

func (f File) LinkFile(src string, dst string, isSymlink bool) error {
	if isSymlink {
		osFs := afero.OsFs{}
		return osFs.SymlinkIfPossible(src, dst)
	} else {
		return os.Link(src, dst)
	}
}

func (f File) DeleteDir(dst string) error {
	return f.Fs.RemoveAll(dst)
}

func (f File) Stat(dst string) bool {
	info, _ := f.Fs.Stat(dst)
	return info != nil
}

func (f File) DeleteFile(dst string) error {
	return f.Fs.Remove(dst)
}

func (f File) WriteFile(dst string, in io.Reader, mode fs.FileMode) error {
	file, err := f.Fs.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err = io.Copy(file, in); err != nil {
		return err
	}

	if _, err = file.Stat(); err != nil {
		return err
	}
	return nil
}

func (f File) SaveFile(dst string, content string, mode fs.FileMode) error {
	if !f.Stat(path.Dir(dst)) {
		_ = f.CreateDir(path.Dir(dst), mode.Perm())
	}
	file, err := f.Fs.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	_, _ = write.WriteString(content)
	write.Flush()
	return nil
}

func (f File) Chmod(dst string, mode fs.FileMode) error {
	return f.Fs.Chmod(dst, mode)
}

func (f File) Chown(dst string, uid int, gid int) error {
	return f.Fs.Chown(dst, uid, gid)
}

func (f File) ChownR(dst string, uid string, gid string, sub bool) error {
	cmdStr := fmt.Sprintf(`chown %s:%s "%s"`, uid, gid, dst)
	if sub {
		cmdStr = fmt.Sprintf(`chown -R %s:%s "%s"`, uid, gid, dst)
	}
	if cmd.HasNoPasswordSudo() {
		cmdStr = fmt.Sprintf("sudo %s", cmdStr)
	}
	if msg, err := cmd.ExecWithTimeOut(cmdStr, 2*time.Second); err != nil {
		if msg != "" {
			return errors.New(msg)
		}
		return err
	}
	return nil
}

func (f File) ChmodR(dst string, mode int64) error {
	cmdStr := fmt.Sprintf(`chmod -R %v "%s"`, fmt.Sprintf("%04o", mode), dst)
	if cmd.HasNoPasswordSudo() {
		cmdStr = fmt.Sprintf("sudo %s", cmdStr)
	}
	if msg, err := cmd.ExecWithTimeOut(cmdStr, 2*time.Second); err != nil {
		if msg != "" {
			return errors.New(msg)
		}
		return err
	}
	return nil
}

func (f File) Rename(oldName string, newName string) error {
	return f.Fs.Rename(oldName, newName)
}

type WriteCounter struct {
	Total   uint64
	Written uint64
	Key     string
	Name    string
}

type Process struct {
	Total   uint64  `json:"total"`
	Written uint64  `json:"written"`
	Percent float64 `json:"percent"`
	Name    string  `json:"name"`
}

func (f File) DownloadFile(url, dst string) error {

	r, err := g.Client().Get(context.Background(), url)
	if err != nil {
		return fmt.Errorf("create download file [%s] error, err %s", dst, err)
	}
	defer r.Close()
	if err := gfile.PutBytes(dst, r.ReadAll()); err != nil {
		return fmt.Errorf("save download file [%s] error, err %s", dst, err.Error())
	}

	return nil
}

func (f File) Cut(oldPaths []string, dst string) error {
	for _, p := range oldPaths {
		base := filepath.Base(p)
		dstPath := filepath.Join(dst, base)
		if err := f.Fs.Rename(p, dstPath); err != nil {
			return err
		}
	}
	return nil
}

func (f File) Copy(src, dst string) error {
	if src = path.Clean("/" + src); src == "" {
		return os.ErrNotExist
	}
	if dst = path.Clean("/" + dst); dst == "" {
		return os.ErrNotExist
	}
	if src == "/" || dst == "/" {
		return os.ErrInvalid
	}
	if dst == src {
		return os.ErrInvalid
	}
	info, err := f.Fs.Stat(src)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return f.CopyDir(src, dst)
	}
	return f.CopyFile(src, dst)
}

func (f File) IsDir(src string) bool {
	info, err := f.Fs.Stat(src)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return true
	}
	return false
}

func (f File) IsExist(src string) bool {
	exist, err := afero.Exists(f.Fs, src)
	if err != nil {
		return false
	}
	if exist {
		return true
	}
	return false
}

func (f File) CopyDir(src, dst string) error {
	srcInfo, err := f.Fs.Stat(src)
	if err != nil {
		return err
	}
	dstDir := filepath.Join(dst, srcInfo.Name())
	if err := f.Fs.MkdirAll(dstDir, srcInfo.Mode()); err != nil {
		return err
	}

	dir, _ := f.Fs.Open(src)
	obs, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	var errs []error

	for _, obj := range obs {
		fSrc := filepath.Join(src, obj.Name())
		if obj.IsDir() {
			err = f.CopyDir(fSrc, dstDir)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			err = f.CopyFile(fSrc, dstDir)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}

	var errString string
	for _, err := range errs {
		errString += err.Error() + "\n"
	}

	if errString != "" {
		return errors.New(errString)
	}

	return nil
}

func (f File) CopyFile(src, dst string) error {
	srcFile, err := f.Fs.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	srcInfo, err := f.Fs.Stat(src)
	if err != nil {
		return err
	}
	dstPath := filepath.Join(dst, srcInfo.Name())
	if src == dstPath {
		return nil
	}

	err = f.Fs.MkdirAll(filepath.Dir(dst), 0666)
	if err != nil {
		return err
	}

	dstFile, err := f.Fs.OpenFile(dstPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0775)
	if err != nil {
		return err
	}
	defer dstFile.Close()

	if _, err = io.Copy(dstFile, srcFile); err != nil {
		return err
	}
	info, err := f.Fs.Stat(src)
	if err != nil {
		return err
	}
	if err = f.Fs.Chmod(dstFile.Name(), info.Mode()); err != nil {
		return err
	}

	return nil
}

func (f File) GetDirSize(path string) (float64, error) {
	var m sync.Map
	var wg sync.WaitGroup

	wg.Add(1)
	go ScanDir(f.Fs, path, &m, &wg)
	wg.Wait()

	var dirSize float64
	m.Range(func(k, v interface{}) bool {
		dirSize = dirSize + v.(float64)
		return true
	})

	return dirSize, nil
}

type CompressType string

const (
	Zip   CompressType = "zip"
	Gz    CompressType = "gz"
	Bz2   CompressType = "bz2"
	Tar   CompressType = "tar"
	TarGz CompressType = "tar.gz"
	Xz    CompressType = "xz"
)

func getFormat(cType CompressType) archiver.CompressedArchive {

	format := archiver.CompressedArchive{}
	switch cType {
	case Tar:
		format.Archival = archiver.Tar{}
	case TarGz, Gz:
		format.Compression = archiver.Gz{}
		format.Archival = archiver.Tar{}
	case Zip:
		format.Archival = archiver.Zip{
			Compression: zip.Deflate,
		}
	case Bz2:
		format.Compression = archiver.Bz2{}
		format.Archival = archiver.Tar{}
	case Xz:
		format.Compression = archiver.Xz{}
		format.Archival = archiver.Tar{}
	}
	return format
}

func (f File) Compress(srcRiles []string, dst string, name string, cType CompressType) error {
	format := getFormat(cType)

	fileMaps := make(map[string]string, len(srcRiles))
	for _, s := range srcRiles {
		base := filepath.Base(s)
		fileMaps[s] = base
	}

	if !f.Stat(dst) {
		_ = f.CreateDir(dst, 0755)
	}

	files, err := archiver.FilesFromDisk(nil, fileMaps)
	if err != nil {
		return err
	}
	dstFile := filepath.Join(dst, name)
	out, err := f.Fs.Create(dstFile)
	if err != nil {
		return err
	}

	switch cType {
	case Zip:
		if err := ZipFile(files, out); err != nil {
			_ = f.DeleteFile(dstFile)
			return err
		}
	default:
		err = format.Archive(context.Background(), out, files)
		if err != nil {
			_ = f.DeleteFile(dstFile)
			return err
		}
	}
	return nil
}

func isIgnoreFile(name string) bool {
	return strings.HasPrefix(name, "__MACOSX") || strings.HasSuffix(name, ".DS_Store") || strings.HasPrefix(name, "._")
}

func decodeGBK(input string) (string, error) {
	decoder := simplifiedchinese.GBK.NewDecoder()
	decoded, _, err := transform.String(decoder, input)
	if err != nil {
		return "", err
	}
	return decoded, nil
}

func (f File) Decompress(srcFile string, dst string, cType CompressType) error {
	format := getFormat(cType)

	handler := func(ctx context.Context, archFile archiver.File) error {
		info := archFile.FileInfo
		if isIgnoreFile(archFile.Name()) {
			return nil
		}
		fileName := archFile.NameInArchive
		var err error
		if header, ok := archFile.Header.(cZip.FileHeader); ok {
			if header.NonUTF8 && header.Flags == 0 {
				fileName, err = decodeGBK(fileName)
				if err != nil {
					return err
				}
			}
		}
		filePath := filepath.Join(dst, fileName)
		if archFile.FileInfo.IsDir() {
			if err := f.Fs.MkdirAll(filePath, info.Mode()); err != nil {
				return err
			}
			return nil
		} else {
			parentDir := path.Dir(filePath)
			if !f.Stat(parentDir) {
				if err := f.Fs.MkdirAll(parentDir, info.Mode()); err != nil {
					return err
				}
			}
		}
		fr, err := archFile.Open()
		if err != nil {
			return err
		}
		defer fr.Close()
		fw, err := f.Fs.OpenFile(filePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, info.Mode())
		if err != nil {
			return err
		}
		defer fw.Close()
		if _, err := io.Copy(fw, fr); err != nil {
			return err
		}

		return nil
	}
	input, err := f.Fs.Open(srcFile)
	if err != nil {
		return err
	}
	return format.Extract(context.Background(), input, nil, handler)
}

func (f File) Backup(srcFile string) (string, error) {
	backupPath := srcFile + "_bak"
	info, _ := f.Fs.Stat(backupPath)
	if info != nil {
		if info.IsDir() {
			_ = f.DeleteDir(backupPath)
		} else {
			_ = f.DeleteFile(backupPath)
		}
	}
	if err := f.Rename(srcFile, backupPath); err != nil {
		return backupPath, err
	}

	return backupPath, nil
}

func (f File) CopyAndBackup(src string) (string, error) {
	backupPath := src + "_bak"
	info, _ := f.Fs.Stat(backupPath)
	if info != nil {
		if info.IsDir() {
			_ = f.DeleteDir(backupPath)
		} else {
			_ = f.DeleteFile(backupPath)
		}
	}
	_ = f.CreateDir(backupPath, 0755)
	if err := f.Copy(src, backupPath); err != nil {
		return backupPath, err
	}
	return backupPath, nil
}

func ZipFile(files []archiver.File, dst afero.File) error {
	zw := zip.NewWriter(dst)
	defer zw.Close()

	for _, file := range files {
		hdr, err := zip.FileInfoHeader(file)
		if err != nil {
			return err
		}
		hdr.Name = file.NameInArchive
		if file.IsDir() {
			if !strings.HasSuffix(hdr.Name, "/") {
				hdr.Name += "/"
			}
			hdr.Method = zip.Store
		}
		w, err := zw.CreateHeader(hdr)
		if err != nil {
			return err
		}
		if file.IsDir() {
			continue
		}

		if file.LinkTarget != "" {
			_, err = w.Write([]byte(filepath.ToSlash(file.LinkTarget)))
			if err != nil {
				return err
			}
		} else {
			fileReader, err := file.Open()
			if err != nil {
				return err
			}
			_, err = io.Copy(w, fileReader)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// 获取指定目录的文件名称列表
func GetSubFileNames(dirPath string) ([]string, error) {
	var subFileNames []string

	files, err := os.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		subFileNames = append(subFileNames, file.Name())
	}

	return subFileNames, nil
}
