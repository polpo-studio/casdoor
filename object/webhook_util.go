// Copyright 2021 The Casdoor Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package object

import (
	"github.com/astaxie/beego/logs"
	"io"
	"net/http"
	"strings"

	"github.com/casdoor/casdoor/util"
)

func sendWebhook(webhook *Webhook, record *Record) error {
	client := &http.Client{}

	body := strings.NewReader(util.StructToJson(record))

	req, err := http.NewRequest(webhook.Method, webhook.Url, body)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", webhook.ContentType)

	for _, header := range webhook.Headers {
		req.Header.Set(header.Name, header.Value)
	}

	resp, err := client.Do(req)

	if err != nil {
		logs.Warning("Send webhook failed, %s, resp: %v", err.Error(), resp)
		logs.Warning("Send webhook failed, req body: %v, url: %s", body, webhook.Url)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, _ := io.ReadAll(resp.Body)

		logs.Warning("Send webhook failed, req body: %v, url: %s", body, webhook.Url)
		logs.Error("Send webhook failed, status code: %d, resp body: %s", resp.StatusCode, string(bodyBytes))
	}

	return err
}
