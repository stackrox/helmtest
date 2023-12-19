package hack

import (
	"testing"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport/spdy"
)

func TestAPIMachineryUsable(_ *testing.T) {
	_, _, _ = spdy.RoundTripperFor(&rest.Config{})
}
