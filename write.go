// Copyright 2013 Unknwon
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package goconfig

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	FORMAT_COMPAT = iota // "key=value"
	FORMAT_NORMAL        // "key = value"
	FORMAT_LONG          // "%-41s = $value"
)

var defaultFormat = FORMAT_COMPAT

// SaveConfigDataFmt writes configuration to a writer
// Default equal sign   : "="
// @@ format            : customize output format
// @@ trimNullSign : support key without value
func SaveConfigDataFmt(c *ConfigFile, out io.Writer, format int, trimNullSign bool) (err error) {
	equalSign := "="

	buf := bytes.NewBuffer(nil)
	for _, section := range c.sectionList {
		// Write section comments.
		if len(c.GetSectionComments(section)) > 0 {
			if _, err = buf.WriteString(c.GetSectionComments(section) + LineBreak); err != nil {
				return err
			}
		}

		if section != DEFAULT_SECTION {
			// Write section name.
			if _, err = buf.WriteString("[" + section + "]" + LineBreak); err != nil {
				return err
			}
		}

		for num, key := range c.keyList[section] {
			if key != " " {
				// Write key comments.
				if len(c.GetKeyComments(section, key)) > 0 {
					// do not write linebreak for first line
					commentStr := c.GetKeyComments(section, key) + LineBreak
					if num != 1 {
						commentStr = LineBreak + commentStr
					}

					if _, err = buf.WriteString(commentStr); err != nil {
						return err
					}
				}

				keyName := key
				// Check if it's auto increment.
				if keyName[0] == '#' {
					keyName = "-"
				}
				//[SWH|+]:支持键名包含等号和冒号
				if strings.Contains(keyName, `=`) || strings.Contains(keyName, `:`) {
					if strings.Contains(keyName, "`") {
						if strings.Contains(keyName, `"`) {
							keyName = `"""` + keyName + `"""`
						} else {
							keyName = `"` + keyName + `"`
						}
					} else {
						keyName = "`" + keyName + "`"
					}
				}
				value := c.data[section][key]
				// In case key value contains "`" or "\"".
				if strings.Contains(value, "`") {
					if strings.Contains(value, `"`) {
						value = `"""` + value + `"""`
					} else {
						value = `"` + value + `"`
					}
				}

				// concat key=value string
				var key_value_string string

				// use pretty format
				switch format {
				case FORMAT_LONG:
					key_value_string = fmt.Sprintf("%-41s %s %s%s", keyName, equalSign, value, LineBreak)
				case FORMAT_NORMAL:
					key_value_string = fmt.Sprintf("%s %s %s%s", keyName, equalSign, value, LineBreak)
				case FORMAT_COMPAT:
					key_value_string = fmt.Sprintf("%s%s%s%s", keyName, equalSign, value, LineBreak)
				default:
					return fmt.Errorf("invalid format number: %d", format)
				}

				// support for key without value
				if trimNullSign && strings.TrimSpace(value) == "" {
					key_value_string = fmt.Sprintf("%s%s", keyName, LineBreak)
				}

				// Write key and value.
				if _, err = buf.WriteString(key_value_string); err != nil {
					return err
				}
			}
		}

		// Put a line between sections.
		if _, err = buf.WriteString(LineBreak); err != nil {
			return err
		}
	}

	if _, err := buf.WriteTo(out); err != nil {
		return err
	}

	return nil
}

// SaveConfigData writes configuration to a writer
func SaveConfigData(c *ConfigFile, out io.Writer) (err error) {
	return SaveConfigDataFmt(c, out, defaultFormat, false)
}

// SaveConfigFile writes configuration file to local file system
func SaveConfigFileFmt(c *ConfigFile, filename string, format int, trimNullSign bool) (err error) {
	// Write configuration file by filename.
	var f *os.File
	if f, err = os.Create(filename); err != nil {
		return err
	}

	if err := SaveConfigDataFmt(c, f, format, trimNullSign); err != nil {
		return err
	}
	return f.Close()
}

// SaveConfigFile writes configuration file to local file system
func SaveConfigFile(c *ConfigFile, filename string) (err error) {
	return SaveConfigFileFmt(c, filename, defaultFormat, false)
}
