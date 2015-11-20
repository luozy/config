/*
 * ------------------------------------------------------------
 * Author: luoziyuan@gmail.com
 * Date: 2015/11/19
 * ------------------------------------------------------------
 */

package config

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

var (
	cfg   *Config    // current instance
	mutex sync.Mutex // current instance mutex
)

const (
	prefix = "["
	suffix = "]"
)

type section struct {
	name    string
	parent  string
	storage map[string]string
}

type Config struct {
	fileName string
	storage  map[string]*section
}

func init() {
}

// new config instance
func NewConfig(file string) (*Config, error) {
	mutex.Lock()
	defer mutex.Unlock()

	c := &Config{
		fileName: file,
		storage:  make(map[string]*section),
	}

	if err := c.Load(); nil != err {
		return nil, err
	} else {
		cfg = c
		return c, nil
	}
}

// get current instance k/v by int
func GetInt(section, key string) int {
	if str := GetStr(section, key); "" != str {
		if v, err := strconv.ParseInt(str, 10, 64); nil == err {
			return int(v)
		}
	}

	return 0
}

func GetUint(section, key string) uint {
	if str := GetStr(section, key); "" != str {
		if v, err := strconv.ParseUint(str, 10, 64); nil == err {
			return uint(v)
		}
	}

	return 0
}

func GetStr(section, key string) string {
	return getString(section, key)
}

func getString(section, key string) string {
	mutex.Lock()
	defer mutex.Unlock()

	if nil != cfg {
		fmt.Println("getString -1")
		return cfg.GetString(section, key)
	} else {
		fmt.Println("getString -2")
		return ""
	}
}

func GetFloat64(section, key string) float64 {
	if str := GetStr(section, key); "" != str {
		if v, err := strconv.ParseFloat(str, 64); nil != err {
			return v
		}
	}

	return 0.0
}

func GetFloat32(section, key string) float32 {
	return float32(GetFloat64(section, key))
}

func IsYes(section, key string) bool {
	str := GetStr(section, key)
	str = strings.ToLower(str)
	return "yes" == str
}

func SetInt(section, key string, value int) error {
	if nil == cfg {
		return fmt.Errorf("config not initialized")
	}

	return nil
}

func SetStr(section, key, value string) error {
	if nil == cfg {
		return fmt.Errorf("config not initialized")
	}
	return nil
}

func DelSection(section string) error {
	if nil == cfg {
		return fmt.Errorf("config not initialized")
	}

	return nil
}

func (this *Config) Load() error {
	fp, err := os.Open(this.fileName)
	if nil != err {
		return err
	}

	defer fp.Close()

	buf := bufio.NewReader(fp)
	var section string
	for {
		line := ""
		if tmp, _, err := buf.ReadLine(); nil != err {
			if io.EOF == err {
				break
			}
			return err
		} else {
			line = strings.TrimSpace(string(tmp))
		}

		if "" == line || strings.HasPrefix(line, "#") {
			continue
		}

		if strings.HasPrefix(line, prefix) && strings.HasSuffix(line, suffix) {
			section = strings.TrimSpace(line[1 : len(line)-1])
			//section = strings.ToLower(section)
			if _, ok := this.storage[section]; !ok {
				this.storage[section] = this.newSection(section)
			}
			continue
		}

		pair := strings.SplitN(line, "=", 2)
		if 2 == len(pair) {
			key := strings.TrimSpace(pair[0])
			val := strings.TrimSpace(pair[1])
			bpos := strings.Index(val, " ")
			epos := strings.Index(val, "#")
			if _, ok := this.storage[section]; !ok {
				this.storage[section] = this.newSection(section)
			}

			if bpos > 0 && epos >= bpos {
				val = strings.TrimSpace(val[0:bpos])
			}

			this.storage[section].storage[key] = val
		}
	}

	return nil
}

func (this *Config) newSection(sec string) *section {
	storage := &section{
		name:    sec,
		parent:  "",
		storage: make(map[string]string),
	}

	return storage
}

func (this *Config) GetInt(section, key string) int {
	return int(this.GetInt64(section, key))
}

func (this *Config) GetInt64(section, key string) int64 {
	if str := this.GetString(section, key); "" != str {
		if v, err := strconv.ParseInt(str, 10, 64); nil == err {
			return v
		}
	}

	return 0
}

func (this *Config) GetInt32(section, key string) int32 {
	return int32(this.GetInt64(section, key))
}

func (this *Config) GetString(section, key string) string {
	var val string
	if sec, ok := this.storage[section]; ok {
		if val, ok = sec.storage[key]; !ok && "" != sec.parent {
			return this.GetString(sec.parent, key)
		}
	}

	return val
}

func (this *Config) SetString(section, key, val string) {
	if _, ok := this.storage[section]; !ok {
		this.storage[section] = this.newSection(section)
	}

	key = strings.TrimSpace(key)
	val = strings.TrimSpace(val)
	this.storage[section].storage[key] = val

}

func (this *Config) DeleteSection(section string) {
	delete(this.storage, section)
}

func (this *Config) DeleteKey(section, key string) {
	if sec, ok := this.storage[section]; ok {
		delete(sec.storage, key)
	}
}
