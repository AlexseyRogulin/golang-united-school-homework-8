package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type userData struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	Age   uint   `json:"age"`
}

type userDataList []userData

func (d *userDataList) Restore(filename string) error {
	rawData, err := getFileRawData(filename)
	if err != nil {
		return err
	}
	if len(rawData) == 0 {
		return nil
	}
	var newData userDataList
	if err := json.Unmarshal(rawData, &newData); err != nil {
		return fmt.Errorf("decode: %w", err)
	}
	*d = newData
	fmt.Println(*d)
	return nil
}

func getFileRawData(filename string) ([]byte, error) {
	file, err := os.OpenFile(filename, os.O_RDONLY, 0666)
	if err != nil {
		if os.IsNotExist(err) {
			return []byte{}, nil
		}
		return []byte{}, fmt.Errorf("open filen %q: %w", filename, err)
	}
	defer file.Close()
	rawData, err := io.ReadAll(file)
	if err != nil {
		return []byte{}, fmt.Errorf("read all file: %w", err)
	}
	return rawData, nil
}

func (d *userDataList) AddString(data string) error {
	var user userData
	if err := json.Unmarshal([]byte(data), &user); err != nil {
		return fmt.Errorf("unmarshal %q: %w", data, err)
	}
	if d.checkExistingId(user.Id) {
		return fmt.Errorf("Item with id %v already exists", user.Id)
	}
	*d = append(*d, user)
	return nil
}

func (d *userDataList) checkExistingId(id string) bool {
	for _, user := range *d {
		if user.Id == id {
			return true
		}
	}
	return false
}

func (d *userDataList) SaveTo(filename string) error {
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return fmt.Errorf("open filen %q: %w", filename, err)
	}
	data, err := json.Marshal(d)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	if _, err := file.Write(data); err != nil {
		return fmt.Errorf("write file %q: %w", filename, err)
	}
	return nil
}

func (d *userDataList) FindById(id string) []byte {
	for _, user := range *d {
		if user.Id == id {
			data, _ := json.Marshal(user)
			return data
		}
	}
	return []byte{}
}

func (d *userDataList) RemoveById(id string) error {
	for i, user := range *d {
		if user.Id == id {
			*d = append((*d)[:i], (*d)[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Item with id %v not found", id)
}
