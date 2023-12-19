package hack

import (
	"k8s.io/client-go/rest"
	"k8s.io/client-go/transport/spdy"
	"testing"
)

func TestAPIMachineryUsable(t *testing.T) {
	_, _, _ = spdy.RoundTripperFor(&rest.Config{})
}
