package shellcoder

import "testing"

func TestShellcode(t *testing.T) {
	t.Run("TestShellcode", func(t *testing.T) {
		t.Log("TestShellcode")
		b, err := GenerateLinuxX64ShellcodeFromBytes([]byte{'E', 'L', 'F', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'})
		if err != nil {
			t.Fail()
		}
		t.Log(b)
	})
}
