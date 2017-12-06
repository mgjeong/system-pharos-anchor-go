/*******************************************************************************
 * Copyright 2017 Samsung Electronics All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 *******************************************************************************/
package logger

import (
	"io/ioutil"
	"log"
	"os"
	"strings"
	"testing"
)

var oldStdout *os.File

func setUpLogging() (func(), *os.File, *os.File) {
	oldStdout = os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	loggers[INFO] = log.New(os.Stdout, "[INFO][EM]", logFlag)
	loggers[DEBUG] = log.New(os.Stdout, "[DEBUG][EM]", logFlag)
	loggers[ERROR] = log.New(os.Stdout, "[ERROR][EM]", logFlag)
	return func() {
		os.Stdout = oldStdout
		loggers[INFO] = log.New(os.Stdout, "[INFO][EM]", logFlag)
		loggers[DEBUG] = log.New(os.Stdout, "[DEBUG][EM]", logFlag)
		loggers[ERROR] = log.New(os.Stdout, "[ERROR][EM]", logFlag)
	}, r, w
}

func getPrintString(r *os.File, w *os.File) string {
	w.Close()
	out, _ := ioutil.ReadAll(r)
	return string(out)
}

func TestLogger(t *testing.T) {
	type testInfo struct {
		name   string
		prefix string
		suffix string
	}

	testStr := "test"
	testCase := []testInfo{
		{"INFO", "[INFO][EM]", "[" + testStr + "]\n"},
		{"DEBUG", "[DEBUG][EM]", "[" + testStr + "]\n"},
		{"ERROR", "[ERROR][EM]", "[" + testStr + "]\n"},
	}

	for i, test := range testCase {
		t.Run(test.name, func(t *testing.T) {
			tearDown, r, w := setUpLogging()

			Logging(i, testStr)
			str := getPrintString(r, w)
			if !strings.HasPrefix(str, test.prefix) {
				t.Error()
			}
			if !strings.HasSuffix(str, test.suffix) {
				t.Error()
			}

			tearDown()
		})
	}
}
