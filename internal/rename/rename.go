package rename

import "errors"

// Manager is the interface required by Run.
type Manager interface {
	Load(name string) (map[string]string, error)
	Save(name string, vars map[string]string) error
	Delete(name string) error
}

// Run renames srcName to dstName in the profile store.
// It loads the source snapshot, saves it under the destination name,
// then deletes the source. Returns an error if src is missing or dst
// already exists (unless overwrite is true).
func Run(m Manager, srcName, dstName string, overwrite bool) error {
	if srcName == "" {
		return errors.New("source name must not be empty")
	}
	if dstName == "" {
		return errors.New("destination name must not be empty")
	}
	if srcName == dstName {
		return errors.New("source and destination names are identical")
	}

	vars, err := m.Load(srcName)
	if err != nil {
		return fmt.Errorf("load %q: %w", srcName, err)
	}

	if !overwrite {
		if _, err := m.Load(dstName); err == nil {
			return fmt.Errorf("destination %q already exists; use --overwrite to replace", dstName)
		}
	}

	if err := m.Save(dstName, vars); err != nil {
		return fmt.Errorf("save %q: %w", dstName, err)
	}

	if err := m.Delete(srcName); err != nil {
		return fmt.Errorf("delete source %q: %w", srcName, err)
	}

	return nil
}
