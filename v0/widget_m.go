package colonycore

import (
	"errors"
	"github.com/eaciit/dbox"
	"github.com/eaciit/orm/v1"
	"github.com/eaciit/toolkit"
	"os"
	"path/filepath"
	"strings"
)

type Widget struct {
	orm.ModelBase
	ID           string     `json:"_id"`
	Title        string     `json:"title"`
	DataSourceId []string   `json:"dataSourceId"`
	Description  string     `json:"description"`
	Config       toolkit.Ms `json:"config"`
	Params       toolkit.M  `json:"params"`
	URL          string     `json:"url"`
}

type DataSourceWidget struct {
	ID   string      `json:"_id"`
	Data []toolkit.M `json:"data"`
}

func (w *Widget) TableName() string {
	return "widgets"
}

func (w *Widget) RecordID() interface{} {
	return w.ID
}

func (w *Widget) Get(search string) ([]Widget, error) {
	var query *dbox.Filter

	if search != "" {
		query = dbox.Contains("_id", search)
	}

	data := []Widget{}
	cursor, err := Find(new(Widget), query)
	if err != nil {
		return nil, err
	}
	if err := cursor.Fetch(&data, 0, false); err != nil {
		return nil, err
	}
	defer cursor.Close()
	return data, nil
}

func (w *Widget) GetById() error {
	if err := Get(w, w.ID); err != nil {
		return err
	}
	return nil
}

func (w *Widget) Save(extractDest string) error {
	path, err := GetWidgetPath(extractDest)
	if path == "" {
		return errors.New("directory doesn't contains index.html")
	}
	if err != nil {
		return err
	}
	urlPath := filepath.ToSlash(path)
	splitPath := strings.SplitAfter(urlPath, "/data-root/widget/")
	w.URL = strings.Join([]string{w.URL, "res-widget", splitPath[1]}, "/")

	if err := Save(w); err != nil {
		return err
	}
	return nil
}

func (w *Widget) ExtractFile(compressedSource string, fileName string) (toolkit.Ms, error) {
	compressedFile := filepath.Join(compressedSource, fileName)
	extractDest := filepath.Join(compressedSource, w.ID)

	if err := os.RemoveAll(extractDest); err != nil {
		return nil, err
	}

	if strings.Contains(fileName, ".tar.gz") {
		if err := toolkit.TarGzExtract(compressedFile, extractDest); err != nil {
			return nil, err
		}
	} else if strings.Contains(fileName, ".gz") {
		if err := toolkit.GzExtract(compressedFile, extractDest); err != nil {
			return nil, err
		}
	} else if strings.Contains(fileName, ".tar") {
		if err := toolkit.TarExtract(compressedFile, extractDest); err != nil {
			return nil, err
		}
	} else if strings.Contains(fileName, ".zip") {
		if err := toolkit.ZipExtract(compressedFile, extractDest); err != nil {
			return nil, err
		}
	}

	if err := os.Remove(compressedFile); err != nil {
		return nil, err
	}

	path, err := GetWidgetPath(extractDest)
	if path == "" {
		return nil, errors.New("directory doesn't contains index.html")
	}
	if err != nil {
		return nil, err
	}

	getConfigFile := filepath.Join(path, "config.json")
	result, err := GetJsonFile(getConfigFile)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (w *Widget) Delete(compressedSource string) error {
	extractDest := filepath.Join(compressedSource, w.ID)
	if err := Delete(w); err != nil {
		return err
	}

	if err := os.RemoveAll(extractDest); err != nil {
		return err
	}
	return nil
}
