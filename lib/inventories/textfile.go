package inventories

import (
	"bufio"
	"fmt"
	"os"
	"regexp"

	"github.com/jtopjian/bagel/lib/utils"
)

// textFileValidEntry is a regular expression to match an valid
// textfile line.
var textFileValidEntry = regexp.MustCompile(`^[0-9A-Za-z\-:\._\[\]]+$`)

// TextFile represents a textfile inventory driver.
type TextFile struct {
	File string `mapstructure:"file" required:"true"`
}

// NewTextFile will return a TextFile.
func NewTextFile(options map[string]interface{}) (*TextFile, error) {
	var textfile TextFile

	err := utils.DecodeAndValidate(options, &textfile)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(textfile.File); os.IsNotExist(err) {
		return nil, fmt.Errorf("file %s does not exist", textfile.File)
	}

	return &textfile, nil
}

// Discover implements the Inventory interface for a textfile driver.
// It returns a set of hosts specified in a text file.
func (r TextFile) Discover() ([]Target, error) {
	f, err := os.Open(r.File)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var targets []Target
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		if textFileValidEntry.MatchString(line) {
			target := Target{
				Name:    line,
				Address: line,
			}

			targets = append(targets, target)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return targets, nil
}
