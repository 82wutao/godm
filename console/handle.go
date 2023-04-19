package console

import (
	"fmt"
	"os"

	"dm.net/datamine/common"
)

// Alg is alg of handle data
type Alg common.Func[interface{}, interface{}]

func parseAlg(filePath string) Alg {
	return nil
}

// HandleFile handle data in console
func HandleFile(dataFilePath string, delim byte, algFilePath string) {
	algObj := parseAlg(algFilePath)

	reader, err := NewLineReader(dataFilePath, delim)
	if err != nil {
		fmt.Fprintf(os.Stderr, "")
		return
	}
	for line, lineErr := reader.Line(); lineErr != nil; line, lineErr = reader.Line() {
		algObj(line)
	}
	reader.Close()
}
