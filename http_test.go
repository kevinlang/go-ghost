package ghost

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const ExampleAdminKey = "5ea1aeb17edc2650468b6554:0f1103f5af0395a73041457eb6928f9e0d143a8dcba187915342e65687e2a589"

func TestNewTokenAuthClient(t *testing.T) {
	c, err := NewTokenAuthClient(ExampleAdminKey)
	require.NoError(t, err)
	require.NotNil(t, c)
}

func TestNewTokenAuthClient_invalidKey(t *testing.T) {
	_, err := NewTokenAuthClient("blah")
	require.Error(t, err)
}

func TestAdminTokenSource(t *testing.T) {
	ts := adminTokenSource{Key: ExampleAdminKey}
	tok, err := ts.Token()
	require.NoError(t, err)
	require.NotNil(t, tok)
	require.True(t, time.Now().Add(time.Minute*5).After(tok.Expiry))
	require.Equal(t, "Ghost", tok.TokenType)
}
