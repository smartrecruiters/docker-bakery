package commons

import (
	"os"
	"path"
	"text/template"
)

func FillTemplate(templatePath, finalPath string, mapping map[string]string) (err error) {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		Debugf("Getting template failed, err: %s", err)
		return
	}

	base := path.Dir(finalPath)
	err = MakeDir(base)
	if err != nil {
		Debugf("Creating dir structure for new file failed, err: %s", err)
		return
	}

	f, err := os.Create(finalPath)
	if err != nil {
		Debugf("Creating new file failed, err: %s", err)
		return
	}
	defer func() {
		if e := f.Close(); e != nil {
			err = e
		}
	}()

	err = t.Execute(f, mapping)
	if err != nil {
		Debugf("Parsing failed, err: %s", err)
		return
	}

	si, err := os.Stat(templatePath)
	if err != nil {
		Debugf("Getting stats for template file failed, err: %s", err)
		return
	}
	err = os.Chmod(finalPath, si.Mode())
	if err != nil {
		Debugf("Permission set for target file failed, err: %s", err)
		return
	}
	return
}

func MakeDir(path string) (err error) {
	err = nil
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			Debugf("Make dir error occurred, err: %s", err)
		}
	}
	return
}
