// Copyright (C) 2022 Specter Ops, Inc.
//
// This file is part of AzureHound.
//
// AzureHound is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// AzureHound is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package sinks

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/certmichelin/azurehound/v2/models"
	"github.com/certmichelin/azurehound/v2/pipeline"
)

func WriteToFile[T any](ctx context.Context, filePath string, stream <-chan T) error {

	if file, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666); err != nil {
		return err
	} else {
		defer file.Close()

		if _, err := file.WriteString("{\n\t\"data\": [\n"); err != nil {
			return err
		} else {
			meta := models.Meta{
				Type:    "azure",
				Version: 5,
				Count:   0,
			}

			format := "\t\t%v"
			for item := range pipeline.OrDone(ctx.Done(), stream) {
				if _, err := file.WriteString(fmt.Sprintf(format, item)); err != nil {
					return err
				}
				meta.Count++
				format = ",\n\t\t%v"
			}

			if bytes, err := json.Marshal(meta); err != nil {
				return err
			} else if _, err := file.WriteString(fmt.Sprintf("\n\t],\n\t\"meta\": %s\n}\n", string(bytes))); err != nil {
				return err
			} else {
				return nil
			}
		}
	}
}
