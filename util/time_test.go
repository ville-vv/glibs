package util

import (
	"testing"
)

func TestTimeParseLocal(t *testing.T) {
	t.Log(TimeParseLocal("2019/03/26"))
	t.Log(TimeParseLocal("03/26/2019"))
	t.Log(TimeParseLocal("2019/03/26 17:30:03"))
	t.Log(TimeParseLocal("03/26/2019 17:30:03"))
}
