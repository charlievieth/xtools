package imports

import (
	"os"
	"testing"

	"github.com/charlievieth/xtools/testenv"
)

func TestMain(m *testing.M) {
	testenv.ExitIfSmallMachine()
	os.Exit(m.Run())
}
