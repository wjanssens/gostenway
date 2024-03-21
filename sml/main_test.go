package sml

import "testing"

func TestIsMethods(t *testing.T) {
	r := NewRoot()
	e, _ := r.AddElement("Test")
	a, _ := r.AddAttribute("Test", nil)
	x, _ := r.AddEmpty()

	if !r.IsRoot() {
		t.Errorf("expected root.IsRoot() == true")
	}
	if !e.IsElement() {
		t.Errorf("expected element.IsElement() == true")
	}
	if !a.IsAttribute() {
		t.Errorf("expected attribute.IsAttribute() == true")
	}
	if !x.IsEmpty() {
		t.Errorf("expected empty.IsEmpty() == true")
	}

}
