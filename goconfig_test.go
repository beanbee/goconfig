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
	"os"

	"testing"
)

func TestLoadConfigFile(t *testing.T) {
	configTest := `; Google
google = www.google.com
search = http://%(google)s

; Here are Comments
; Second line
[Demo]
# This symbol can also make this line to be comments
key1 = Let's us goconfig!!!
key2 = rewrite this key of conf.ini
key3 = this is based on key2:%(key2)s
quote = "special case for quote
中国 = China
chinese-var = hello %(中国)s!
array_key = 1,2,3,4,5

[What's this?]
; Not Enough Comments!!
# line 2
name = try one more value ^-^
empty_value = 
no_value`
	c, err := LoadFromData([]byte(configTest))
	if err != nil {
		t.Logf("load data failed %v", err)
	}
	tmpFile := "test_save.ini"

	// TrimNullValueSign = false
	c.SetKeyComments("Demo", "chinese-var", "comment1")
	c.SetKeyComments("Demo", "key2", "comment2")

	SaveConfigFile(c, tmpFile)

	os.Remove(tmpFile)
}
