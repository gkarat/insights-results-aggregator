/*
Copyright © 2020 Red Hat, Inc.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package migration

import (
	"database/sql"

	"github.com/RedHatInsights/insights-results-aggregator/types"
)

var mig0010AddTagsFieldToRuleErrorKeyTable = Migration{
	StepUp: func(tx *sql.Tx, _ types.DBDriver) error {
		_, err := tx.Exec(`
			ALTER TABLE rule_error_key ADD COLUMN tags VARCHAR NOT NULL DEFAULT ''
		`)
		return err
	},
	StepDown: func(tx *sql.Tx, driver types.DBDriver) error {
		if driver == types.DBDriverSQLite3 {
			return downgradeTable(tx, ruleErrorKeyTable, `
				CREATE TABLE rule_error_key (
					"error_key"     VARCHAR NOT NULL,
					"rule_module"   VARCHAR NOT NULL REFERENCES rule(module) ON DELETE CASCADE,
					"condition"     VARCHAR NOT NULL,
					"description"   VARCHAR NOT NULL,
					"impact"        INTEGER NOT NULL,
					"likelihood"    INTEGER NOT NULL,
					"publish_date"  TIMESTAMP NOT NULL,
					"active"        BOOLEAN NOT NULL,
					"generic"       VARCHAR NOT NULL,
					PRIMARY KEY("error_key", "rule_module")
				)
			`, []string{"error_key", "rule_module", "condition", "description", "impact", "likelihood", "publish_date", "active", "generic"})
		}

		_, err := tx.Exec(`
			ALTER TABLE rule_error_key DROP COLUMN tags
		`)
		return err
	},
}