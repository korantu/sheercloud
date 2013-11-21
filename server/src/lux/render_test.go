/**
 * Created with IntelliJ IDEA.
 */
package lux

import "testing"

func TestLuxTest(t * testing.T) {
	if err := CheckLux(); err != nil {
	   t.Error(err.Error())
	}
}

