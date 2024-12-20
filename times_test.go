package toolbox

import (
	"fmt"
	"testing"
)

func TestNowTime(t *testing.T) {
	fmt.Println(NewUtils().NewTimesUtils().NowTime())
}

func TestTimeString(t *testing.T) {
	fmt.Println(NewUtils().NewTimesUtils().NowTimeString())
}
