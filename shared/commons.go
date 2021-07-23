package shared

import (
	"encoding/json"
	"fmt"
	maps "github.com/mitchellh/mapstructure"
	"io/ioutil"
	"os"
)

func Backup(i interface{}, new string) error {
	if _, err := os.Stat(new); err == nil {
		// file does exist
		return os.ErrExist
	} else if os.IsNotExist(err) {
		// ok, file does not exist
	} else {
		// file may or may not exist. See err for details.
		return err
	}
	bytes, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(new, bytes, 0644)
}

func Save(p Printer, path string) error {
	bytes, err := p.Print()
	if err != nil {
		return err
	}
	return ioutil.WriteFile(path, bytes, 0644)
}

func Print(v Validator) ([]byte, error) {
	if err := v.Validate(); err != nil {
		return nil, err
	}
	return json.MarshalIndent(v, "", "\t")
}

func Load(i interface{}, path, version string) (interface{}, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}

	// TODO backup raw here

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var input map[string]interface{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err
	}

	if err := checkVersion(input, version); err != nil {
		return nil, err
	}

	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{Metadata: &md, TagName: "json", Result: &i})
	if err != nil {
		return nil, err
	}
	err = d.Decode(input)
	if err != nil {
		return nil, err
	}
	if u, ok := i.(WithUnused); ok {
		u.SetUnused(md.Unused)
	}

	return i, nil
}

func Upgrade(u Upgrader, path string) (interface{}, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return nil, err
	}

	// TODO backup raw here

	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var input map[string]interface{}
	err = json.Unmarshal(b, &input)
	if err != nil {
		return nil, err
	}

	err = u.UpgradeFunc(input)
	if err != nil {
		return nil, err
	}

	var md maps.Metadata
	d, err := maps.NewDecoder(&maps.DecoderConfig{Metadata: &md, TagName: "json", Result: &u})
	if err != nil {
		return nil, err
	}
	err = d.Decode(input)
	if err != nil {
		return nil, err
	}
	if u, ok := u.(WithUnused); ok {
		u.SetUnused(md.Unused)
	}

	return u, nil
}

func GetVersion(input map[string]interface{}) (string, error) {
	meta, ok := input["meta"].(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("structure doesn't look like one we can understand - does not have meta object")
	}
	v, ok := meta["version"].(string)
	if !ok {
		return "", fmt.Errorf("structure doesn't look like one we can understand - meta object does not have version field")
	}
	return v, nil
}

func checkVersion(input map[string]interface{}, version string) error {
	v, err := GetVersion(input)
	if err != nil {
		return err
	}
	if version != v {
		return NotCurrentVersionError{Version: v}
	}
	return nil
}
