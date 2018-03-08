package commons

import (
	"os"
	"path"
	"text/template"
)

// FillTemplate fills source template from file and stores it under provided destination. To fill the template provided map with mappings is used.
func FillTemplate(templatePath, finalPath string, mapping map[string]string) error {
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		Debugf("Getting template failed, err: %s", err)
		return err
	}

	base := path.Dir(finalPath)
	err = MakeDir(base)
	if err != nil {
		Debugf("Creating dir structure for new file failed, err: %s", err)
		return err
	}

	f, err := os.Create(finalPath)
	if err != nil {
		Debugf("Creating new file failed, err: %s", err)
		return err
	}
	defer func() {
		f.Close()
	}()

	err = t.Execute(f, mapping)
	if err != nil {
		Debugf("Parsing failed, err: %s", err)
		return err
	}

	si, err := os.Stat(templatePath)
	if err != nil {
		Debugf("Getting stats for template file failed, err: %s", err)
		return err
	}

	err = os.Chmod(finalPath, si.Mode())
	if err != nil {
		Debugf("Permission set for target file failed, err: %s", err)
	}

	return err
}

// MakeDir creates directory under a given path.
func MakeDir(path string) error {
	var err error
	if _, err = os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
		if err != nil {
			Debugf("Make dir error occurred, err: %s", err)
		}
	}
	return err
}
