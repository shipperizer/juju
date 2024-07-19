// Code generated by triggergen. DO NOT EDIT.

package triggers

import (
	"fmt"

	"github.com/juju/juju/core/database/schema"
)


// ChangeLogTriggersForCharm generates the triggers for the 
// charm table.
func ChangeLogTriggersForCharm(columnName string, namespaceID int) func() schema.Patch {
	return func() schema.Patch {
		return schema.MakePatch(fmt.Sprintf(`
-- insert trigger for Charm
CREATE TRIGGER trg_log_charm_insert
AFTER INSERT ON charm FOR EACH ROW
BEGIN
    INSERT INTO change_log (edit_type_id, namespace_id, changed, created_at)
    VALUES (1, %[2]d, NEW.%[1]s, DATETIME('now'));
END;

-- update trigger for Charm
CREATE TRIGGER trg_log_charm_update
AFTER UPDATE ON charm FOR EACH ROW
WHEN 
	(NEW.name != OLD.name OR (NEW.name IS NOT NULL AND OLD.name IS NULL) OR (NEW.name IS NULL AND OLD.name IS NOT NULL)) OR
	(NEW.description != OLD.description OR (NEW.description IS NOT NULL AND OLD.description IS NULL) OR (NEW.description IS NULL AND OLD.description IS NOT NULL)) OR
	(NEW.summary != OLD.summary OR (NEW.summary IS NOT NULL AND OLD.summary IS NULL) OR (NEW.summary IS NULL AND OLD.summary IS NOT NULL)) OR
	(NEW.subordinate != OLD.subordinate OR (NEW.subordinate IS NOT NULL AND OLD.subordinate IS NULL) OR (NEW.subordinate IS NULL AND OLD.subordinate IS NOT NULL)) OR
	(NEW.min_juju_version != OLD.min_juju_version OR (NEW.min_juju_version IS NOT NULL AND OLD.min_juju_version IS NULL) OR (NEW.min_juju_version IS NULL AND OLD.min_juju_version IS NOT NULL)) OR
	(NEW.run_as_id != OLD.run_as_id OR (NEW.run_as_id IS NOT NULL AND OLD.run_as_id IS NULL) OR (NEW.run_as_id IS NULL AND OLD.run_as_id IS NOT NULL)) OR
	(NEW.assumes != OLD.assumes OR (NEW.assumes IS NOT NULL AND OLD.assumes IS NULL) OR (NEW.assumes IS NULL AND OLD.assumes IS NOT NULL)) OR
	(NEW.lxd_profile != OLD.lxd_profile OR (NEW.lxd_profile IS NOT NULL AND OLD.lxd_profile IS NULL) OR (NEW.lxd_profile IS NULL AND OLD.lxd_profile IS NOT NULL)) OR
	(NEW.archive_path != OLD.archive_path OR (NEW.archive_path IS NOT NULL AND OLD.archive_path IS NULL) OR (NEW.archive_path IS NULL AND OLD.archive_path IS NOT NULL)) 
BEGIN
    INSERT INTO change_log (edit_type_id, namespace_id, changed, created_at)
    VALUES (2, %[2]d, OLD.%[1]s, DATETIME('now'));
END;

-- delete trigger for Charm
CREATE TRIGGER trg_log_charm_delete
AFTER DELETE ON charm FOR EACH ROW
BEGIN
    INSERT INTO change_log (edit_type_id, namespace_id, changed, created_at)
    VALUES (4, %[2]d, OLD.%[1]s, DATETIME('now'));
END;`, columnName, namespaceID))
	}
}

