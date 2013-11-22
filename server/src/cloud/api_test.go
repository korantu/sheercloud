package cloud

import (
	"testing"
)

func TestApiTrivial(t *testing.T) {
	if false {
		t.Fail()
	}
}

func TestApiSession(t *testing.T) {
	var a, b SessionID
	if a, b = GenerateSessionID(), GenerateSessionID(); a == b {
		t.Error("Same id generated")
	}

	info_a, info_b := &SessionInfo{"User1"}, &SessionInfo{"User1"}

	a.PutInfo(info_a)
	if b.GetInfo() != nil {
		t.Error("Session must not be empty")
	}

	if !info_b.Same(a.GetInfo()) {
		t.Error("Session retrivial failed")
	}

	if info_b.Same(nil) {
		t.Error("nil is never same")
	}
}

func TestApiLogin(t *testing.T) {
	// t.Fail()

	good := &ApiLoginRequest{"sheer", "all"}
	bad := &ApiLoginRequest{"sheer", "none"}

	in, out := good.Process(), bad.Process()

	if !in.Success {
		t.Error("Failed to login proper user")
	}

	if out.Success {
		t.Error("Logged in wrong user")
	}
}
