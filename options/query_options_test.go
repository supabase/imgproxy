package options

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/imgproxy/imgproxy/v3/config"
	"github.com/imgproxy/imgproxy/v3/imagetype"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type QueryProcessingOptionsTestSuite struct{ suite.Suite }

func (s *QueryProcessingOptionsTestSuite) SetupTest() {
	config.Reset()
	config.QueryStringOpts = true
	// Reset presets
	presets = make(map[string]urlOptions)
}

func (s *QueryProcessingOptionsTestSuite) TestParseBase64URL() {
	originURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/%s.png?size=100:100", base64.RawURLEncoding.EncodeToString([]byte(originURL)))
	po, imageURL, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *QueryProcessingOptionsTestSuite) TestParseBase64URLWithoutExtension() {
	originURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/%s?size=100:100", base64.RawURLEncoding.EncodeToString([]byte(originURL)))
	po, imageURL, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.Unknown, po.Format)
}

func (s *QueryProcessingOptionsTestSuite) TestParseBase64URLWithBase() {
	config.BaseURL = "http://images.dev/"

	originURL := "lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/%s.png?size=100:100", base64.RawURLEncoding.EncodeToString([]byte(originURL)))
	po, imageURL, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), fmt.Sprintf("%s%s", config.BaseURL, originURL), imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *QueryProcessingOptionsTestSuite) TestParsePlainURL() {
	originURL := "http://images.dev/lorem/ipsum.jpg"
	path := fmt.Sprintf("/plain/%s@png?size=100:100", originURL)
	po, imageURL, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *QueryProcessingOptionsTestSuite) TestParsePlainURLWithoutExtension() {
	originURL := "http://images.dev/lorem/ipsum.jpg"
	path := fmt.Sprintf("/plain/%s?size=100:100", originURL)

	po, imageURL, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.Unknown, po.Format)
}
func (s *QueryProcessingOptionsTestSuite) TestParsePlainURLEscaped() {
	originURL := "http://images.dev/lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/plain/%s@png?size=100:100", url.PathEscape(originURL))
	po, imageURL, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), originURL, imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *QueryProcessingOptionsTestSuite) TestParsePlainURLWithBase() {
	config.BaseURL = "http://images.dev/"

	originURL := "lorem/ipsum.jpg"
	path := fmt.Sprintf("/plain/%s@png?size=100:100", originURL)
	po, imageURL, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), fmt.Sprintf("%s%s", config.BaseURL, originURL), imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *QueryProcessingOptionsTestSuite) TestParsePlainURLEscapedWithBase() {
	config.BaseURL = "http://images.dev/"

	originURL := "lorem/ipsum.jpg?param=value"
	path := fmt.Sprintf("/plain/%s@png?size=100:100", url.PathEscape(originURL))
	po, imageURL, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)
	require.Equal(s.T(), fmt.Sprintf("%s%s", config.BaseURL, originURL), imageURL)
	require.Equal(s.T(), imagetype.PNG, po.Format)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryFormat() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?format=webp"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), imagetype.WEBP, po.Format)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryResize() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?resize=fill:100:200:1"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), ResizeFill, po.ResizingType)
	require.Equal(s.T(), 100, po.Width)
	require.Equal(s.T(), 200, po.Height)
	require.True(s.T(), po.Enlarge)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryResizingType() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?resizing_type=fill"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), ResizeFill, po.ResizingType)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQuerySize() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?size=100:200:1"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 100, po.Width)
	require.Equal(s.T(), 200, po.Height)
	require.True(s.T(), po.Enlarge)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryWidth() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?width=100"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 100, po.Width)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryHeight() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?height=100"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 100, po.Height)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryEnlarge() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?enlarge=1"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Enlarge)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryExtend() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?extend=1:so:10:20"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), true, po.Extend.Enabled)
	require.Equal(s.T(), GravitySouth, po.Extend.Gravity.Type)
	require.Equal(s.T(), 10.0, po.Extend.Gravity.X)
	require.Equal(s.T(), 20.0, po.Extend.Gravity.Y)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryGravity() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?gravity=soea"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), GravitySouthEast, po.Gravity.Type)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryGravityFocuspoint() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?gravity=fp:0.5:0.75"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), GravityFocusPoint, po.Gravity.Type)
	require.Equal(s.T(), 0.5, po.Gravity.X)
	require.Equal(s.T(), 0.75, po.Gravity.Y)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryQuality() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?quality=55"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 55, po.Quality)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryBackground() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?background=128:129:130"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Flatten)
	require.Equal(s.T(), uint8(128), po.Background.R)
	require.Equal(s.T(), uint8(129), po.Background.G)
	require.Equal(s.T(), uint8(130), po.Background.B)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryBackgroundHex() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?background=ffddee"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Flatten)
	require.Equal(s.T(), uint8(0xff), po.Background.R)
	require.Equal(s.T(), uint8(0xdd), po.Background.G)
	require.Equal(s.T(), uint8(0xee), po.Background.B)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryBackgroundDisable() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?background=fff&background="
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.False(s.T(), po.Flatten)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQueryBlur() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?blur=0.2"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), float32(0.2), po.Blur)
}

func (s *QueryProcessingOptionsTestSuite) TestParseQuerySharpen() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?sharpen=0.2"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), float32(0.2), po.Sharpen)
}
func (s *QueryProcessingOptionsTestSuite) TestParseQueryDpr() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?dpr=2"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.Equal(s.T(), 2.0, po.Dpr)
}
func (s *QueryProcessingOptionsTestSuite) TestParseQueryWatermark() {
	path := "/plain/http://images.dev/lorem/ipsum.jpg?watermark=0.5:soea:10:20:0.6"
	po, _, err := ParseQuery(path, make(http.Header))

	require.Nil(s.T(), err)

	require.True(s.T(), po.Watermark.Enabled)
	require.Equal(s.T(), GravitySouthEast, po.Watermark.Gravity.Type)
	require.Equal(s.T(), 10.0, po.Watermark.Gravity.X)
	require.Equal(s.T(), 20.0, po.Watermark.Gravity.Y)
	require.Equal(s.T(), 0.6, po.Watermark.Scale)
}

func TestQueryProcessingOptions(t *testing.T) {
	suite.Run(t, new(QueryProcessingOptionsTestSuite))
}
